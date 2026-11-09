package appointment_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

var AppointmentPopulateVendorAndBuyer = []bson.M{
	{
		"$lookup": bson.M{
			"from":         enums.PROFILE.GetCollectionName(),
			"localField":   "vendorId",
			"foreignField": "_id",
			"as":           "vendor",
		},
	},
	{
		"$unwind": bson.M{
			"path":                       "$vendor",
			"preserveNullAndEmptyArrays": true,
		},
	},
	{
		"$lookup": bson.M{
			"from":         enums.PROFILE.GetCollectionName(),
			"localField":   "buyerId",
			"foreignField": "_id",
			"as":           "buyer",
		},
	},
	{
		"$unwind": bson.M{
			"path":                       "$buyer",
			"preserveNullAndEmptyArrays": true,
		},
	},
}

var AppointmentPopulateTranslation = []bson.M{
	{
		"$lookup": bson.M{
			//TODO: Fix from
			"from":         enums.PROFILE.GetCollectionName(),
			"localField":   "translation.translateFromId",
			"foreignField": "_id",
			"as":           "translation.translateFrom",
		},
	},
	{
		//TODO: Fix from
		"$lookup": bson.M{
			"from":         enums.PROFILE.GetCollectionName(),
			"localField":   "translation.translateToId",
			"foreignField": "_id",
			"as":           "translation.translateTo",
		},
	},
}
