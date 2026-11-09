package product_enums

type AttributeType string

const (
	TEXT           AttributeType = "text"
	NUMBER         AttributeType = "number"
	EMAIL          AttributeType = "email"
	PASSWORD       AttributeType = "password"
	DATE           AttributeType = "date"
	TIME           AttributeType = "time"
	DATETIME_LOCAL AttributeType = "datetime_local"
	RADIO          AttributeType = "radio"
	SWITCH         AttributeType = "switch"
	CHECKBOX       AttributeType = "checkbox"
	FILE           AttributeType = "file"
	COLOR          AttributeType = "color"
	RANGE          AttributeType = "range"
	TEL            AttributeType = "tel"
	URL            AttributeType = "url"
	SEARCH         AttributeType = "search"
	SELECT         AttributeType = "select"
	TEXTAREA       AttributeType = "textarea"
	WYSIWYG_EDITOR AttributeType = "wysiwyg_editor"
	CUSTOM_GROUP   AttributeType = "custom"
)

type AttributeTag string

const (
	SYSTEM AttributeTag = "system"
	CUSTOM AttributeTag = "custom"
)

type BelongTo string

const (
	PRODUCT         BelongTo = "product"
	QUOTATION       BelongTo = "quotation"
	PRODUCTION_LINE BelongTo = "production_line"
)

type AttributeRefField string

const (
	BRANCH AttributeRefField = "branch"
	ADDON  AttributeRefField = "addon"
)

type AttributeCode string

const (
	SYSTEM_BRANCH AttributeCode = "system_branch"
	SYSTEM_ADDON  AttributeCode = "system_addon"
)

type AttributeGroupCode string

const (
	EXTRA AttributeGroupCode = "extra"
)
