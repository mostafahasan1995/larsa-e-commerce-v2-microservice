package seed_data

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
)

var refFieldBranch = product_enums.AttributeRefField(product_enums.BRANCH)
var refFieldAddon = product_enums.AttributeRefField(product_enums.ADDON)
var Attributes = []product_dtos.CreateAttributeDto{
	{
		Name: types.Localizable[string]{
			"en": "Branch",
			"ar": "الفرع",
		},
		Type:           product_enums.SELECT,
		IsFiltrable:    true,
		Required:       false,
		IsComparable:   true,
		Visibility:     true,
		Sortable:       true,
		Code:           string(product_enums.SYSTEM_BRANCH),
		BelongsTo:      []product_enums.BelongTo{product_enums.PRODUCT},
		RefField:       &refFieldBranch,
		Multi:          true,
		SelectionStyle: &[]string{"dropdown"}[0],
	},
	{
		Name: types.Localizable[string]{
			"en": "Addon",
			"ar": "الإضافي",
		},
		Type:           product_enums.SELECT,
		IsFiltrable:    true,
		Required:       false,
		IsComparable:   true,
		Visibility:     true,
		Sortable:       true,
		Code:           string(product_enums.SYSTEM_ADDON),
		BelongsTo:      []product_enums.BelongTo{product_enums.PRODUCT},
		RefField:       &refFieldAddon,
		Multi:          true,
		SelectionStyle: &[]string{"dropdown"}[0],
	},
}

var groupCode = string(product_enums.EXTRA)
var AttributesGroup = product_dtos.CreateGroupDto{

	Name: types.Localizable[string]{
		"en": "Extra",
		"ar": "إضافي",
	},
	Code: &groupCode,
}
