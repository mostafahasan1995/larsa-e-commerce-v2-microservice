package base_repos

import (
	"context"
	"errors"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	"net/http"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseRepo provides generic CRUD operations for MongoDB collections
type BaseRepo[Model any, Response any] struct {
	db       *mongo.Client
	collName enums.CollectionName
}

// NewBaseRepo creates a new base repository instance
func NewBaseRepo[Model any, Response any](db *mongo.Client, collName enums.CollectionName) *BaseRepo[Model, Response] {
	return &BaseRepo[Model, Response]{
		db:       db,
		collName: collName,
	}
}

// GetCollection returns the collection for the current context
func (repo *BaseRepo[Model, Response]) GetCollection(ctx context.Context) (*mongo.Collection, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	if cfg == nil || cfg.Db == "" {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "missing app config in context",
		}
	}
	if repo.db == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "mongo client is nil",
		}
	}

	collection := repo.db.Database(cfg.Db).Collection(repo.collName.GetCollectionName())

	if collection == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get collection",
		}
	}
	return collection, nil
}

// GetDB returns the MongoDB client
func (repo *BaseRepo[Model, Response]) GetDB() *mongo.Client {
	return repo.db
}

// GetCollectionName returns the collection name
func (repo *BaseRepo[Model, Response]) GetCollectionName() enums.CollectionName {
	return repo.collName
}

func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

// Create creates a new document in the collection
func (repo *BaseRepo[Model, Response]) Create(ctx context.Context, document Model) (*Response, *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	if _, err := collection.InsertOne(ctx, document); err != nil {

		if IsDup(err) {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusConflict,
				Message:    "document already exists",
			}
		}

		return nil, helpers.DefaultError(err)
	}

	// Convert the inserted document to Response type
	var response Response
	documentBytes, err := bson.Marshal(document)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}
	if err := bson.Unmarshal(documentBytes, &response); err != nil {
		return nil, helpers.DefaultError(err)
	}

	return &response, nil
}

// CreateMany creates a new documents in the collection
func (repo *BaseRepo[Model, Response]) CreateMany(ctx context.Context, documents []Model) (*[]Response, *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	// Convert []Model to []interface{}
	docs := make([]interface{}, len(documents))
	for i, doc := range documents {
		docs[i] = doc
	}

	if _, err := collection.InsertMany(ctx, docs); err != nil {

		if IsDup(err) {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusConflict,
				Message:    "document already exists",
			}
		}

		return nil, helpers.DefaultError(err)
	}

	// Convert the inserted document to Response type
	var responses []Response
	for _, document := range documents {
		documentBytes, err := bson.Marshal(document)
		if err != nil {
			return nil, helpers.DefaultError(err)
		}
		var response Response
		if err := bson.Unmarshal(documentBytes, &response); err != nil {
			return nil, helpers.DefaultError(err)
		}
		responses = append(responses, response)
	}

	return &responses, nil
}

type UpsertArgs[Model any, Response any] struct {
	Filter []bson.M
	Create []Model
	Update []Model
}

func (repo *BaseRepo[Model, Response]) Upsert(ctx context.Context, args UpsertArgs[Model, Response]) *helpers.ApiError {
	collection, _ := repo.GetCollection(ctx)
	if len(args.Create) == 0 && len(args.Update) == 0 {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "no documents to upsert",
		}
	}

	if len(args.Create) != len(args.Update) {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "create and update documents must have the same length",
		}
	}

	storedDoc, err := repo.FindOne(ctx, args.Filter)
	if err != nil {
		return helpers.DefaultError(err)
	}

	if storedDoc == nil {
		// Create new document
		_, apiErr := repo.CreateMany(ctx, args.Create)
		if apiErr != nil {
			return apiErr
		}
	} else {
		// Update existing document
		var docModel map[string]interface{}
		docBytes, err := bson.Marshal(storedDoc)
		if err != nil {
			return helpers.DefaultError(err)
		}
		if err := bson.Unmarshal(docBytes, &docModel); err != nil {
			return helpers.DefaultError(err)
		}

		for _, updateDoc := range args.Update {
			docBytes, err := bson.Marshal(updateDoc)
			if err != nil {
				return helpers.DefaultError(err)
			}
			var updateMap bson.M
			if err := bson.Unmarshal(docBytes, &updateMap); err != nil {
				return helpers.DefaultError(err)
			}

			_, updateErr := collection.UpdateOne(ctx, bson.M{"_id": docModel["_id"]}, bson.M{"$set": updateMap})
			if updateErr != nil {
				return helpers.DefaultError(updateErr)
			}
		}
	}
	return nil
}

type BulkWriteArgs[Model any, Response any] struct {
	Filter bson.M
	Upsert Model
}

func (repo *BaseRepo[Model, Response]) BulkUpsert(ctx context.Context, args []BulkWriteArgs[Model, Response]) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	collection, _ := repo.GetCollection(ctx)

	// Create bulk write models
	for _, doc := range args {

		result, err := repo.FindOne(ctx, []bson.M{
			{
				"$match": doc.Filter,
			},
		})
		if err != nil {
			return helpers.DefaultError(err)
		}

		docMap := make(bson.M)
		docBytes, docErr := bson.Marshal(doc.Upsert)
		if docErr != nil {
			return helpers.DefaultError(err)
		}
		if err := bson.Unmarshal(docBytes, &docMap); err != nil {
			return helpers.DefaultError(err)
		}

		docMap["updatedAt"] = time.Now()
		docMap["updatedBy"] = cfg.User.Id

		if result == nil {
			docMap["createdAt"] = time.Now()
			docMap["createdBy"] = cfg.User.Id
			docMap["_id"] = primitive.NewObjectID()
			if _, err := collection.InsertOne(ctx, docMap); err != nil {
				if IsDup(err) {
					return &helpers.ApiError{
						StatusCode: http.StatusConflict,
						Message:    "document already exists",
					}
				}

				return helpers.DefaultError(err)
			}
		} else {
			// Update existing document
			if _, err := collection.UpdateOne(ctx, doc.Filter, bson.M{"$set": docMap}); err != nil {
				return helpers.DefaultError(err)
			}
		}
	}

	return nil
}

// GetById retrieves a single document by ID
func (repo *BaseRepo[Model, Response]) GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*Response, *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	pipeline = append(pipeline, bson.M{"$match": bson.M{"_id": id}})
	cursor, aggregateErr := collection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, helpers.DefaultError(aggregateErr)
	}
	defer cursor.Close(ctx)

	var result Response
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, helpers.DefaultError(err)
		}
	} else {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    repo.collName.GetDocumentName() + " not found",
		}
	}

	return &result, nil
}

func (repo *BaseRepo[Model, Response]) GetByIds(ctx context.Context, ids []primitive.ObjectID, pipeline []bson.M, additionalOptions base_utils.AdditionalOptions) ([]Response, *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	uniqueIds := make([]primitive.ObjectID, 0, len(ids))
	seen := make(map[primitive.ObjectID]bool)
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			uniqueIds = append(uniqueIds, id)
		}
	}

	pipeline = append(pipeline, bson.M{"$match": bson.M{"_id": bson.M{"$in": ids}}})
	cursor, aggregateErr := collection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, helpers.DefaultError(aggregateErr)
	}
	defer cursor.Close(ctx)

	var results []Response
	if err := cursor.All(ctx, &results); err != nil {
		return nil, helpers.DefaultError(err)
	}

	if additionalOptions.ReturnError {
		if len(results) != len(uniqueIds) {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "one or more " + repo.collName.GetDocumentName() + " not found",
			}
		}
		return results, nil
	}

	return results, nil
}

func (repo *BaseRepo[Model, Response]) FindOne(ctx context.Context, pipeline []bson.M, opts ...*options.AggregateOptions) (*Response, *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	cursor, aggregateErr := collection.Aggregate(ctx, pipeline, opts...)
	if aggregateErr != nil {
		return nil, helpers.DefaultError(aggregateErr)
	}
	defer cursor.Close(ctx)

	var result Response
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, helpers.DefaultError(err)
		}
	} else {
		return nil, nil
	}

	return &result, nil
}

func (repo *BaseRepo[Model, Response]) GetAll(ctx context.Context, pipeline []bson.M) ([]Response, *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	cursor, aggregateErr := collection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, helpers.DefaultError(aggregateErr)
	}
	defer cursor.Close(ctx)

	var results []Response

	if err := cursor.All(ctx, &results); err != nil {
		return nil, helpers.DefaultError(err)
	}

	if results == nil {
		return []Response{}, nil
	}

	return results, nil
}

func (repo *BaseRepo[Model, Response]) FindAll(ctx context.Context, pipeline []bson.M) ([]interface{}, *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	cursor, aggregateErr := collection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, helpers.DefaultError(aggregateErr)
	}
	defer cursor.Close(ctx)

	var results []interface{}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, helpers.DefaultError(err)
	}

	if results == nil {
		return []interface{}{}, nil
	}

	return results, nil
}

// GetWithPagination retrieves paginated documents with an aggregation pipeline
func (repo *BaseRepo[Model, Response]) GetWithPagination(ctx context.Context, buildPaginationPipelineProps helpers.BuildPaginationPipelineProps) (*types.DataWithPagination[Response], *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	//TODO: we need to move pagination pipeline to the last step
	// Add pagination pipeline
	pipeline := helpers.BuildPaginationPipeline(repo.collName, buildPaginationPipelineProps)

	cursor, aggregateErr := collection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, helpers.DefaultError(aggregateErr)
	}
	defer cursor.Close(ctx)

	var results types.DataWithPagination[Response]

	if cursor.Next(ctx) {
		if err := cursor.Decode(&results); err != nil {
			return nil, helpers.DefaultError(err)
		}
	}

	return &results, nil
}

// GetWithOptimizedPagination retrieves paginated documents using separate count and data queries
// This avoids the $facet document size limit issues
func (repo *BaseRepo[Model, Response]) GetWithOptimizedPagination(ctx context.Context, buildPaginationPipelineProps helpers.BuildPaginationPipelineProps) (*types.DataWithPagination[Response], *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	// Get separate count and data pipelines
	countPipeline, dataPipeline := helpers.BuildOptimizedPaginationPipeline(repo.collName, buildPaginationPipelineProps)

	// Execute count query
	countCursor, countErr := collection.Aggregate(ctx, countPipeline)
	if countErr != nil {
		return nil, helpers.DefaultError(countErr)
	}
	defer countCursor.Close(ctx)

	var totalCount int64 = 0
	if countCursor.Next(ctx) {
		var countResult struct {
			TotalCount int64 `bson:"totalCount"`
		}
		if err := countCursor.Decode(&countResult); err != nil {
			return nil, helpers.DefaultError(err)
		}
		totalCount = countResult.TotalCount
	}

	// Execute data query
	dataCursor, dataErr := collection.Aggregate(ctx, dataPipeline)
	if dataErr != nil {
		return nil, helpers.DefaultError(dataErr)
	}
	defer dataCursor.Close(ctx)

	var data []Response
	if err := dataCursor.All(ctx, &data); err != nil {
		return nil, helpers.DefaultError(err)
	}

	if data == nil {
		data = []Response{}
	}

	// Calculate pagination info
	perPage := buildPaginationPipelineProps.Limit
	totalPages := float64(0)
	if totalCount > 0 && perPage > 0 {
		totalPages = float64((totalCount + perPage - 1) / perPage) // Ceiling division
	}

	result := &types.DataWithPagination[Response]{
		Data: data,
		Pagination: types.Pagination{
			TotalCount: totalCount,
			PerPage:    perPage,
			TotalPages: totalPages,
		},
	}

	return result, nil
}

// UpdateById updates an existing document
func (repo *BaseRepo[Model, Response]) UpdateById(ctx context.Context, id primitive.ObjectID, document Model, filter bson.M) (*Response, *helpers.ApiError) {
	collection, _ := repo.GetCollection(ctx)

	update := bson.M{
		"$set": document,
	}

	// Merge the provided filter with the another filters
	filter["_id"] = id

	var response Response

	upsert := false
	after := options.After

	options := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	result := collection.FindOneAndUpdate(ctx, filter, update, options)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    repo.collName.GetDocumentName() + " not found",
			}
		}

		return nil, helpers.DefaultError(result.Err())
	}

	result.Decode(&response)

	return &response, nil
}

func (r *BaseRepo[Model, Response]) UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*Response, *helpers.ApiError) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	updateFields["updatedAt"] = time.Now()
	if userId != primitive.NilObjectID {
		updateFields["updatedBy"] = userId
	}

	update := bson.M{
		"$set": updateFields,
	}

	filter := bson.M{"_id": id}

	var response Response

	upsert := false
	after := options.After

	options := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	result := collection.FindOneAndUpdate(ctx, filter, update, options)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "document not found",
			}
		}
		return nil, helpers.DefaultError(result.Err())
	}

	result.Decode(&response)

	return &response, nil
}

func (r *BaseRepo[Model, Response]) UpdateMany(ctx context.Context, filter bson.M, updateFields map[string]interface{}, userId primitive.ObjectID) *helpers.ApiError {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return err
	}

	updateFields["updatedAt"] = time.Now()
	if userId != primitive.NilObjectID {
		updateFields["updatedBy"] = userId
	}

	update := bson.M{
		"$set": updateFields,
	}

	_, updateErr := collection.UpdateMany(ctx, filter, update)
	if updateErr != nil {
		return helpers.DefaultError(updateErr)
	}

	return nil
}

func (r *BaseRepo[Model, Response]) UpdateOne(ctx context.Context, filter bson.M, updateFields map[string]interface{}, userId primitive.ObjectID) (*Response, *helpers.ApiError) {
	collection, err := r.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	updateFields["updatedAt"] = time.Now()
	updateFields["updatedBy"] = userId

	update := bson.M{
		"$set": updateFields,
	}

	var response Response

	upsert := false
	after := options.After

	options := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	result := collection.FindOneAndUpdate(ctx, filter, update, options)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "document not found",
			}
		}
		return nil, helpers.DefaultError(result.Err())
	}

	result.Decode(&response)

	return &response, nil
}

// DeleteById soft deletes a document
func (repo *BaseRepo[Model, Response]) DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	collection, _ := repo.GetCollection(ctx)

	// Merge the provided filter with the another filters
	filter["_id"] = id

	now := time.Now()
	update := bson.M{"$set": bson.M{
		"deletedAt": now,
		"deletedBy": cfg.User.Id,
	}}

	result, updateErr := collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return helpers.DefaultError(updateErr)
	}

	if result.MatchedCount == 0 {
		return &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "document not found",
		}
	}

	return nil
}

func (repo *BaseRepo[Model, Response]) DeleteOne(ctx context.Context, filter bson.M) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	collection, _ := repo.GetCollection(ctx)

	now := time.Now()
	update := bson.M{"$set": bson.M{
		"deletedAt": now,
		"deletedBy": cfg.User.Id,
	}}

	result, updateErr := collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return helpers.DefaultError(updateErr)
	}

	if result.MatchedCount == 0 {
		return &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "document not found",
		}
	}

	return nil
}

// BulkDelete soft deletes multiple documents
func (repo *BaseRepo[Model, Response]) BulkDelete(ctx context.Context, ids []primitive.ObjectID, filter bson.M) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	collection, _ := repo.GetCollection(ctx)

	// First check if all IDs exist
	filter["_id"] = bson.M{"$in": ids}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to count documents",
		}
	}

	if count != int64(len(ids)) {
		return &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "some documents were not found",
		}
	}

	now := time.Now()
	update := bson.M{"$set": bson.M{
		"deletedAt": now,
		"deletedBy": cfg.User.Id,
		"updatedAt": now,
		"updatedBy": cfg.User.Id,
	}}

	_, err = collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to delete documents",
		}
	}

	return nil
}

// BulkDelete soft deletes multiple documents
func (repo *BaseRepo[Model, Response]) BulkDeleteFilter(ctx context.Context, filter bson.M) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	collection, _ := repo.GetCollection(ctx)

	now := time.Now()
	update := bson.M{"$set": bson.M{
		"deletedAt": now,
		"deletedBy": cfg.User.Id,
		"updatedAt": now,
		"updatedBy": cfg.User.Id,
	}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to delete documents",
		}
	}

	return nil
}

// HardBulkDelete hard deletes multiple documents
func (repo *BaseRepo[Model, Response]) HardBulkDelete(ctx context.Context, filter bson.M) *helpers.ApiError {
	collection, _ := repo.GetCollection(ctx)

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to count documents",
		}
	}

	if count == 0 {
		return &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "documents not found",
		}
	}

	_, err = collection.DeleteMany(ctx, filter)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to delete documents",
		}
	}

	return nil
}

// Count returns the number of documents matching the filter
func (repo *BaseRepo[Model, Response]) Count(ctx context.Context, filter bson.M) (int64, *helpers.ApiError) {
	collection, err := repo.GetCollection(ctx)
	if err != nil {
		return 0, err
	}

	count, errCount := collection.CountDocuments(ctx, filter)
	if errCount != nil {
		return 0, helpers.DefaultError(err)
	}

	return count, nil
}

// Sum returns the sum of a numeric field for documents matching the filter
func (repo *BaseRepo[Model, Response]) Sum(ctx context.Context, filter bson.M, field string) (float64, *helpers.ApiError) {
	collection, err := repo.GetCollection(ctx)
	if err != nil {
		return 0, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"total": bson.M{
				"$sum": "$" + field,
			},
		}}},
	}

	cursor, aggErr := collection.Aggregate(ctx, pipeline)
	if aggErr != nil {
		return 0, helpers.DefaultError(err)
	}
	defer cursor.Close(ctx)

	var result struct {
		Total float64 `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, helpers.DefaultError(err)
		}
		return result.Total, nil
	}

	return 0, nil
}

func (repo *BaseRepo[Model, Response]) Aggregate(ctx context.Context, pipeline any, callback func(cur *mongo.Cursor) error) error {
	collection, err := repo.GetCollection(ctx)
	if err != nil {
		return err
	}

	cur, errAgg := collection.Aggregate(ctx, pipeline)
	if errAgg != nil {
		return errAgg
	}

	if err := callback(cur); err != nil {
		return err
	}

	return nil
}

func (repo *BaseRepo[Model, Response]) BulkWrite(ctx context.Context, writeOps []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	collection, err := repo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	opts := options.BulkWrite().SetOrdered(true)
	result, errIns := collection.BulkWrite(ctx, writeOps, opts)
	if errIns != nil {
		return nil, errIns
	}

	return result, nil
}
