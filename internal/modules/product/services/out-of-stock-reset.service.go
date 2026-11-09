package product_services

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OutOfStockResetService interface {
	Start()
	ResetOutOfStockProducts(ctx context.Context, databaseName string)
}

type outOfStockResetService struct {
	branchRepo        BranchRepo
	branchProductRepo BranchProductRepo
	settingRepo       config_services.SettingRepo
}

func NewOutOfStockResetService(i *do.Injector) (OutOfStockResetService, error) {
	return &outOfStockResetService{
		branchRepo:        do.MustInvoke[BranchRepo](i),
		branchProductRepo: do.MustInvoke[BranchProductRepo](i),
		settingRepo:       do.MustInvoke[config_services.SettingRepo](i),
	}, nil
}

func (service *outOfStockResetService) Start() {
	// Run the reset every day at 6:00 AM
	c := cron.New()
	c.AddFunc("0 6 * * *", func() {
		fmt.Println("Starting out-of-stock reset process")
		dbs, err := service.settingRepo.GetRelevantDbs()
		if err != nil {
			helpers.DefaultError(err)
			return
		}
		for _, db := range dbs {
			service.ResetOutOfStockProducts(nil, db)
		}
	})

	c.Start()
}

func (service *outOfStockResetService) ResetOutOfStockProducts(ctx context.Context, databaseName string) {
	if ctx == nil {
		context := service.settingRepo.AddContext(databaseName)
		ctx = context
	}

	// Get all active branches
	branches, err := service.branchRepo.GetAll(ctx, "", filters.BranchFilters{})
	if err != nil {
		helpers.DefaultError(err)
		return
	}

	for _, branch := range branches {
		// Check if today is a working day for this branch
		if !service.isWorkingDay(branch.OpeningDays) {
			fmt.Printf("Branch %s is not working today, skipping out-of-stock reset\n", branch.Id.Hex())
			continue
		}

		// Reset out-of-stock products for this branch
		service.resetOutOfStockForBranch(ctx, branch.Id)
	}
}

func (service *outOfStockResetService) isWorkingDay(openingDays *models.OpeningDays) bool {
	if openingDays == nil {
		return false
	}

	today := time.Now().Weekday()
	daySchedule, exists := openingDays.Days[today]

	if !exists {
		return false
	}

	// Check if there are any time ranges for today
	return len(daySchedule.TimeRanges) > 0
}

func (service *outOfStockResetService) resetOutOfStockForBranch(ctx context.Context, branchId primitive.ObjectID) {
	// Find all out-of-stock products for this branch that have auto-reset enabled
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"branchId":            branchId,
				"outOfStock":          true,
				"autoResetOutOfStock": true,
			},
		},
	}

	branchProducts, err := service.branchProductRepo.GetAll(ctx, pipeline)
	if err != nil {
		helpers.DefaultError(err)
		return
	}

	if len(branchProducts) == 0 {
		fmt.Printf("No out-of-stock products found for branch %s\n", branchId.Hex())
		return
	}

	// Update all out-of-stock products to in-stock
	err = service.branchProductRepo.UpdateMany(ctx, bson.M{
		"branchId":            branchId,
		"outOfStock":          true,
		"autoResetOutOfStock": true,
	}, map[string]interface{}{
		"outOfStock": false,
	}, primitive.ObjectID{})

	if err != nil {
		helpers.DefaultError(err)
		return
	}

	fmt.Printf("Reset out-of-stock products for branch %s\n", branchId.Hex())
}
