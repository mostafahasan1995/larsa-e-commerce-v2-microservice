package profile_handlers

import (
	"net/http"
	"strconv"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	profile_filters "larsa-e-commerce-microservice-v2/internal/modules/profile/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/profile/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProfileHandler handles HTTP requests for profiles
type ProfileHandler struct {
	profileService     services.ProfileService
	validationInstance *validator.Validate
}

// NewProfileHandler creates a new profile handler and sets up routes
func NewProfileHandler(i *do.Injector, r *chi.Mux) {
	h := ProfileHandler{
		profileService:     do.MustInvoke[services.ProfileService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}
	r.Route(enums.PROFILE.GetWebRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				OwnerCapability: enums.CapabilitiesRegistry.Profile.Vendor.Put.String(),
				BuyerCapability: enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(dtos.CreateProfileDto{}, middleware.BODY, i),
		).Put("/me", helpers.Make(h.UpdateMyProfile))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability:      enums.CapabilitiesRegistry.Profile.SuperAdmin.Get.String(),
				VendorCapability:          enums.CapabilitiesRegistry.Profile.Vendor.Get.String(),
				AgentCapability:           enums.CapabilitiesRegistry.Profile.Agent.GetOne.String(),
				ShippingCompanyCapability: enums.CapabilitiesRegistry.Profile.ShippingCompany.GetOne.String(),
				OwnerCapability:           enums.CapabilitiesRegistry.Profile.Vendor.GetOne.String(),
				BuyerCapability:           enums.CapabilitiesRegistry.General,
			}),
		).Get("/me", h.GetMe)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.Post.String(),
			}),
			middleware.ValidateDTO(dtos.CreateVendorProfileDto{}, middleware.BODY, i),
		).Post("/vendor", helpers.Make(h.CreateVendorProfile))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.Put.String(),
			}),
			middleware.ValidateDTO(dtos.CreateVendorProfileDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/vendor/{id}", helpers.Make(h.UpdateVendorProfile))

		// Buyer profile routes
		r.With(
			// middleware.AdminAuth(middleware.Restriction{
			// 	SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.Post.String(),
			// }),
			middleware.ValidateDTO(dtos.CreateBuyerProfileDto{}, middleware.BODY, i),
		).Post("/buyer", helpers.Make(h.CreateBuyerProfile))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.Put.String(),
			}),
			middleware.ValidateDTO(dtos.CreateBuyerProfileDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/buyer/{id}", helpers.Make(h.UpdateBuyerProfile))

		// Guest profile routes
		r.With(
			middleware.ValidateDTO(dtos.CreateGuestProfileDto{}, middleware.BODY, i),
		).Post("/guest", helpers.Make(h.CreateGuestProfile))

		r.Get("/guest/{fingerprintId}", h.GetGuestProfile)
	})
	r.Route(enums.PROFILE.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.Get.String(),
			}),
			middleware.ValidateDTO(profile_filters.ProfileFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.Get.String(),
			}),
		).Get("/", h.GetWithPagination)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.Get.String(),
			}),
		).Get("/user/{userId}", h.GetProfileByUserId)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.Delete.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.DeleteById)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.DeleteMany.String(),
			}),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		).Delete("/", helpers.Make(h.BulkDelete))

		// Delete all guest profiles endpoint
		r.With(
		// middleware.AdminAuth(middleware.Restriction{
		// 	SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.DeleteMany.String(),
		// }),
		).Delete("/guests", h.DeleteAllGuestProfiles)

		// Migration endpoint
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Profile.SuperAdmin.Post.String(),
			}),
		).Post("/migrate-users", helpers.Make(h.MigrateUsersToProfiles))
	})
}

// CreateVendorProfile creates a new vendor profile
func (h *ProfileHandler) CreateVendorProfile(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateVendorProfileDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.profileService.CreateVendorProfile(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

// UpdateVendorProfile updates an existing vendor profile
func (h *ProfileHandler) UpdateVendorProfile(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	dto, err := middleware.GetValidatedDTO[dtos.CreateVendorProfileDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.profileService.UpdateVendorProfile(ctx, objectId, dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

// CreateBuyerProfile creates a new buyer profile
func (h *ProfileHandler) CreateBuyerProfile(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateBuyerProfileDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.profileService.CreateBuyerProfile(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

// UpdateBuyerProfile updates an existing buyer profile
func (h *ProfileHandler) UpdateBuyerProfile(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	dto, err := middleware.GetValidatedDTO[dtos.CreateBuyerProfileDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.profileService.UpdateBuyerProfile(ctx, objectId, dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

// GetById retrieves a profile by ID
func (h *ProfileHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, getErr := h.profileService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// GetWithPagination retrieves profiles with pagination
func (h *ProfileHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.profileService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// GetAll retrieves all profiles
func (h *ProfileHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, getErr := h.profileService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// GetProfileByUserId retrieves a profile by user ID
func (h *ProfileHandler) GetProfileByUserId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	userId := chi.URLParam(r, "userId")
	objectId, _ := primitive.ObjectIDFromHex(userId)
	result, getErr := h.profileService.GetProfileByUserId(ctx, objectId)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	if result == nil {
		// If no profile is found, return 404
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "profile not found for user",
		})
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// DeleteById deletes a profile by ID
func (h *ProfileHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	if deleteErr := h.profileService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// BulkDelete deletes multiple profiles
func (h *ProfileHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if bulkDeleteErr := h.profileService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// DeleteAllGuestProfiles deletes all guest profiles (where firstName equals "Guest")
func (h *ProfileHandler) DeleteAllGuestProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	deletedCount, deleteErr := h.profileService.DeleteAllGuestProfiles(ctx)
	if deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	response := dtos.DeleteGuestProfilesResponseDto{
		DeletedCount: deletedCount,
		Message:      "Guest profiles deleted successfully",
	}

	responseJson, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

// UpdateMyProfile updates the authenticated user's profile
func (h *ProfileHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateProfileDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.profileService.UpdateMyProfile(ctx, dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

// GetMe retrieves the authenticated user's profile
func (h *ProfileHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	result, getErr := h.profileService.GetMe(ctx)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	if result == nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "profile not found",
		})
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// CreateGuestProfile creates a new guest profile
func (h *ProfileHandler) CreateGuestProfile(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateGuestProfileDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.profileService.CreateGuestProfile(&ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}
	accessToken := ctx.Value("accessToken").(string)
	refreshToken := ctx.Value("refreshToken").(string)

	response, _ := json.MarshalContext(ctx, result)
	w.Header().Set("x-access-token", accessToken)
	w.Header().Set("x-refresh-token", refreshToken)
	w.Write(response)

	return nil
}

// GetGuestProfile retrieves a profile by fingerprint ID
func (h *ProfileHandler) GetGuestProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	fingerprintId := chi.URLParam(r, "fingerprintId")

	result, getErr := h.profileService.GetProfileByFingerprintId(ctx, fingerprintId)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	if result == nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "guest profile not found",
		})
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// MigrateUsersToProfiles migrates user data from the user collection to the profile collection
func (h *ProfileHandler) MigrateUsersToProfiles(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	limitStr := r.URL.Query().Get("limit")
	limitInt := int64(100) // default limit
	offsetStr := r.URL.Query().Get("offset")
	offsetInt := int64(0) // default offset

	if limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil && limit > 0 {
			limitInt = limit
		}
	}

	if offsetStr != "" {
		if offset, err := strconv.ParseInt(offsetStr, 10, 64); err == nil && offset >= 0 {
			offsetInt = offset
		}
	}
	result, migrateErr := h.profileService.MigrateUsersToProfilesDirect(ctx, limitInt, offsetInt)
	if migrateErr != nil {
		helpers.HandleApiError(w, migrateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}
