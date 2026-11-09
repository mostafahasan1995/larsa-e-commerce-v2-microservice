package cron

import (
	"context"
	"fmt"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/events"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	cart_dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
	cart_filters "larsa-e-commerce-microservice-v2/internal/modules/carts/filters"
	cart_models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	cart_services "larsa-e-commerce-microservice-v2/internal/modules/carts/services"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"

	"github.com/robfig/cron/v3"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReminderCronService interface {
	Start()
	CheckPaymentReminders(ctx context.Context, databaseName string)
	CheckCartAbandonmentReminders(ctx context.Context, databaseName string)
}

type reminderCronService struct {
	orderRepo   order_repos.OrderRepo
	cartRepo    cart_services.CartRepo
	settingRepo config_services.SettingRepo
	mongoClient *mongo.Client
}

func NewReminderCronService(i *do.Injector) (ReminderCronService, error) {
	return &reminderCronService{
		orderRepo:   do.MustInvoke[order_repos.OrderRepo](i),
		cartRepo:    do.MustInvoke[cart_services.CartRepo](i),
		settingRepo: do.MustInvoke[config_services.SettingRepo](i),
		mongoClient: do.MustInvoke[*mongo.Client](i),
	}, nil
}

func (service *reminderCronService) Start() {
	// Run the check every 5 minutes (configurable)
	// Ensure distributed lock TTL index exists (idempotent).
	// We use a Mongo collection with a TTL index to implement a simple distributed lock.
	// Each job+database pair creates a lock document with a short expiry, so only
	// one instance can insert it at any moment. Other instances will skip processing.
	service.ensureLocksIndex(context.Background())
	c := cron.New()
	c.AddFunc("@every 5m", func() {
		fmt.Println("Starting reminder cron check process")
		dbs, err := service.settingRepo.GetRelevantDbs()
		if err != nil {
			helpers.DefaultError(err)
			return
		}
		for _, db := range dbs {
			// Attempt to acquire a per-db lock before processing.
			// The key uniquely identifies this cron job for a given tenant database.
			lockKey := fmt.Sprintf("reminder-cron:%s", db)
			acquired := service.tryAcquireLock(context.Background(), lockKey, 4*time.Minute+30*time.Second)
			if !acquired {
				// Another instance already holds the lock for this db in this window.
				fmt.Printf("Skipping reminders for database %s: lock held by another instance\n", db)
				continue
			}
			fmt.Printf("Checking reminders for database: %s\n", db)
			// Check payment reminders
			service.CheckPaymentReminders(context.TODO(), db)
			// Check cart abandonment reminders
			service.CheckCartAbandonmentReminders(context.TODO(), db)
		}
	})

	c.Start()
}

// ensureLocksIndex creates a TTL index on the locks collection (idempotent).
// A TTL index automatically deletes expired lock documents, so a crashed instance
// won't hold a lock forever. Setting expireAfterSeconds to 0 makes Mongo use the
// expiresAt field as an absolute expiry time.
func (service *reminderCronService) ensureLocksIndex(ctx context.Context) {
	collection := service.mongoClient.Database("admin").Collection("cron_locks")
	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0).SetName("ttl_expiresAt"),
	}
	_, err := collection.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		// Best-effort; ignore if index already exists or insufficient permissions
		helpers.DefaultError(err)
	}
}

// tryAcquireLock attempts to insert a lock document with a TTL. If the document already exists,
// another instance holds the lock and we return false. TTL ensures automatic release.
func (service *reminderCronService) tryAcquireLock(ctx context.Context, key string, ttl time.Duration) bool {
	collection := service.mongoClient.Database("admin").Collection("cron_locks")
	now := time.Now()
	expiresAt := now.Add(ttl)

	// Acquire the lock by updating the document if it is expired, or inserting if it does not exist.
	// This avoids relying on Mongo's TTL cleanup timing.
	filter := bson.M{
		"_id": key,
		"$or": []bson.M{
			{"expiresAt": bson.M{"$lte": now}},
			{"expiresAt": bson.M{"$exists": false}},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"expiresAt": expiresAt,
		},
		"$setOnInsert": bson.M{
			"createdAt": now,
		},
	}
	opts := options.Update().SetUpsert(true)

	res, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		helpers.DefaultError(err)
		return false
	}

	if res.UpsertedCount == 1 || res.ModifiedCount == 1 {
		//fmt.Printf("Acquired lock %s until %s\n", key, expiresAt.Format(time.RFC3339))
		return true
	}

	fmt.Printf("Skipped lock %s (already held)\n", key)
	return false
}

func (service *reminderCronService) CheckPaymentReminders(ctx context.Context, databaseName string) {
	// Ensure tenant-aware context always exists
	if ctx == nil {
		ctx = service.settingRepo.AddContext(databaseName)
	}
	if cfg, _ := util.GetReqAppCfg(ctx); cfg == nil || cfg.Db == "" {
		ctx = service.settingRepo.AddContext(databaseName)
	}

	// Calculate time range: orders created in the last 5 minutes
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	// Build the filter to find orders that match our criteria
	match := bson.M{
		"updatedAt":                 bson.M{"$gte": fiveMinutesAgo},
		"paymentInfo.paymentMethod": order_enums.ONLINE_PAYMENT,
		"paymentInfo.status": bson.M{"$in": []order_enums.PaymentStatus{
			order_enums.UNPAID,
			order_enums.PARTIALLY_PAID,
		}},
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"updatedAt": -1}}, // latest first
	}

	orders, apiErr := service.orderRepo.GetAll(ctx, pipeline)
	if apiErr != nil {
		fmt.Printf("Failed to get unpaid online orders: %v\n", apiErr)
		helpers.DefaultError(fmt.Errorf("failed to get unpaid online orders: %v", apiErr))
		return
	}

	//fmt.Printf("Found %d unpaid online orders created in the last 5 minutes\n", len(orders))

	// Send reminder notifications for each order
	for _, order := range orders {
		service.sendPaymentReminder(order, databaseName)
	}
}

func (service *reminderCronService) CheckCartAbandonmentReminders(ctx context.Context, databaseName string) {
	// Ensure tenant-aware context always exists
	if ctx == nil {
		ctx = service.settingRepo.AddContext(databaseName)
	}
	if cfg, _ := util.GetReqAppCfg(ctx); cfg == nil || cfg.Db == "" {
		ctx = service.settingRepo.AddContext(databaseName)
	}

	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	match := bson.M{
		"updatedAt": bson.M{"$gte": fiveMinutesAgo},
		"isActive":  true,
		"deletedAt": bson.M{"$exists": false},
	}

	customPipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"updatedAt": -1}}, // oldest first
	}

	cartFilter := cart_filters.CartFilters{}
	cartFilter.PopFields = &[]string{"cartItemIds"}

	pipeline := cartFilter.BuildPipeline(match, ctx, true)
	pipeline = append(pipeline, customPipeline...)

	cfg, _ := util.GetReqAppCfg(ctx)
	if cfg == nil {
		fmt.Printf("Failed to get app config\n")
		helpers.DefaultError(fmt.Errorf("failed to get app config"))
		return
	}

	// Get the cart collection
	collection := service.mongoClient.Database(cfg.Db).Collection(enums.CART.GetCollectionName())

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Printf("Failed to get abandoned carts: %v\n", err)
		helpers.DefaultError(fmt.Errorf("failed to get abandoned carts: %v", err))
		return
	}
	defer cursor.Close(ctx)

	var carts []cart_models.CartResponse
	if err := cursor.All(ctx, &carts); err != nil {
		fmt.Printf("Failed to decode abandoned carts: %v\n", err)
		helpers.DefaultError(fmt.Errorf("failed to decode abandoned carts: %v", err))
		return
	}

	//fmt.Printf("Found %d carts updated within the last 5 minutes\n", len(carts))
	for _, cart := range carts {
		service.sendCartAbandonmentReminder(ctx, cart, databaseName)
	}
}

func (service *reminderCronService) sendPaymentReminder(order order_models.OrderResponse, databaseName string) {
	events.PublishAccessCheckoutEvent(databaseName, order)

	fmt.Printf("Sent payment reminder for order: %s (Code: %s)\n", order.Id.Hex(), order.Code)
}

func (service *reminderCronService) sendCartAbandonmentReminder(ctx context.Context, cart cart_models.CartResponse, databaseName string) {
	if len(cart.CartItems) > 0 {
		firstItem := cart.CartItems[0]

		// Create UpsertCartItemDto from the first cart item
		cartItemDto := &cart_dtos.UpsertCartItemDto{
			ItemId:   firstItem.ItemId,
			Quantity: 1,
			Type:     firstItem.Type,
			BranchId: firstItem.BranchId,
			UnitId:   firstItem.UnitId,
		}

		events.PublishAddToCartEvent(databaseName, cart.UserId.Hex(), "test@test.com", cartItemDto)

		fmt.Printf("Sent cart abandonment reminder for cart: %s (User: %s)\n", cart.Id.Hex(), cart.UserId.Hex())
	}
}
