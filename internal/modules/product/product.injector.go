package product_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/product/handlers"
	repos "larsa-e-commerce-microservice-v2/internal/modules/product/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewProductRepo)
	do.Provide(injector, repos.NewProductionLineRepo)
	do.Provide(injector, repos.NewAttributeRepo)
	do.Provide(injector, repos.NewWebProductRepo)
	do.Provide(injector, repos.NewBranchRepo)
	do.Provide(injector, repos.NewBranchWebRepo)
	do.Provide(injector, repos.NewGroupRepo)
	do.Provide(injector, repos.NewBranchProductRepo)
	do.Provide(injector, repos.NewAddonRepo)
	do.Provide(injector, repos.NewWebAddonRepo)
	do.Provide(injector, repos.NewWebReviewRepo)
	do.Provide(injector, repos.NewWebBannerRepo)
	do.Provide(injector, repos.NewBannerRepo)
	do.Provide(injector, repos.NewDiscountRepo)

	// services
	do.Provide(injector, services.NewProductService)
	do.Provide(injector, services.NewProductionLineService)
	do.Provide(injector, services.NewAttributeService)
	do.Provide(injector, services.NewWebProductService)
	do.Provide(injector, services.NewBranchService)
	do.Provide(injector, services.NewGroupService)
	do.Provide(injector, services.NewBranchProductService)
	do.Provide(injector, services.NewInventoryProductService)
	do.Provide(injector, services.NewAddonService)
	do.Provide(injector, services.NewBranchWebService)
	do.Provide(injector, services.NewWebAddonService)
	do.Provide(injector, services.NewWebReviewService)
	do.Provide(injector, services.NewWebBannerService)
	do.Provide(injector, services.NewBannerService)
	do.Provide(injector, services.NewDiscountService)
	do.Provide(injector, services.NewOutOfStockResetService)

	// Start the out-of-stock reset service
	outOfStockResetService := do.MustInvoke[services.OutOfStockResetService](injector)
	outOfStockResetService.Start()

	// handlers
	handlers.NewDiscountHandler(injector, r)
	handlers.NewBannerHandler(injector, r)
	handlers.NewWebBannerHandler(injector, r)
	handlers.NewProductHandler(injector, r)
	handlers.NewProductionLineHandler(injector, r)
	handlers.NewAttributeHandler(injector, r)
	handlers.NewWebProductHandler(injector, r)
	handlers.NewBranchHandler(injector, r)
	handlers.NewGroupHandler(injector, r)
	handlers.NewBranchProductHandler(injector, r)
	handlers.NewAddonHandler(injector, r)
	handlers.NewBranchWebHandler(injector, r)
	handlers.NewWebAddonHandler(injector, r)

	return injector
}
