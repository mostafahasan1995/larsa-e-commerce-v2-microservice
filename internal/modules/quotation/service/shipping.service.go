package quotation_service

// import (
// 	"context"
// 	"larsa-e-commerce-microservice-v2/internal/common/helpers"
// 	dtos "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
// 	models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"

// 	"github.com/samber/do"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type ShippingRepo interface {
// 	SubmitToAgent(ctx context.Context, shipping *dtos.SubmitToAgentDto) (*models.ShippingResponse, *helpers.ApiError)
// 	UpdateSubmitToAgentById(ctx context.Context, id primitive.ObjectID, shipping *dtos.SubmitToAgentDto) (*models.ShippingResponse, *helpers.ApiError)
// 	UpsertShippingFromAgent(ctx context.Context, shipping models.ShippingResponse, shippingAgent *dtos.UpsertShippingAgentDto) *helpers.ApiError
// 	UpsertShippingToAgent(ctx context.Context, shipping models.ShippingResponse, shippingAgent *dtos.UpsertShippingAgentDto) *helpers.ApiError
// 	SubmitToVendor(ctx context.Context, shipping models.ShippingResponse, shippingVendor *dtos.SubmitToVendorDto) *helpers.ApiError
// }

// type ShippingService interface {
// 	SubmitToAgent(ctx context.Context, shipping *dtos.SubmitToAgentDto) (*models.ShippingResponse, *helpers.ApiError)
// 	UpdateSubmitToAgentById(ctx context.Context, id primitive.ObjectID, shipping *dtos.SubmitToAgentDto) (*models.ShippingResponse, *helpers.ApiError)
// 	UpsertShippingFromAgent(ctx context.Context, shipping models.ShippingResponse, shippingAgent *dtos.UpsertShippingAgentDto) *helpers.ApiError
// 	UpsertShippingToAgent(ctx context.Context, shipping models.ShippingResponse, shippingAgent *dtos.UpsertShippingAgentDto) *helpers.ApiError
// 	SubmitToVendor(ctx context.Context, shipping models.ShippingResponse, shippingVendor *dtos.SubmitToVendorDto) *helpers.ApiError
// }

// type shipping struct {
// 	repo ShippingRepo
// }

// func NewShippingService(i *do.Injector) (ShippingService, error) {
// 	return &shipping{
// 		repo: do.MustInvoke[ShippingRepo](i),
// 	}, nil
// }

// func (p *shipping) SubmitToAgent(ctx context.Context, shipping *dtos.SubmitToAgentDto) (*models.ShippingResponse, *helpers.ApiError) {
// 	return p.repo.SubmitToAgent(ctx, shipping)
// }

// func (p *shipping) UpdateSubmitToAgentById(ctx context.Context, id primitive.ObjectID, shipping *dtos.SubmitToAgentDto) (*models.ShippingResponse, *helpers.ApiError) {
// 	return p.repo.UpdateSubmitToAgentById(ctx, id, shipping)
// }

// func (p *shipping) UpsertShippingFromAgent(ctx context.Context, shipping models.ShippingResponse, shippingAgent *dtos.UpsertShippingAgentDto) *helpers.ApiError {
// 	return p.repo.UpsertShippingFromAgent(ctx, shipping, shippingAgent)
// }

// func (p *shipping) UpsertShippingToAgent(ctx context.Context, shipping models.ShippingResponse, shippingAgent *dtos.UpsertShippingAgentDto) *helpers.ApiError {
// 	return p.repo.UpsertShippingToAgent(ctx, shipping, shippingAgent)
// }

// func (p *shipping) SubmitToVendor(ctx context.Context, shipping models.ShippingResponse, shippingVendor *dtos.SubmitToVendorDto) *helpers.ApiError {
// 	return p.repo.SubmitToVendor(ctx, shipping, shippingVendor)
// }
