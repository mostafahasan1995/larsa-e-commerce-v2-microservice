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

	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type ProductHandler struct {
	productService       services.ProductService
	branchProductService services.BranchProductService
	validationInstance   *validator.Validate
}

func NewProductHandler(i *do.Injector, r *chi.Mux) {
	h := ProductHandler{
		productService:       do.MustInvoke[services.ProductService](i),
		branchProductService: do.MustInvoke[services.BranchProductService](i),
		validationInstance:   do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.PRODUCT.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.Post.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.Post.String(),
			}),
			middleware.ValidateDTO(dtos.CreateProductDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.GetOne.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(product_filters.ProductFilters{}, middleware.QUERY, i)).Get("/{id}", h.GetById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.Get.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.Get.String(),
			}),
			middleware.ValidateDTO(product_filters.ProductFilters{}, middleware.QUERY, i)).Get("/", h.GetWithPagination)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.Get.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.Get.String(),
			}),
			middleware.ValidateDTO(product_filters.ProductFilters{}, middleware.QUERY, i)).Get("/all", h.GetAll)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.Put.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.Put.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.Put.String(),
			}),
			middleware.ValidateDTO(dtos.CreateProductDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.Delete.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.Delete.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.Delete.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.DeleteById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.DeleteMany.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.DeleteMany.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.DeleteMany.String(),
			}),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		).Delete("/", helpers.Make(h.BulkDelete))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.Put.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.Put.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.Put.String(),
			}),
			middleware.ValidateDTO(dtos.LinkProductionLineDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}/production-line", helpers.Make(h.LinkToProductionLine))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.Delete.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.Delete.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.Delete.String(),
			}),
		).Delete("/duplicates", h.DeleteDuplicatesByName)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.GetOne.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}/inventory-products", helpers.Make(h.GenerateInventoryProducts))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.Post.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.Post.String(),
			}),
			middleware.ValidateDTO(dtos.SaveInventoryProductsDto{}, middleware.BODY, i),
		).Post("/inventory-products", helpers.Make(h.SaveInventoryProducts))
	})
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateProductDto](r)

	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.productService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.productService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

}

func (h *ProductHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.productService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, marshalErr := json.MarshalContext(ctx, result)
	if marshalErr != nil {
		helpers.DefaultError(marshalErr)
		return
	}

	w.Write(response)

}
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)
	// Base match stage with deletedAt filter
	matchStage := bson.M{
		"deletedAt": bson.M{"$exists": false},
	}

	// Create filters from query parameters
	productFilters, _ := product_filters.NewProductFilters(query)

	// Build pipeline with filters
	pipeline := productFilters.BuildPipeline(matchStage, ctx, product_filters.BuildPipelineArgs{
		AddExtraOperations: true,
	})
	result, getErr := h.productService.GetAll(ctx, pipeline)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, marshalErr := json.MarshalContext(ctx, result)
	if marshalErr != nil {
		helpers.DefaultError(marshalErr)
		return
	}

	w.Write(response)
}

func (h *ProductHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateProductDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.productService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *ProductHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.productService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if bulkDeleteErr := h.productService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *ProductHandler) LinkToProductionLine(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.LinkProductionLineDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	linkErr := h.productService.LinkToProductionLine(ctx, objectId, *data)
	if linkErr != nil {
		helpers.HandleApiError(w, linkErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *ProductHandler) DeleteDuplicatesByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	deleteErr := h.productService.DeleteDuplicatesByName(ctx)
	if deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) GenerateInventoryProducts(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid product ID",
		}
	}

	// Get the product first
	product, getErr := h.productService.GetById(ctx, id, "")
	if getErr != nil {
		return getErr
	}

	// Check if product is configurable
	if product.ProductType != product_enums.PRODUCT_TYPE_CONFIGURABLE {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Product is not configurable",
		}
	}

	// Generate inventory products
	inventoryProducts, generateErr := h.productService.GenerateInventoryProducts(ctx, product)
	if generateErr != nil {
		return generateErr
	}

	response, _ := json.MarshalContext(ctx, inventoryProducts)
	w.Write(response)

	return nil
}

func (h *ProductHandler) SaveInventoryProducts(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[dtos.SaveInventoryProductsDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	// Save inventory products
	updatedProduct, saveErr := h.productService.SaveInventoryProducts(ctx, data)
	if saveErr != nil {
		helpers.HandleApiError(w, saveErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, updatedProduct)
	w.Write(response)

	return nil
}
