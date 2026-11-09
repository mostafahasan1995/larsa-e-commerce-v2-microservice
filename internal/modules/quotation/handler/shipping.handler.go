package quotation_handler

// import (
// 	"larsa-e-commerce-microservice-v2/internal/common/enums"
// 	"larsa-e-commerce-microservice-v2/internal/common/helpers"
// 	"larsa-e-commerce-microservice-v2/internal/common/middleware"
// 	"larsa-e-commerce-microservice-v2/internal/common/types"
// 	"larsa-e-commerce-microservice-v2/internal/common/util"
// 	dtos "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
// 	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
// 	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
// 	quotation_service "larsa-e-commerce-microservice-v2/internal/modules/quotation/service"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-playground/validator/v10"
// 	"github.com/goccy/go-json"
// 	"github.com/samber/do"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type ShippingHandler struct {
// 	shippingService    quotation_service.ShippingService
// 	validationInstance *validator.Validate
// }

// func NewShippingHandler(i *do.Injector, r *chi.Mux) {
// 	h := ShippingHandler{
// 		shippingService:    do.MustInvoke[quotation_service.ShippingService](i),
// 		validationInstance: do.MustInvoke[*validator.Validate](i),
// 	}

// 	r.Route(enums.SHIPPING.GetAdminRoute(), func(r chi.Router) {
// 		r.With(
// 			middleware.AdminAuth(middleware.Restriction{
// 				SuperAdminCapability: enums.CapabilitiesRegistry.Shipping.SuperAdmin.Post.String(),
// 				OwnerCapability:      enums.CapabilitiesRegistry.Shipping.Vendor.Post.String(),
// 				MemberCapability:     enums.CapabilitiesRegistry.Shipping.Member.Post.String(),
// 			}),
// 			middleware.ValidateDTO(dtos.SubmitToAgentDto{}, middleware.BODY, i),
// 		).Post("/agents", helpers.Make(h.SubmitToAgent))
// 		r.With(
// 			middleware.AdminAuth(middleware.Restriction{
// 				SuperAdminCapability: enums.CapabilitiesRegistry.Shipping.SuperAdmin.Put.String(),
// 				OwnerCapability:      enums.CapabilitiesRegistry.Shipping.Vendor.Put.String(),
// 				MemberCapability:     enums.CapabilitiesRegistry.Shipping.Member.Put.String(),
// 			}),
// 			middleware.ValidateDTO(dtos.SubmitToAgentDto{}, middleware.BODY, i),
// 			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
// 		).Put("/{id}/agents", helpers.Make(h.UpdateSubmitToAgentById))
// 		r.With(
// 			middleware.AdminAuth(middleware.Restriction{
// 				SuperAdminCapability: enums.CapabilitiesRegistry.Shipping.SuperAdmin.Post.String(),
// 				OwnerCapability:      enums.CapabilitiesRegistry.Shipping.Vendor.Post.String(),
// 				MemberCapability:     enums.CapabilitiesRegistry.Shipping.Member.Post.String(),
// 			}),
// 			middleware.ValidateDTO(dtos.UpsertShippingAgentDto{}, middleware.BODY, i),
// 			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
// 		).Put("/{id}/agents/from", helpers.Make(h.UpsertShippingFromAgent))
// 		r.With(
// 			middleware.AdminAuth(middleware.Restriction{
// 				SuperAdminCapability: enums.CapabilitiesRegistry.Shipping.SuperAdmin.Post.String(),
// 				OwnerCapability:      enums.CapabilitiesRegistry.Shipping.Vendor.Post.String(),
// 				MemberCapability:     enums.CapabilitiesRegistry.Shipping.Member.Post.String(),
// 			}),
// 			middleware.ValidateDTO(dtos.UpsertShippingAgentDto{}, middleware.BODY, i),
// 			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
// 		).Put("/{id}/agents/to", helpers.Make(h.UpsertShippingToAgent))
// 		r.With(
// 			middleware.AdminAuth(middleware.Restriction{
// 				SuperAdminCapability: enums.CapabilitiesRegistry.Shipping.SuperAdmin.Put.String(),
// 				OwnerCapability:      enums.CapabilitiesRegistry.Shipping.Vendor.Put.String(),
// 				MemberCapability:     enums.CapabilitiesRegistry.Shipping.Member.Put.String(),
// 			}),
// 			middleware.ValidateDTO(dtos.SubmitToVendorDto{}, middleware.BODY, i),
// 			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
// 		).Put("/{id}/vendors", helpers.Make(h.SubmitToVendor))

// 	})
// }

// func (h *ShippingHandler) SubmitToAgent(w http.ResponseWriter, r *http.Request) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	dto, err := middleware.GetValidatedDTO[dtos.SubmitToAgentDto](r)
// 	if err != nil {
// 		helpers.DefaultError(err)
// 		return err
// 	}

// 	result, createErr := h.shippingService.SubmitToAgent(ctx, dto)
// 	if createErr != nil {
// 		helpers.HandleApiError(w, createErr)
// 		return nil
// 	}

// 	response, _ := json.MarshalContext(ctx, *result)
// 	w.Write(response)

// 	return nil
// }

// func (h *ShippingHandler) UpdateSubmitToAgentById(w http.ResponseWriter, r *http.Request) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	ctx, _ := util.AddCtxAppCfg(r)
// 	id := chi.URLParam(r, "id")
// 	objectId, _ := primitive.ObjectIDFromHex(id)

// 	data, err := middleware.GetValidatedDTO[dtos.SubmitToAgentDto](r)
// 	if err != nil {
// 		helpers.DefaultError(err)
// 		return err
// 	}

// 	result, updateErr := h.shippingService.UpdateSubmitToAgentById(ctx, objectId, data)
// 	if updateErr != nil {
// 		helpers.HandleApiError(w, updateErr)
// 		return nil
// 	}

// 	response, _ := json.MarshalContext(ctx, *result)
// 	w.Write(response)

// 	return nil
// }

// func (h *ShippingHandler) UpsertShippingFromAgent(w http.ResponseWriter, r *http.Request) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	dto, err := middleware.GetValidatedDTO[dtos.UpsertShippingAgentDto](r)
// 	if err != nil {
// 		helpers.DefaultError(err)
// 		return err
// 	}

// 	shipping, err := middleware.GetContext[quotation_models.ShippingResponse](r, string(quotation_enums.SHIPPING_CONTEXT_KEY))
// 	if err != nil {
// 		helpers.DefaultError(err)
// 		return err
// 	}

// 	err = h.shippingService.UpsertShippingFromAgent(ctx, *shipping, dto)
// 	if err != nil {
// 		helpers.HandleApiError(w, err)
// 		return nil
// 	}

// 	w.WriteHeader(http.StatusNoContent)

// 	return nil
// }

// func (h *ShippingHandler) UpsertShippingToAgent(w http.ResponseWriter, r *http.Request) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	dto, err := middleware.GetValidatedDTO[dtos.UpsertShippingAgentDto](r)
// 	if err != nil {
// 		helpers.DefaultError(err)
// 		return err
// 	}

// 	shipping, err := middleware.GetContext[quotation_models.ShippingResponse](r, string(quotation_enums.SHIPPING_CONTEXT_KEY))
// 	if err != nil {
// 		helpers.DefaultError(err)
// 		return err
// 	}

// 	err = h.shippingService.UpsertShippingToAgent(ctx, *shipping, dto)
// 	if err != nil {
// 		helpers.HandleApiError(w, err)
// 		return nil
// 	}

// 	w.WriteHeader(http.StatusNoContent)

// 	return nil
// }

// func (h *ShippingHandler) SubmitToVendor(w http.ResponseWriter, r *http.Request) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	dto, err := middleware.GetValidatedDTO[dtos.SubmitToVendorDto](r)
// 	if err != nil {
// 		helpers.DefaultError(err)
// 		return err
// 	}

// 	shipping, err := middleware.GetContext[quotation_models.ShippingResponse](r, string(quotation_enums.SHIPPING_CONTEXT_KEY))
// 	if err != nil {
// 		helpers.DefaultError(err)
// 		return err
// 	}

// 	err = h.shippingService.SubmitToVendor(ctx, *shipping, dto)
// 	if err != nil {
// 		helpers.HandleApiError(w, err)
// 		return nil
// 	}

// 	w.WriteHeader(http.StatusNoContent)

// 	return nil

// }
