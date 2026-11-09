package product_enums

type ProductFlag string

const (
	HOT_PICK     ProductFlag = "hot_pick"
	TRENDING     ProductFlag = "trending"
	INNOVATIVE   ProductFlag = "innovative"
	NEW          ProductFlag = "new"
	MOST_POPULAR ProductFlag = "most_popular"
	LOW_MOQ      ProductFlag = "low_moq"
	OEM          ProductFlag = "oem"
)

type ProductType string

const (
	PRODUCT_TYPE_INVENTORY    ProductType = "inventory"
	PRODUCT_TYPE_SIMPLE       ProductType = "simple"
	PRODUCT_TYPE_CONFIGURABLE ProductType = "configurable"
)

type ProductAttributeType string

const (
	PRODUCT_ATTRIBUTE_TYPE_STRING  ProductAttributeType = "string"
	PRODUCT_ATTRIBUTE_TYPE_NUMBER  ProductAttributeType = "number"
	PRODUCT_ATTRIBUTE_TYPE_BOOLEAN ProductAttributeType = "boolean"
	PRODUCT_ATTRIBUTE_TYPE_ARRAY   ProductAttributeType = "array"
	PRODUCT_ATTRIBUTE_TYPE_OBJECT  ProductAttributeType = "object"
)

type OutOfStock string

const (
	OUT_OF_STOCK_IN_STOCK     OutOfStock = "in_stock"
	OUT_OF_STOCK_OUT_OF_STOCK OutOfStock = "out_of_stock"
	OUT_OF_STOCK_UNKNOWN      OutOfStock = "unknown"
)
