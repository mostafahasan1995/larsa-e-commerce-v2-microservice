package base_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
)

// RefField represents a reference field configuration
type SearchQueryField struct {
	FieldName   string
	IsLocalized bool
	IsArray     bool
}

// GetRefFields returns the reference fields configuration for a given module
func GetSearchQueryField(moduleName enums.CollectionName) []SearchQueryField {
	switch moduleName {
	case enums.CATEGORY:
		return []SearchQueryField{
			{
				FieldName:   "name",
				IsLocalized: true,
				IsArray:     false,
			},
			{
				FieldName:   "codes",
				IsLocalized: true,
				IsArray:     true,
			},
		}
	case enums.PRODUCT:
		return []SearchQueryField{
			{
				FieldName:   "name",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.ORDER:
		return []SearchQueryField{
			{
				FieldName:   "code",
				IsLocalized: false,
				IsArray:     false,
			},
		}
	case enums.PROFILE:
		return []SearchQueryField{
			{
				FieldName:   "firstName",
				IsLocalized: false,
				IsArray:     false,
			},
			{
				FieldName:   "lastName",
				IsLocalized: false,
				IsArray:     false,
			},
		}
	case enums.PRODUCTION_LINE:
		return []SearchQueryField{
			{
				FieldName:   "name",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.UNIT:
		return []SearchQueryField{
			{
				FieldName:   "name",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.BRANCH:
		return []SearchQueryField{
			{
				FieldName:   "name",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.ATTRIBUTE:
		return []SearchQueryField{
			{
				FieldName:   "name",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.BROKER:
		return []SearchQueryField{
			{
				FieldName:   "basicInformation.fullName",
				IsLocalized: true,
				IsArray:     false,
			},
			{
				FieldName:   "basicInformation.companyName",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.AGENT:
		return []SearchQueryField{
			{
				FieldName:   "basicInformation.fullName",
				IsLocalized: true,
				IsArray:     false,
			},
			{
				FieldName:   "basicInformation.companyName",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.CUSTOMER:
		return []SearchQueryField{
			{
				FieldName:   "basicInformation.fullName",
				IsLocalized: true,
				IsArray:     false,
			},
			{
				FieldName:   "basicInformation.ownerName",
				IsLocalized: true,
				IsArray:     false,
			},
			{
				FieldName:   "basicInformation.userName",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.QUOTATION:
		return []SearchQueryField{
			{
				FieldName:   "title",
				IsLocalized: false,
				IsArray:     false,
			},
			{
				FieldName:   "code",
				IsLocalized: false,
				IsArray:     false,
			},
		}
	case enums.TRANSACTION:
		return []SearchQueryField{
			{
				FieldName:   "subject",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.ADDON:
		return []SearchQueryField{
			{
				FieldName:   "label",
				IsLocalized: true,
				IsArray:     false,
			},
			{
				FieldName:   "addonValues.label",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.WAITER:
		return []SearchQueryField{
			{
				FieldName:   "name",
				IsLocalized: true,
				IsArray:     false,
			},
		}
	case enums.ZONE:
		return []SearchQueryField{
			{
				FieldName:   "name",
				IsLocalized: false,
				IsArray:     false,
			},
		}

	default:
		return nil
	}
}
