package unserved_delivery_request_services

import (
	"context"
	"log"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	unserved_delivery_request_dtos "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/dtos"
	unserved_delivery_request_filters "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/filters"
	unserved_delivery_request_models "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/models"
	unserved_delivery_request_repos "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/repos"
)

type UnservedDeliveryRequestService interface {
	RecordRequest(ctx context.Context, lat, lng float64, userId *primitive.ObjectID, country, city, area, source *string)
	GetWithPagination(ctx context.Context, skip, limit int64, filter *unserved_delivery_request_filters.UnservedDeliveryRequestFilters) (*types.DataWithPagination[unserved_delivery_request_models.UnservedDeliveryRequestResponse], *helpers.ApiError)
	GetAll(ctx context.Context, filter *unserved_delivery_request_filters.UnservedDeliveryRequestFilters) ([]unserved_delivery_request_models.UnservedDeliveryRequestResponse, *helpers.ApiError)
	GetAggregated(ctx context.Context, groupBy string, filter *unserved_delivery_request_filters.UnservedDeliveryRequestFilters) (*unserved_delivery_request_dtos.UnservedDeliveryRequestAggregateResponse, *helpers.ApiError)
}

type unservedDeliveryRequestService struct {
	repo *unserved_delivery_request_repos.UnservedDeliveryRequestRepo
}

func NewUnservedDeliveryRequestService(i *do.Injector) (UnservedDeliveryRequestService, error) {
	repo := do.MustInvoke[*unserved_delivery_request_repos.UnservedDeliveryRequestRepo](i)
	return &unservedDeliveryRequestService{repo: repo}, nil
}

// RecordRequest persists the unserved delivery request to the database (permanent storage).
// It is called when no branch is found for the given coordinates (e.g. GET /branches/by-location).
func (s *unservedDeliveryRequestService) RecordRequest(ctx context.Context, lat, lng float64, userId *primitive.ObjectID, country, city, area, source *string) {
	model := unserved_delivery_request_models.UnservedDeliveryRequest{
		Latitude:  lat,
		Longitude: lng,
		UserId:    userId,
		Country:   country,
		City:      city,
		Area:      area,
		Source:    source,
	}
	_, err := s.repo.Create(ctx, model)
	if err != nil {
		log.Printf("[UNSERVED_DELIVERY_REQUEST] failed to persist unserved delivery request: %v (lat=%.4f lng=%.4f)", err, lat, lng)
	}
}

func (s *unservedDeliveryRequestService) GetWithPagination(ctx context.Context, skip, limit int64, filter *unserved_delivery_request_filters.UnservedDeliveryRequestFilters) (*types.DataWithPagination[unserved_delivery_request_models.UnservedDeliveryRequestResponse], *helpers.ApiError) {
	return s.repo.GetWithPagination(ctx, skip, limit, filter)
}

func (s *unservedDeliveryRequestService) GetAll(ctx context.Context, filter *unserved_delivery_request_filters.UnservedDeliveryRequestFilters) ([]unserved_delivery_request_models.UnservedDeliveryRequestResponse, *helpers.ApiError) {
	return s.repo.GetAll(ctx, filter)
}

func (s *unservedDeliveryRequestService) GetAggregated(ctx context.Context, groupBy string, filter *unserved_delivery_request_filters.UnservedDeliveryRequestFilters) (*unserved_delivery_request_dtos.UnservedDeliveryRequestAggregateResponse, *helpers.ApiError) {
	if groupBy == "" {
		groupBy = "date"
	}
	validGroupBy := map[string]bool{"country": true, "city": true, "area": true, "date": true}
	if !validGroupBy[groupBy] {
		groupBy = "date"
	}

	data, err := s.repo.AggregateByLocationAndDate(ctx, groupBy, filter)
	if err != nil {
		return nil, err
	}
	return &unserved_delivery_request_dtos.UnservedDeliveryRequestAggregateResponse{
		GroupBy: groupBy,
		Data:    data,
	}, nil
}
