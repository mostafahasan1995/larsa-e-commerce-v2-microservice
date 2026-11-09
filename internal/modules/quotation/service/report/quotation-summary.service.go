package quotation_report_service

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	quotation_report_repos "larsa-e-commerce-microservice-v2/internal/modules/quotation/repo/report"
	quotation_types "larsa-e-commerce-microservice-v2/internal/modules/quotation/types"

	"github.com/samber/do"
)

type QuotationSummaryService interface {
	GetQuickSummary(ctx context.Context) (*quotation_types.QuotationSummary, *helpers.ApiError)
	GetRecentActivity(ctx context.Context, fromDate *time.Time, toDate *time.Time) ([]*quotation_types.QuotationActivity, *helpers.ApiError)
	GetMostRequestedProducts(ctx context.Context, fromDate *time.Time, toDate *time.Time) ([]*quotation_types.ProductRequestStats, *helpers.ApiError)
	GetExpiringQuotations(ctx context.Context) ([]*quotation_models.QuotationResponse, *helpers.ApiError)
	GetBuyerQuotationsReplyInsights(ctx context.Context) (*quotation_types.QuotationInsights, *helpers.ApiError)
}

type quotationSummary struct {
	repo quotation_report_repos.QuotationSummaryRepo
}

func NewQuotationSummaryService(i *do.Injector) (QuotationSummaryService, error) {
	return &quotationSummary{
		repo: do.MustInvoke[quotation_report_repos.QuotationSummaryRepo](i),
	}, nil
}

func (s *quotationSummary) GetQuickSummary(ctx context.Context) (*quotation_types.QuotationSummary, *helpers.ApiError) {
	config, _ := util.GetReqAppCfg(ctx)
	vendorId := config.User.Id
	// Create a pipeline to count quotations by status
	summary, apiErr := s.repo.GetQuickSummary(ctx, vendorId)
	if apiErr != nil {
		return nil, apiErr
	}
	return summary, nil
}

func (s *quotationSummary) GetRecentActivity(ctx context.Context, fromDate *time.Time, toDate *time.Time) ([]*quotation_types.QuotationActivity, *helpers.ApiError) {
	config, _ := util.GetReqAppCfg(ctx)
	vendorId := config.User.Id
	activities, apiErr := s.repo.GetRecentActivity(ctx, vendorId, fromDate, toDate)
	if apiErr != nil {
		return nil, apiErr
	}
	return activities, nil
}

func (s *quotationSummary) GetMostRequestedProducts(ctx context.Context, fromDate *time.Time, toDate *time.Time) ([]*quotation_types.ProductRequestStats, *helpers.ApiError) {
	config, _ := util.GetReqAppCfg(ctx)
	vendorId := config.User.Id
	products, apiErr := s.repo.GetMostRequestedProducts(ctx, vendorId, fromDate, toDate)
	if apiErr != nil {
		return nil, apiErr
	}
	return products, nil
}

func (s *quotationSummary) GetExpiringQuotations(ctx context.Context) ([]*quotation_models.QuotationResponse, *helpers.ApiError) {
	config, _ := util.GetReqAppCfg(ctx)
	vendorId := config.User.Id
	return s.repo.GetExpiringQuotations(ctx, vendorId)
}

func (s *quotationSummary) GetBuyerQuotationsReplyInsights(ctx context.Context) (*quotation_types.QuotationInsights, *helpers.ApiError) {
	//TODO: implement this
	return &quotation_types.QuotationInsights{
		PositiveFeedbackPercentage: 0,
		NegativeFeedbackPercentage: 0,
	}, nil
}
