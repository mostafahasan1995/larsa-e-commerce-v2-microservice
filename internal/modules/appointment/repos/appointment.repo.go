package appointment_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	appointment_filters "larsa-e-commerce-microservice-v2/internal/modules/appointment/filters"
	appointment_models "larsa-e-commerce-microservice-v2/internal/modules/appointment/models"
	appointment_service "larsa-e-commerce-microservice-v2/internal/modules/appointment/service"
	appointment_utils "larsa-e-commerce-microservice-v2/internal/modules/appointment/utils"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/samber/do"
)

type AppointmentRepository struct {
	*base_repos.BaseRepo[appointment_models.Appointment, appointment_models.AppointmentResponse]
}

func NewAppointmentRepository(i *do.Injector) (appointment_service.AppointmentRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &AppointmentRepository{
		BaseRepo: base_repos.NewBaseRepo[appointment_models.Appointment, appointment_models.AppointmentResponse](db, enums.APPOINTMENT),
	}, nil
}

func (r *AppointmentRepository) GetWithPagination(ctx context.Context, skip, limit int64, filter appointment_filters.AppointmentFilters) (*appointment_models.AppointmentsWithPagination, *helpers.ApiError) {
	pipeline := filter.BuildPipeline([]bson.M{}, ctx, true)

	props := helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	}

	return r.BaseRepo.GetWithPagination(ctx, props)
}

func (r *AppointmentRepository) GetById(ctx context.Context, id primitive.ObjectID) (*appointment_models.AppointmentResponse, *helpers.ApiError) {
	pipeline := []bson.M{}
	pipeline = append(pipeline, appointment_utils.AppointmentPopulateVendorAndBuyer...)
	//TODO: FEX ME
	// pipeline = append(pipeline, appointment_utils.AppointmentPopulateTranslation...)

	return r.BaseRepo.GetById(ctx, id, pipeline)
}

func (r *AppointmentRepository) UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*appointment_models.AppointmentResponse, *helpers.ApiError) {
	collection, err := r.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	updateFields["updatedAt"] = time.Now()
	updateFields["updatedBy"] = userId

	result, _err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if _err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to update appointment",
		}
	}
	if result.MatchedCount == 0 {
		return nil, &helpers.ApiError{
			StatusCode: 404,
			Message:    "Document not found",
		}
	}

	return r.BaseRepo.GetById(ctx, id, []bson.M{})
}
