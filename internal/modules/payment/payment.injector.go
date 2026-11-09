package payment_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/payment/handlers"
	repos "larsa-e-commerce-microservice-v2/internal/modules/payment/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/payment/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewPaymentRepo)
	do.Provide(injector, repos.NewInvoiceRepo)
	do.Provide(injector, repos.NewPayoutRepo)
	do.Provide(injector, repos.NewTransactionRepo)

	// services
	do.Provide(injector, services.NewPaymentService)
	do.Provide(injector, services.NewInvoiceService)
	do.Provide(injector, services.NewPayoutService)
	do.Provide(injector, services.NewTransactionService)

	// handlers
	handlers.NewPaymentHandler(injector, r)
	handlers.NewInvoiceHandler(injector, r)
	handlers.NewPayoutHandler(injector, r)
	handlers.NewTransactionHandler(injector, r)

	return injector
}
