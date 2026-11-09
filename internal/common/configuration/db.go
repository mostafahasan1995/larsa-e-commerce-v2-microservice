package configuration

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func InitDB() (*mongo.Client, error) {
	dbFullHost := EnvironmentVariables.DBFullHost
	dbPass := EnvironmentVariables.DBPassword
	dbUser := EnvironmentVariables.DBUser
	host := EnvironmentVariables.DBHost
	dbPort := EnvironmentVariables.DBPort

	var uri string
	if dbFullHost != "" {
		uri = fmt.Sprintf("mongodb://%s", dbFullHost)
	} else if dbUser != "" && dbPass != "" {
		dbPortInt, _ := strconv.Atoi(dbPort)
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d", dbUser, dbPass, host, dbPortInt)
	} else {
		dbPortInt, _ := strconv.Atoi(dbPort)
		uri = fmt.Sprintf("mongodb://%s:%d", host, dbPortInt)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// logger := log.New(os.Stdout, "mongo: ", log.LstdFlags)

	// Configure connection pool settings
	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second).
		SetConnectTimeout(10 * time.Second).
		SetSocketTimeout(30 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	// Initialize mutation hook and set it as the monitor BEFORE creating the client
	moniteringHook := NewMoniteringHook(nil) // Pass nil since client doesn't exist yet
	clientOptions.SetMonitor(moniteringHook.GetCommandMonitor())

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Update the mutation hook with the client after connection
	moniteringHook.UpdateClient(client)

	// Verify the connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("MongoDB successfully connected and pinged.")

	fmt.Printf("Config Mode: %v\n", EnvironmentVariables.Mode)
	// Set up a background goroutine to monitor connection health
	// if EnvironmentVariables.Mode == "debug" {
	// 	go func() {
	// 		ticker := time.NewTicker(60 * time.Second)
	// 		defer ticker.Stop()

	// 		for range ticker.C {
	// 			pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 			if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
	// 				logger.Printf("MongoDB connection health check failed: %v", err)
	// 			}
	// 			pingCancel()
	// 		}
	// 	}()
	// }

	return client, nil
}

type WithTxn struct {
	db *mongo.Client
}

func (w *WithTxn) Exec(ctx context.Context, callback func(ctx mongo.SessionContext) (any, error)) (any, error) {
	sessionCtx, ok := ctx.(mongo.SessionContext)
	if ok {
		fmt.Println("already session context")
		return callback(sessionCtx)
	}

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	// Starts a session on the client
	session, err := w.db.StartSession()
	if err != nil {
		return nil, err
	}
	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(ctx)
	return session.WithTransaction(ctx, callback, txnOptions)
}

func NewWithTxn(i *do.Injector) (*WithTxn, error) {
	return &WithTxn{
		db: do.MustInvoke[*mongo.Client](i),
	}, nil
}
