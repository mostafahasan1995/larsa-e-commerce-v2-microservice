package seed_services

import (
	"context"
	"fmt"
	"sync"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	seed_utils "larsa-e-commerce-microservice-v2/internal/modules/seed/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type SeedService interface {
	SeedAll(ctx context.Context) error
	SeedSystemData(ctx context.Context) error
	TruncateAllData(ctx context.Context) error
}

type seedService struct {
	db *mongo.Client
}

func NewSeedService(i *do.Injector) (SeedService, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &seedService{
		db: db,
	}, nil
}

func (s *seedService) SeedAll(ctx context.Context) error {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Start a session
	session, err := s.db.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	// Start a transaction
	wc := &writeconcern.WriteConcern{W: "majority"}
	txnOpts := options.Transaction().SetWriteConcern(wc)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := s.db.Database(cfg.Db)
		if err := seed_utils.SeedAll(sessCtx, db, false); err != nil {
			return nil, err
		}
		return nil, nil
	}, txnOpts)

	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}

func (s *seedService) SeedSystemData(ctx context.Context) error {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Start a session
	session, err := s.db.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	// Start a transaction
	wc := &writeconcern.WriteConcern{W: "majority"}
	txnOpts := options.Transaction().SetWriteConcern(wc)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := s.db.Database(cfg.Db)
		if err := seed_utils.SeedAll(sessCtx, db, true); err != nil {
			return nil, err
		}
		return nil, nil
	}, txnOpts)

	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}

func (s *seedService) TruncateAllData(ctx context.Context) error {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Start a session
	session, err := s.db.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	// Start a transaction
	wc := &writeconcern.WriteConcern{W: "majority"}
	txnOpts := options.Transaction().SetWriteConcern(wc)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := s.db.Database(cfg.Db)

		// Create a channel to collect errors
		errChan := make(chan error, 3)

		// Create a wait group to wait for all goroutines to complete
		var wg sync.WaitGroup
		wg.Add(3)

		// Delete from category collection
		go func() {
			defer wg.Done()
			_, err := db.Collection(enums.CATEGORY.GetCollectionName()).DeleteMany(sessCtx, bson.M{})
			errChan <- err
		}()

		// Delete from product collection
		go func() {
			defer wg.Done()
			_, err := db.Collection(enums.PRODUCT.GetCollectionName()).DeleteMany(sessCtx, bson.M{})
			errChan <- err
		}()

		// Update production line collection
		go func() {
			defer wg.Done()
			_, err := db.Collection(enums.PRODUCTION_LINE.GetCollectionName()).UpdateMany(sessCtx, bson.M{}, bson.M{"$set": bson.M{"productIds": []bson.M{}}})
			errChan <- err
		}()

		// Wait for all goroutines to complete
		wg.Wait()
		close(errChan)

		// Check for any errors
		for err := range errChan {
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}, txnOpts)

	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}
