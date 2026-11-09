package product_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	"net/http"

	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type BranchProductHandler struct {
	branchProductService services.BranchProductService
	validationInstance   *validator.Validate
}

func NewBranchProductHandler(i *do.Injector, r *chi.Mux) {
	h := BranchProductHandler{
		branchProductService: do.MustInvoke[services.BranchProductService](i),
		validationInstance:   do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.BRANCH_PRODUCT.GetAdminRoute(), func(r chi.Router) {

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Product.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Product.Vendor.GetOne.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Product.Member.GetOne.String(),
			}),
			middleware.ValidateDTO(product_filters.BranchProductFilters{}, middleware.QUERY, i)).Get("/branches/{id}", h.GetByBranchIdAndProductId)

	})
}

func (h *BranchProductHandler) GetByBranchIdAndProductId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	branchProductFilters, _ := product_filters.NewBranchProductFilters(query)

	result, getErr := h.branchProductService.GetByBranchIdAndProductId(ctx, product_utils.GetByBranchIdAndProductIdArgs{
		BranchId:  objectId,
		ProductId: *branchProductFilters.ProductId,
	}, []bson.M{})
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

}
