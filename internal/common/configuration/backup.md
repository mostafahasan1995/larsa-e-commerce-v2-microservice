package configuration

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
)

// MoniteringHook represents a hook that will be called on database mutations
type MoniteringHook struct {
	client      *mongo.Client
	requestInfo sync.Map // Map to store per-request info (collection name, start time)
}

// NewMoniteringHook creates a new mutation hook instance
func NewMoniteringHook(client *mongo.Client) *MoniteringHook {
	return &MoniteringHook{
		client: client,
	}
}

// UpdateClient updates the client reference in the mutation hook
func (h *MoniteringHook) UpdateClient(client *mongo.Client) {
	h.client = client
}

// GetCommandMonitor returns a command monitor that will track mutations
func (h *MoniteringHook) GetCommandMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			// Store collection name from the command
			var collectionName string

			// Handle getMore command specifically as it has collection name in "collection" field
			if evt.CommandName == "getMore" {
				if val, err := evt.Command.LookupErr("collection"); err == nil {
					if str, ok := val.StringValueOK(); ok {
						collectionName = str
					}
				}
			} else {
				// For most CRUD operations (insert, update, delete, find, aggregate),
				// the first element (command name) contains the collection name
				if val, err := evt.Command.LookupErr(evt.CommandName); err == nil {
					if str, ok := val.StringValueOK(); ok {
						collectionName = str
					}
				}
			}

			h.requestInfo.Store(evt.RequestID, collectionName)

			if isMutationOperation(evt.CommandName) {
				// log.Printf("Mutation operation started: %s %s", evt.CommandName, evt.DatabaseName)
				// log.Printf("Command: %s", evt.Command)
			} else {
				// log.Printf("Non-mutation operation started: %s %s", evt.CommandName, evt.DatabaseName)
				// log.Printf("Command: %s", evt.Command)
			}
		},
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			if evt != nil {
				// Retrieve and clean up request info
				var collectionName string
				if val, ok := h.requestInfo.Load(evt.RequestID); ok {
					collectionName = val.(string)
					h.requestInfo.Delete(evt.RequestID)
				}

				h.handleMonitering(ctx, HandleMoniteringArgs{
					CommandName:          &evt.CommandName,
					DatabaseName:         evt.DatabaseName,
					ServerConnectionID64: evt.ServerConnectionID64,
					RequestID:            &evt.RequestID,
					ServiceID:            evt.ServiceID,
					Command:              &evt.Reply,
					Duration:             &evt.Duration,
					ConnectionID:         &evt.ConnectionID,
					CollectionName:       collectionName,
				})

				if !isExcludedOperation(evt.Reply) {

					if isMutationOperation(evt.CommandName) {
						// log.Printf("Mutation operation succeeded: %s %s %d", evt.CommandName, evt.DatabaseName, evt.Duration.Milliseconds())
					} else {
						// log.Printf("Non-mutation operation succeeded: %s %s %d", evt.CommandName, evt.DatabaseName, evt.Duration.Milliseconds())
					}
				}
			}
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			if evt != nil {
				// Clean up request info on failure
				h.requestInfo.Delete(evt.RequestID)

				h.handleMonitering(ctx, HandleMoniteringArgs{
					CommandName:          &evt.CommandName,
					DatabaseName:         evt.DatabaseName,
					ServerConnectionID64: evt.ServerConnectionID64,
					RequestID:            &evt.RequestID,
					ServiceID:            evt.ServiceID,
					Command:              nil,
					Duration:             &evt.Duration,
					ConnectionID:         &evt.ConnectionID,
				})
			}

			if isMutationOperation(evt.CommandName) {
				// 	log.Printf("Mutation operation failed: %s %s %d - %v", evt.CommandName, evt.DatabaseName, evt.Duration.Milliseconds(), evt.Failure)
			} else {
				// log.Printf("Non-mutation operation failed: %s %s %d - %v", evt.CommandName, evt.DatabaseName, evt.Duration.Milliseconds(), evt.Failure)
			}

		},
	}
}

type HandleMoniteringArgs struct {
	Duration             *time.Duration
	CommandName          *string
	DatabaseName         string
	RequestID            *int64
	ConnectionID         *string
	ServerConnectionID64 *int64
	ServiceID            *primitive.ObjectID
	Command              *bson.Raw
	CollectionName       string
}

// handleMonitering processes a successful mutation operation
func (h *MoniteringHook) handleMonitering(ctx context.Context, args HandleMoniteringArgs) {
	// Skip if this is a monitoring log operation to prevent recursion
	if args.CollectionName == enums.MONITERING_LOG.GetCollectionName() {
		return
	}

	// Create a new document to store mutation information
	cfg, _ := util.GetReqAppCfg(ctx)
	moniteringLog := map[string]interface{}{
		"_id":        primitive.NewObjectID(),
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	if args.ServerConnectionID64 != nil {
		moniteringLog["server_connection_id64"] = args.ServerConnectionID64
	}

	if args.RequestID != nil {
		moniteringLog["request_id"] = args.RequestID
	}

	if args.CommandName != nil {
		moniteringLog["operation"] = args.CommandName
	}

	if args.DatabaseName != "" {
		moniteringLog["database"] = args.DatabaseName
	}

	if args.CollectionName != "" {
		moniteringLog["collection_name"] = args.CollectionName
	}

	if args.Duration != nil {
		moniteringLog["duration_ms"] = args.Duration.Milliseconds()

		// Enhanced performance logging for slow queries with query details
		if args.Duration.Milliseconds() > 2000 { // Log queries taking more than 2000ms
			queryDetails := "N/A"
			collectionName := args.CollectionName
			if collectionName == "" {
				collectionName = "unknown"
			}

			// Extract query details based on operation type
			if args.Command != nil {
				commandStr := args.Command.String()

				// Only stringify for slow queries when we actually need to extract details from the response/cursor
				// Note: We're still stringifying the response here, which is suboptimal for 'find'
				// but at least it's only for slow queries > 2s, not every query.
				// Ideally, we should capture the query from the 'Started' event too.

				if args.CommandName != nil {
					switch *args.CommandName {
					case "find":
						if idx := strings.Index(commandStr, `"filter"`); idx != -1 {
							if end := strings.Index(commandStr[idx:], `}, "`); end != -1 {
								queryDetails = commandStr[idx : idx+end+1]
							}
						}
					case "aggregate":
						if idx := strings.Index(commandStr, `"pipeline"`); idx != -1 {
							if end := strings.Index(commandStr[idx:], `], "`); end != -1 {
								pipeline := commandStr[idx : idx+end+1]
								// Truncate very long pipelines
								if len(pipeline) > 200 {
									pipeline = pipeline[:200] + "..."
								}
								queryDetails = pipeline
							}
						}
					case "update", "delete":
						if idx := strings.Index(commandStr, `"q"`); idx != -1 {
							if end := strings.Index(commandStr[idx:], `}, "`); end != -1 {
								queryDetails = commandStr[idx : idx+end+1]
							}
						}
					}
				}
			}

			fmt.Printf("[DB PERFORMANCE] SLOW QUERY - %s on %s.%s (%v)\n  Query: %s\n",
				*args.CommandName,
				args.DatabaseName,
				collectionName,
				args.Duration,
				queryDetails,
			)
		}
	}

	if args.ConnectionID != nil {
		moniteringLog["connection_id"] = args.ConnectionID
	}

	if cfg != nil {
		if cfg.User != nil {
			moniteringLog["user"] = cfg.User
		}
		if cfg.Lang != "" {
			moniteringLog["lang"] = cfg.Lang
		}
		if cfg.Hp != nil {
			moniteringLog["hp"] = cfg.Hp
		}
	}

	if args.ServiceID != nil {
		moniteringLog["service_id"] = args.ServiceID.Hex()
	}

	// Insert the mutation log into a collection using the new context
	collection := h.client.Database(args.DatabaseName).Collection(enums.MONITERING_LOG.GetCollectionName())
	_, err := collection.InsertOne(ctx, moniteringLog)
	if err != nil {
		// log.Printf("Failed to log monitering: %v", err)
	}
}

// isMutationOperation checks if the given command is a mutation operation
func isMutationOperation(commandName string) bool {
	mutationCommands := map[string]bool{
		"insert":        true,
		"update":        true,
		"delete":        true,
		"findAndModify": true,
		"aggregate":     true, // Only if it includes $out or $merge stages
	}
	return mutationCommands[commandName]
}

func isExcludedOperation(command bson.Raw) bool {
	isExcluded := false
	excludedOperations := []string{
		enums.MONITERING_LOG.GetCollectionName(),
	}
	for _, operation := range excludedOperations {
		if strings.Contains(command.String(), operation) {
			isExcluded = true
			break
		}
	}
	return isExcluded
}
