package quotation_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/quotation/handler"
	handlers_report "larsa-e-commerce-microservice-v2/internal/modules/quotation/handler/report"
	repos "larsa-e-commerce-microservice-v2/internal/modules/quotation/repo"
	quotation_report_repos "larsa-e-commerce-microservice-v2/internal/modules/quotation/repo/report"
	quotation_service "larsa-e-commerce-microservice-v2/internal/modules/quotation/service"
	quotation_report_service "larsa-e-commerce-microservice-v2/internal/modules/quotation/service/report"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewQuotationRepository)
	do.Provide(injector, quotation_report_repos.NewQuotationSummaryRepository)

	// services
	do.Provide(injector, quotation_service.NewQuotationBuyerService)
	do.Provide(injector, quotation_service.NewQuotationVendorService)
	do.Provide(injector, quotation_service.NewQuotationAdminService)
	do.Provide(injector, quotation_report_service.NewQuotationSummaryService)
	do.Provide(injector, quotation_service.NewQuotationOrderConvertorService)

	// handlers
	handlers.NewQuotationBuyerHandler(injector, r)
	handlers.NewQuotationAdminHandler(injector, r)
	handlers.NewQuotationSettingsHandler(injector, r)
	handlers.NewQuotationAdminSettingsHandler(injector, r)
	handlers_report.NewQuotationSummaryHandler(injector, r)

	return injector
}
