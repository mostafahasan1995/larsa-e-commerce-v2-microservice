package stats_repos

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	payment_enums "larsa-e-commerce-microservice-v2/internal/modules/payment/enums"
	payment_models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatsInvoiceRepository struct {
	*base_repos.BaseRepo[payment_models.Invoice, payment_models.InvoiceResponse]
}

type StatsInvoiceRepo interface {
	CountInv(ctx context.Context, status payment_enums.InvoicePaymentStatus, startDate, endDate time.Time) (int, error)
}

func NewStatsInvoiceRepo(i *do.Injector) (StatsInvoiceRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &StatsInvoiceRepository{
		BaseRepo: base_repos.NewBaseRepo[payment_models.Invoice, payment_models.InvoiceResponse](db, enums.INVOICE),
	}, nil
}

func (r *StatsInvoiceRepository) CountInv(ctx context.Context, status payment_enums.InvoicePaymentStatus, startDate, endDate time.Time) (int, error) {
	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return 0, fmt.Errorf("someting wrong happend while getting collection: %v", apiErr)
	}
	if !startDate.Before(endDate) {
		return 0, fmt.Errorf("startDate must be before endDate")
	}

	filter := bson.M{
		"paymentStatus": status,
		"status":        bson.M{"$ne": payment_enums.INVOICE_REJECTED},
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
