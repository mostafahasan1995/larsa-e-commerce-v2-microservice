package product_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"net/http"

	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BranchHandler struct {
	branchService      services.BranchService
	validationInstance *validator.Validate
}

func NewBranchHandler(i *do.Injector, r *chi.Mux) {
	h := BranchHandler{
		branchService:      do.MustInvoke[services.BranchService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.BRANCH.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Post.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Post.String(),
			}),
			middleware.ValidateDTO(dtos.CreateBranchDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.GetOne.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(product_filters.BranchFilters{}, middleware.QUERY, i)).Get("/{id}", h.GetById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(product_filters.BranchFilters{}, middleware.QUERY, i)).Get("/all", h.GetAll)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Get.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Get.String(),
			}),
			middleware.ValidateDTO(product_filters.BranchFilters{}, middleware.QUERY, i)).Get("/", h.GetWithPagination)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Put.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Put.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Put.String(),
			}),
			middleware.ValidateDTO(dtos.CreateBranchDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Delete.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Delete.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Delete.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.DeleteById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.DeleteMany.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.DeleteMany.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.DeleteMany.String(),
			}),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		).Delete("/", helpers.Make(h.BulkDelete))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Post.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Post.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Post("/{id}/zones/import", h.ImportZones)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Get.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Get.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}/zones/export", h.ExportZones)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Get.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Get.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}/zones/import-example", h.GetZoneImportExample)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Post.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Post.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Post("/{id}/products/import", h.ImportProducts)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Get.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Get.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}/products/export", h.ExportProducts)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Branch.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Branch.Vendor.Get.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Branch.Member.Get.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}/products/import-example", h.GetProductImportExample)

	})
}

func (h *BranchHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateBranchDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.branchService.Create(ctx, *dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *BranchHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, getErr := h.branchService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *BranchHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.branchService.GetAll(ctx, query)
	if err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *BranchHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.branchService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *BranchHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateBranchDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.branchService.UpdateById(ctx, objectId, *data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *BranchHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.branchService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BranchHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if bulkDeleteErr := h.branchService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *BranchHandler) ImportZones(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Get branch ID from URL
	branchIdStr := chi.URLParam(r, "id")
	branchId, err := primitive.ObjectIDFromHex(branchIdStr)
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid branch ID",
		})
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Unable to parse form",
		})
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing file",
		})
		return
	}
	defer file.Close()

	// Import zones
	apiErr := h.branchService.ImportZones(ctx, branchId, fileHeader)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BranchHandler) GetZoneImportExample(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	// Get example Excel file
	excelData, apiErr := h.branchService.GetZoneImportExample(ctx)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	// Set headers for Excel file download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=zone_import_template.xlsx")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(excelData)))

	// Write Excel data
	w.Write(excelData)
}

func (h *BranchHandler) ImportProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Get branch ID from URL
	branchIdStr := chi.URLParam(r, "id")
	branchId, err := primitive.ObjectIDFromHex(branchIdStr)
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid branch ID",
		})
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Unable to parse form",
		})
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing file",
		})
		return
	}
	defer file.Close()

	// Import products
	apiErr := h.branchService.ImportProducts(ctx, branchId, fileHeader)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BranchHandler) ExportProducts(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	// Get branch ID from URL
	branchIdStr := chi.URLParam(r, "id")
	branchId, err := primitive.ObjectIDFromHex(branchIdStr)
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid branch ID",
		})
		return
	}

	// Get exported Excel file
	excelData, apiErr := h.branchService.ExportProducts(ctx, branchId)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	// Set headers for Excel file download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=branch_products_export.xlsx")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(excelData)))

	// Write Excel data
	w.Write(excelData)
}

func (h *BranchHandler) ExportZones(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	// Get branch ID from URL
	branchIdStr := chi.URLParam(r, "id")
	branchId, err := primitive.ObjectIDFromHex(branchIdStr)
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid branch ID",
		})
		return
	}

	// Get exported Excel file
	excelData, apiErr := h.branchService.ExportZones(ctx, branchId)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	// Set headers for Excel file download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=branch_zones_export.xlsx")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(excelData)))

	// Write Excel data
	w.Write(excelData)
}

func (h *BranchHandler) GetProductImportExample(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	// Get example Excel file
	excelData, apiErr := h.branchService.GetProductImportExample(ctx)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	// Set headers for Excel file download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=product_import_template.xlsx")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(excelData)))

	// Write Excel data
	w.Write(excelData)
}

func (h *BranchHandler) LinkAllProductsToTheFirstBranch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	linkErr := h.branchService.LinkAllProductsToTheFirstBranch(ctx)
	if linkErr != nil {
		helpers.HandleApiError(w, linkErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
