package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	config_repos "larsa-e-commerce-microservice-v2/internal/modules/config/repos"
	seed_utils "larsa-e-commerce-microservice-v2/internal/modules/seed/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func main() {
	// Define command-line flags
	onlySystemData := flag.Bool("onlySystemData", false, "only System Data seed")
	seedOrders := flag.Bool("seedOrders", false, "seed 100k orders for testing (outside of a transaction)")
	timeout := flag.Duration("timeout", 0, "optional timeout, e.g. 2h, 30m. 0 means default 30s for system seed and no timeout when seedOrders is set")
	flag.Parse()

	// Load configuration
	// cfg := configuration.EnvironmentVariables

	// Initialize MongoDB client
	client, err := configuration.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Create context
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if *seedOrders {
		if *timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), *timeout)
		} else {
			// No timeout by default for large order seeding
			ctx, cancel = context.WithCancel(context.Background())
		}
	} else {
		to := 30 * time.Second
		if *timeout > 0 {
			to = *timeout
		}
		ctx, cancel = context.WithTimeout(context.Background(), to)
	}
	defer cancel()

	settingRepo := config_repos.SettingRepo{
		BaseRepo: base_repos.NewBaseRepo[config_models.Setting, config_models.SettingResponse](client, enums.SETTING),
	}
	relevantDbs, _ := settingRepo.GetRelevantDbs()

	if *seedOrders {
		// Directly seed orders without wrapping in a transaction (safer for large writes)
		for _, dbName := range relevantDbs {
			fmt.Printf("DBClient: %v\n", dbName)
			// database := client.Database(dbName)
			log.Println("Starting orders seeding (100k)...")
			// if err := seed_utils.SeedOrders(ctx, database, primitive.NewObjectID()); err != nil {
			// 	log.Fatalf("Error seeding orders: %v", err)
			// }
			log.Println("Orders seeding completed successfully!")
		}
	} else {
		// Start a session
		session, err := client.StartSession()
		if err != nil {
			log.Fatalf("failed to start session: %v", err)
		}
		defer session.EndSession(ctx)

		// Start a transaction
		wc := &writeconcern.WriteConcern{W: "majority"}
		txnOpts := options.Transaction().SetWriteConcern(wc)

		_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
			for _, db := range relevantDbs {
				fmt.Printf("DBClient: %v\n", db)
				// Get database instance
				database := client.Database(db)
				// Run seeds
				log.Println("Starting seed process...")
				if err := seed_utils.SeedAll(sessCtx, database, *onlySystemData); err != nil {
					return nil, err
				}
				log.Println("Seed process completed successfully!")
			}
			return nil, nil
		}, txnOpts)

		if err != nil {
			log.Fatalf("Transaction failed: %v", err)
		}
	}
	os.Exit(0)
}
