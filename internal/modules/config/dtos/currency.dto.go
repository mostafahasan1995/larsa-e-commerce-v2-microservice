package config_dtos

import "larsa-e-commerce-microservice-v2/internal/common/types"

type CreateCurrencyDto struct {
	Name   types.Localizable[string] `json:"name" validate:"required,localizable"`
	Symbol string                    `json:"symbol" validate:"required"`
}
