package order

import (
	order_handlers "larsa-e-commerce-microservice-v2/internal/modules/order/handlers"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"
	order_service "larsa-e-commerce-microservice-v2/internal/modules/order/service"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) {
	do.Provide(injector, order_repos.NewOrderRepository)
	do.Provide(injector, order_repos.NewOrderWebRepository)
	do.Provide(injector, order_repos.NewOrderReportsRepository)

	//services
	do.Provide(injector, order_service.NewOrderService)
	do.Provide(injector, order_service.NewOrderWebService)
	do.Provide(injector, order_service.NewOrderReportsService)
	do.Provide(injector, order_service.NewOrderSummaryService)
	do.Provide(injector, order_service.NewOrderUpdateService)
	do.Provide(injector, order_service.NewOrderCancelRefundService)

	order_handlers.NewOrderHandler(injector, r)
	order_handlers.NewOrderWebHandler(injector, r)

}
