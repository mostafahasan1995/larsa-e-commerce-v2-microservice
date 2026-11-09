package seed_utils

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// SeedAll runs all seed functions for the product module within a single transaction
func SeedAll(ctx context.Context, db *mongo.Database, onlySystemData bool) error {
	// Start a session
	session, err := db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)
	// Create a default user ID for seeding
	defaultUserId := primitive.NewObjectID()
	// Start a transaction
	wc := &writeconcern.WriteConcern{W: "majority"}
	txnOpts := options.Transaction().SetWriteConcern(wc)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if onlySystemData {
			if err := SeedSystemData(sessCtx, db, defaultUserId); err != nil {
				return nil, err
			}
			return nil, nil
		}

		// Seed attributes
		// if err := SeedAttributes(sessCtx, db, defaultUserId); err != nil {
		// 	return nil, err
		// }

		if err := SeedSystemData(sessCtx, db, defaultUserId); err != nil {
			return nil, err
		}

		// Seed products
		// if err := SeedProducts(sessCtx, db, defaultUserId); err != nil {
		// 	return nil, err
		// }

		return nil, nil
	}, txnOpts)

	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}

func SeedSystemData(sessCtx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	// if err := SeedCategories(sessCtx, db, defaultUserId); err != nil {
	// 	return err
	// }

	if err := SeedCurrencies(sessCtx, db, defaultUserId); err != nil {
		return err
	}

	if err := SeedSalesSettings(sessCtx, db, defaultUserId); err != nil {
		return err
	}

	if err := SeedInventorySettings(sessCtx, db, defaultUserId); err != nil {
		return err
	}

	if err := SeedGeneralSettings(sessCtx, db, defaultUserId); err != nil {
		return err
	}

	if err := SeedCountries(sessCtx, db, defaultUserId); err != nil {
		return err
	}

	if err := SeedCities(sessCtx, db, defaultUserId); err != nil {
		return err
	}

	if err := SeedAreas(sessCtx, db, defaultUserId); err != nil {
		return err
	}

	// if err := SeedAttributes(sessCtx, db, defaultUserId); err != nil {
	// 	return err
	// }

	return nil
}
