package contact_filters

import (
	shipping_filters "larsa-e-commerce-microservice-v2/internal/modules/contact/filters/shipping"
)

type ShippingFilters = shipping_filters.ShippingFilters

func NewShippingFilters(query string) (*ShippingFilters, error) {
	return shipping_filters.NewShippingFilters(query)
}
