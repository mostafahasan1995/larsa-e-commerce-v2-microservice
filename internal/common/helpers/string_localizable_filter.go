package helpers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

type StringLocalizableFilter string

func (sf *StringLocalizableFilter) BuildFilter(fieldPath string) bson.M {
	if sf == nil || *sf == "" {
		return nil
	}

	orConditions := []bson.M{}
	for _, lang := range enums.LanguageList {
		langKey := lang.GetLanguageKey()
		fullPath := fieldPath + "." + langKey
		orConditions = append(orConditions, bson.M{
			fullPath: bson.M{"$regex": *sf, "$options": "i"},
		})
	}

	if len(orConditions) == 0 {
		return nil
	}

	return bson.M{"$or": orConditions}
}
