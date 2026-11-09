package quotation_repos

import (
	"context"
	enums "larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	repo "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_filters "larsa-e-commerce-microservice-v2/internal/modules/quotation/filters"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/samber/do"
)

type QuotationRepo interface {
	Create(ctx context.Context, quotation quotation_models.Quotation) (*quotation_models.QuotationResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter quotation_filters.QuotationFilters) (*quotation_models.QuotationsWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, filter quotation_filters.QuotationFilters) ([]quotation_models.QuotationResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, filter quotation_filters.QuotationFilters) (*quotation_models.QuotationResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, quotation quotation_models.Quotation, filter bson.M) (*quotation_models.QuotationResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	DeleteAll(ctx context.Context) *helpers.ApiError
	UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*quotation_models.QuotationResponse, *helpers.ApiError)
	UpdateQuotationStatus(ctx context.Context, id primitive.ObjectID, status quotation_enums.QuotationStatus, newHistory []quotation_models.QuotationHistory) *helpers.ApiError
	CreateOrderFromQuotation(ctx context.Context, quotationId primitive.ObjectID, order order_models.Order) (*order_models.OrderResponse, *helpers.ApiError)
}

type QuotationRepository struct {
	*repo.BaseRepo[quotation_models.Quotation, quotation_models.QuotationResponse]
	orderRepo order_repos.OrderRepo
}

func NewQuotationRepository(i *do.Injector) (QuotationRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &QuotationRepository{
		BaseRepo:  repo.NewBaseRepo[quotation_models.Quotation, quotation_models.QuotationResponse](db, enums.QUOTATION),
		orderRepo: do.MustInvoke[order_repos.OrderRepo](i),
	}, nil
}

func (repo *QuotationRepository) GetById(ctx context.Context, id primitive.ObjectID, filter quotation_filters.QuotationFilters) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	pipeline := filter.BuildPipeline(ctx, bson.M{}, true)
	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *QuotationRepository) GetAll(ctx context.Context, filter quotation_filters.QuotationFilters) ([]quotation_models.QuotationResponse, *helpers.ApiError) {
	pipeline := filter.BuildPipeline(ctx, bson.M{}, true)
	return repo.BaseRepo.GetAll(ctx, pipeline)
}

func (repo *QuotationRepository) GetWithPagination(ctx context.Context, skip, limit int64, filter quotation_filters.QuotationFilters) (*quotation_models.QuotationsWithPagination, *helpers.ApiError) {
	pipeline := filter.BuildPipeline(ctx, bson.M{}, true)
	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *QuotationRepository) UpdateById(ctx context.Context, id primitive.ObjectID, quotation quotation_models.Quotation, filter bson.M) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	return repo.BaseRepo.UpdateById(ctx, id, quotation, filter)
}

func (repo *QuotationRepository) Create(ctx context.Context, quotation quotation_models.Quotation) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	quotation.History = append(quotation.History, quotation_models.QuotationHistory{
		Date:           time.Now(),
		Action:         quotation_enums.QUOTATION_CREATED,
		AdditionalInfo: &quotation.Title,
	})
	return repo.BaseRepo.Create(ctx, quotation)
}

func (repo *QuotationRepository) UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	updateFields["updatedAt"] = time.Now()
	updateFields["updatedBy"] = userId

	result, updateErr := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if updateErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to update document: " + updateErr.Error(),
		}
	}

	if result.MatchedCount == 0 {
		return nil, &helpers.ApiError{
			StatusCode: 404,
			Message:    "Document not found",
		}
	}

	return repo.BaseRepo.GetById(ctx, id, []bson.M{})
}

func (repo *QuotationRepository) UpdateQuotationStatus(ctx context.Context, id primitive.ObjectID, status quotation_enums.QuotationStatus, newHistory []quotation_models.QuotationHistory) *helpers.ApiError {
	//1- create session
	session, err := repo.BaseRepo.GetDB().StartSession()
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start transaction",
		}
	}
	defer session.EndSession(ctx)

	cfg, _ := util.GetReqAppCfg(ctx)

	//2- update quotation status and history
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, addHistoryErr := repo.UpdateFieldsById(sessCtx, id, map[string]interface{}{
			"history": newHistory,
		}, cfg.User.Id)
		if addHistoryErr != nil {
			return nil, addHistoryErr
		}

		_, updateStatusErr := repo.UpdateFieldsById(sessCtx, id, map[string]interface{}{
			"status": status,
		}, cfg.User.Id)
		if updateStatusErr != nil {
			return nil, updateStatusErr
		}

		return nil, nil
	})
	if sessionErr != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to update quotation status",
		}
	}

	return nil
}

func (repo *QuotationRepository) CreateOrderFromQuotation(ctx context.Context, quotationId primitive.ObjectID, order order_models.Order) (*order_models.OrderResponse, *helpers.ApiError) {
	session, err := repo.BaseRepo.GetDB().StartSession()
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start transaction",
		}
	}
	defer session.EndSession(ctx)

	cfg, _ := util.GetReqAppCfg(ctx)

	var orderResponse *order_models.OrderResponse

	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		_orderResponse, orderErr := repo.orderRepo.Create(sessCtx, order)
		if orderErr != nil {
			return nil, orderErr
		}

		_, updateErr := repo.UpdateFieldsById(sessCtx, quotationId, map[string]interface{}{
			"status": quotation_enums.CONVERTED_TO_ORDER,
		}, cfg.User.Id)
		if updateErr != nil {
			return nil, updateErr
		}

		orderResponse = _orderResponse
		return orderResponse, nil
	}, nil)

	if sessionErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to create order from quotation",
		}
	}
	return orderResponse, nil
}

func (repo *QuotationRepository) DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError {
	// If filter is empty, create a basic one
	if filter == nil {
		filter = bson.M{}
	}

	return repo.BaseRepo.DeleteById(ctx, id, filter)
}

func (repo *QuotationRepository) DeleteAll(ctx context.Context) *helpers.ApiError {
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return err
	}

	// Delete all quotations
	_, deleteErr := collection.DeleteMany(ctx, bson.M{})
	if deleteErr != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete all quotations: " + deleteErr.Error(),
		}
	}

	return nil
}
