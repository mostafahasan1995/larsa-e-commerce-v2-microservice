package unit_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/coupon/handlers"
	repos "larsa-e-commerce-microservice-v2/internal/modules/coupon/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/coupon/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

// ProvideServices registers all coupon services and repos without creating handlers
func ProvideServices(injector *do.Injector) {
	// repos
	do.Provide(injector, repos.NewCouponRepo)
	do.Provide(injector, repos.NewWebCouponRepo)

	// services
	do.Provide(injector, services.NewCouponService)
	do.Provide(injector, services.NewCouponRuleChecker)
	do.Provide(injector, services.NewWebCouponService)
}

// InitHandlers creates handlers after all services are registered
func InitHandlers(injector *do.Injector, r *chi.Mux) {
	// Start the coupon rule checker service
	// couponRuleChecker := do.MustInvoke[*services.CouponRuleCheckerService](injector)
	// couponRuleChecker.Start()

	// handlers
	handlers.NewCouponHandler(injector, r)
	handlers.NewWebCouponHandler(injector, r)
}

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	ProvideServices(injector)
	InitHandlers(injector, r)
	return injector
}
