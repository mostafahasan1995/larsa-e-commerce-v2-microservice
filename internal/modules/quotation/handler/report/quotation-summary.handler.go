package quotation_report_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	quotation_report_service "larsa-e-commerce-microservice-v2/internal/modules/quotation/service/report"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type QuotationSummaryHandler struct {
	quotationSummaryService quotation_report_service.QuotationSummaryService
	validationInstance      *validator.Validate
}

func NewQuotationSummaryHandler(i *do.Injector, r *chi.Mux) {
	h := QuotationSummaryHandler{
		quotationSummaryService: do.MustInvoke[quotation_report_service.QuotationSummaryService](i),
		validationInstance:      do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.QUOTATION_SUMMARY.GetAdminRoute(), func(r chi.Router) {
		// Get Quick Summary
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Quotation, enums.Get)),
		).Get("/summary", h.GetQuickSummary)

		// Get Recent Activity
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Quotation, enums.Get)),
		).Get("/activity", h.GetRecentActivity)

		// Get Most Requested Products
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Quotation, enums.Get)),
		).Get("/products/most-requested", h.GetMostRequestedProducts)

		// Get Expiring Quotations
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Quotation, enums.Get)),
		).Get("/expiring", h.GetExpiringQuotations)

		// Get Buyer Quotations Reply Insights
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Quotation, enums.Get)),
		).Get("/insights", h.GetBuyerQuotationsReplyInsights)
	})
}

func (h *QuotationSummaryHandler) GetQuickSummary(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	summary, apiErr := h.quotationSummaryService.GetQuickSummary(ctx)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	helpers.WriteJson(w, http.StatusOK, summary)
}

func (h *QuotationSummaryHandler) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	fromDateStr := r.URL.Query().Get("fromDate")
	toDateStr := r.URL.Query().Get("toDate")

	var fromDate *time.Time
	var toDate *time.Time

	if fromDateStr != "" {
		fromDateParsed, err := time.Parse(time.RFC3339, fromDateStr)
		if err == nil {
			fromDate = &fromDateParsed
		}
	}

	if toDateStr != "" {
		toDateParsed, err := time.Parse(time.RFC3339, toDateStr)
		if err == nil {
			toDate = &toDateParsed
		}
	}

	activity, apiErr := h.quotationSummaryService.GetRecentActivity(ctx, fromDate, toDate)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	helpers.WriteJson(w, http.StatusOK, activity)
}

func (h *QuotationSummaryHandler) GetMostRequestedProducts(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	fromDateStr := r.URL.Query().Get("fromDate")
	toDateStr := r.URL.Query().Get("toDate")

	var fromDate *time.Time
	var toDate *time.Time

	if fromDateStr != "" {
		fromDateParsed, err := time.Parse(time.RFC3339, fromDateStr)
		if err == nil {
			fromDate = &fromDateParsed
		}
	}

	if toDateStr != "" {
		toDateParsed, err := time.Parse(time.RFC3339, toDateStr)
		if err == nil {
			toDate = &toDateParsed
		}
	}

	products, apiErr := h.quotationSummaryService.GetMostRequestedProducts(ctx, fromDate, toDate)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	helpers.WriteJson(w, http.StatusOK, products)
}

func (h *QuotationSummaryHandler) GetExpiringQuotations(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	quotations, apiErr := h.quotationSummaryService.GetExpiringQuotations(ctx)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	helpers.WriteJson(w, http.StatusOK, quotations)
}

func (h *QuotationSummaryHandler) GetBuyerQuotationsReplyInsights(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	insights, apiErr := h.quotationSummaryService.GetBuyerQuotationsReplyInsights(ctx)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return
	}

	helpers.WriteJson(w, http.StatusOK, insights)
}
