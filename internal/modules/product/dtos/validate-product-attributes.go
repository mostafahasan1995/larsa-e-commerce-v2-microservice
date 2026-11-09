package product_dtos

import (
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

func ValidateProductAttributes(attributes []AttributeDto) *helpers.ApiError {
	for _, attr := range attributes {
		attribute := product_models.ProductAttribute{}
		attribute.Code = attr.Code
		attribute.Type = attr.Type

		// Validate type matches value
		switch attr.Type {
		case "string":
			isString := util.IsString(attr.Value)
			if !isString {
				return &helpers.ApiError{
					Message:    fmt.Sprintf("invalid value type for attribute %s: expected string, got %v", attr.Code, attr.Value),
					StatusCode: 400,
				}
			}
			attribute.Value = attr.Value
		case "number":
			isNumber := util.IsNumber(attr.Value)
			if !isNumber {
				return &helpers.ApiError{
					Message:    fmt.Sprintf("invalid value type for attribute %s: expected number, got %v", attr.Code, attr.Value),
					StatusCode: 400,
				}
			}
			attribute.Value = attr.Value
		case "boolean":
			isBoolean := util.IsBoolean(attr.Value)
			if !isBoolean {
				return &helpers.ApiError{
					Message:    fmt.Sprintf("invalid value type for attribute %s: expected boolean, got %v", attr.Code, attr.Value),
					StatusCode: 400,
				}
			}
			attribute.Value = attr.Value
		case "array":
			if arr, ok := attr.Value.([]interface{}); ok {
				// Validate each element in the array
				isArray := util.IsArray(arr)
				if !isArray {
					return &helpers.ApiError{
						Message:    fmt.Sprintf("invalid value type for attribute %s: expected array, got %v", attr.Code, attr.Value),
						StatusCode: 400,
					}
				}
				attribute.Value = arr
			} else {
				return &helpers.ApiError{
					Message:    fmt.Sprintf("invalid value type for attribute %s: expected array, got %v", attr.Code, attr.Value),
					StatusCode: 400,
				}
			}
		case "object":
			isObject := util.IsObject(attr.Value)
			if !isObject {
				return &helpers.ApiError{
					Message:    fmt.Sprintf("invalid value type for attribute %s: expected object, got %v", attr.Code, attr.Value),
					StatusCode: 400,
				}
			}
			attribute.Value = attr.Value
		default:
			return &helpers.ApiError{
				Message:    fmt.Sprintf("unknown type %s for attribute %s", attr.Type, attr.Code),
				StatusCode: 400,
			}
		}
	}
	return nil
}
