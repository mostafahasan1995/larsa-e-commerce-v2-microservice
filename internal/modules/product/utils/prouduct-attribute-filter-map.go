package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
)

func ProductAttributeFilterMap(attributeType product_enums.ProductAttributeType) []enums.Operator {

	switch attributeType {
	case product_enums.PRODUCT_ATTRIBUTE_TYPE_STRING:
		return []enums.Operator{
			enums.EQ,
			enums.NE,
			enums.IN,
			enums.NOT_IN,
		}
	case product_enums.PRODUCT_ATTRIBUTE_TYPE_NUMBER:
		return []enums.Operator{
			enums.GT,
			enums.GTE,
			enums.LT,
			enums.LTE,
		}
	case product_enums.PRODUCT_ATTRIBUTE_TYPE_BOOLEAN:
		return []enums.Operator{
			enums.EQ,
			enums.NE,
		}
	case product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY:
		return []enums.Operator{
			enums.IN,
			enums.NOT_IN,
		}
	case product_enums.PRODUCT_ATTRIBUTE_TYPE_OBJECT:
		return []enums.Operator{}
	}

	return []enums.Operator{}
}
