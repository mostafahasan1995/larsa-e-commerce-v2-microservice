package product_dtos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExtractProductAttributeKeysResult struct {
	BranchIds  []primitive.ObjectID
	AddonIds   []primitive.ObjectID
	Attributes []product_models.AttributeResponse
}

func ExtractProductAttributeKeys(ctx context.Context, injector *do.Injector, candidateAttributes []product_models.ProductAttribute) (*ExtractProductAttributeKeysResult, *helpers.ApiError) {
	db := do.MustInvoke[*mongo.Client](injector)
	attributeRepo := base_repos.NewBaseRepo[product_models.Attribute, product_models.AttributeResponse](db, enums.ATTRIBUTE)

	attributeCodes := make([]string, 0)
	for _, attr := range candidateAttributes {
		attributeCodes = append(attributeCodes, attr.Code)
	}

	uniqueAttributeCodes := util.ExtractUniqueValues(attributeCodes)

	if len(uniqueAttributeCodes) != len(attributeCodes) {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "some attributes are not unique",
		}
	}

	matchStage := bson.M{
		"code": bson.M{"$in": attributeCodes},
		// "refField":  bson.M{"$in": []product_enums.AttributeRefField{product_enums.BRANCH, product_enums.ADDON}},
		"belongsTo": bson.M{"$in": []product_enums.BelongTo{product_enums.PRODUCT}},
	}
	matchStage = base_utils.AddSoftDeletedFilter(matchStage)
	filter := bson.M{"$match": matchStage}
	storedAttributes, _ := attributeRepo.GetAll(ctx, []bson.M{filter})

	if len(storedAttributes) != len(attributeCodes) {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "some attributes not found",
		}
	}

	result := ExtractProductAttributeKeysResult{}
	result.Attributes = storedAttributes

	for _, storedAttribute := range storedAttributes {

		candidateAttribute := util.Find(candidateAttributes, func(attr product_models.ProductAttribute) bool {
			return attr.Code == storedAttribute.Code
		})

		if storedAttribute.RefField == nil {
			continue
		}

		if *storedAttribute.RefField != product_enums.BRANCH && *storedAttribute.RefField != product_enums.ADDON {
			continue
		}

		switch product_enums.AttributeCode(storedAttribute.Code) {
		case product_enums.SYSTEM_BRANCH:
			ids, err := ExtractValueIds(*candidateAttribute)
			if err != nil {
				return nil, err
			}
			result.BranchIds = append(result.BranchIds, ids...)
		case product_enums.SYSTEM_ADDON:
			ids, err := ExtractValueIds(*candidateAttribute)
			if err != nil {
				return nil, err
			}
			result.AddonIds = append(result.AddonIds, ids...)
		}

	}

	// for _, candidateAttribute := range candidateAttributes {
	// 	switch candidateAttribute.Code {
	// 	case string(product_enums.SYSTEM_BRANCH):
	// 		candidateAttribute.Type = product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY
	// 		candidateAttribute.Value = result.BranchIds
	// 	case string(product_enums.SYSTEM_ADDON):
	// 		candidateAttribute.Type = product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY
	// 		candidateAttribute.Value = result.AddonIds
	// 	}

	// }

	return &result, nil
}

func ExtractValueIds(attribute product_models.ProductAttribute) ([]primitive.ObjectID, *helpers.ApiError) {
	ids := make([]primitive.ObjectID, 0)
	switch attribute.Type {
	case product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY:
		switch v := attribute.Value.(type) {
		case primitive.A:
			for _, v := range v {
				id, err := convertToIds(v)
				if err != nil {
					return nil, err
				}
				ids = append(ids, *id)
			}
		case []interface{}:
			for _, v := range v {
				id, err := convertToIds(v)
				if err != nil {
					return nil, err
				}
				ids = append(ids, *id)
			}
		}

	case product_enums.PRODUCT_ATTRIBUTE_TYPE_STRING:
		id, err := convertToIds(attribute.Value)
		if err != nil {
			return nil, err
		}
		ids = append(ids, *id)
	}

	return util.ExtractUniqueValues(ids), nil
}

func convertToIds(value interface{}) (*primitive.ObjectID, *helpers.ApiError) {

	switch v := value.(type) {
	case primitive.ObjectID:
		return &v, nil
	case string:
		objID, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid object id",
			}
		}
		return &objID, nil

	}

	return nil, nil
}
