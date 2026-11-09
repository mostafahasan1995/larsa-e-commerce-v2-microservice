package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/coupon/dtos"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/coupon/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/coupon/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CouponHandler struct {
	couponService      services.CouponService
	validationInstance *validator.Validate
}

func NewCouponHandler(i *do.Injector, r *chi.Mux) {
	h := CouponHandler{
		couponService:      do.MustInvoke[services.CouponService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.COUPON.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Coupon.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Coupon.Vendor.Post.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Coupon.Member.Post.String(),
			}),
			middleware.ValidateDTO(dtos.CreateCouponDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Coupon.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Coupon.Vendor.GetOne.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Coupon.Member.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(config_filters.CouponFilters{}, middleware.QUERY, i)).Get("/{id}", h.GetById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(config_filters.CouponFilters{}, middleware.QUERY, i)).Get("/", h.GetWithPagination)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(config_filters.CouponFilters{}, middleware.QUERY, i)).Get("/all", h.GetAll)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Coupon.SuperAdmin.Put.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Coupon.Vendor.Put.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Coupon.Member.Put.String(),
			}),
			middleware.ValidateDTO(dtos.CreateCouponDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Coupon.SuperAdmin.Delete.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Coupon.Vendor.Delete.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Coupon.Member.Delete.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.DeleteById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Coupon.SuperAdmin.DeleteMany.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Coupon.Vendor.DeleteMany.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Coupon.Member.DeleteMany.String(),
			}),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		).Delete("/", helpers.Make(h.BulkDelete))
	})
}

func (h *CouponHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateCouponDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.couponService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CouponHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.couponService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *CouponHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.couponService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *CouponHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, getErr := h.couponService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *CouponHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateCouponDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.couponService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CouponHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.couponService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CouponHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if bulkDeleteErr := h.couponService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
