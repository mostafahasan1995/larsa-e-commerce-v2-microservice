package config_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	indexKeyCreatedAt  = "createdAt"
	indexKeyUpdatedAt  = "updatedAt"
	indexKeyDeletedAt  = "deletedAt"
	indexKeyUserId     = "userId"
	indexKeyCartId     = "cartId"
	indexKeyBranchId   = "branchId"
	indexKeyProductId  = "productId"
	indexKeyCode       = "code"
	indexKeyBelongsTo  = "belongsTo"
	indexKeyCouponCode = "couponCode"
	indexKeyIsActive   = "isActive"
)

type GetCollectionIndexesResult struct {
	ToBeInserted mongo.IndexModel
	Inserted     string
}

func GetCollectionIndexes(collectionName string) []GetCollectionIndexesResult {
	createdAtIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: indexKeyCreatedAt, Value: 1}},
		Options: options.Index().SetName(indexKeyCreatedAt + "_1"),
	}

	createdAtIndex := GetCollectionIndexesResult{
		ToBeInserted: createdAtIndexModel,
		Inserted:     indexKeyCreatedAt + "_1",
	}

	// Add compound index for _id and deletedAt for collections that use this pattern
	var idDeletedAtIndex = bson.D{{Key: "_id", Value: 1}, {Key: "deletedAt", Value: 1}}
	var idDeletedAtIndexName = "_id_1_deletedAt_1"

	switch collectionName {
	case enums.GROUP_MODEL_NAME:
		var codeIndex = indexKeyCode + "_1"
		return []GetCollectionIndexesResult{
			createdAtIndex,
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    bson.D{{Key: indexKeyCode, Value: 1}},
					Options: options.Index().SetUnique(true).SetName(codeIndex),
				},
				Inserted: codeIndex,
			},
		}
	case enums.PRODUCT_MODEL_NAME:
		return []GetCollectionIndexesResult{
			createdAtIndex,
		}
	case enums.REVIEW_MODEL_NAME:
		var userIdProductIdDeletedAtIndex = indexKeyUserId + "_1_" + indexKeyProductId + "_1_" + indexKeyDeletedAt + "_1"
		return []GetCollectionIndexesResult{
			createdAtIndex,
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    bson.D{{Key: indexKeyUserId, Value: 1}, {Key: indexKeyProductId, Value: 1}, {Key: indexKeyDeletedAt, Value: 1}},
					Options: options.Index().SetUnique(true).SetName(userIdProductIdDeletedAtIndex),
				},
				Inserted: userIdProductIdDeletedAtIndex,
			},
		}
	case enums.BRANCH_PRODUCT_MODEL_NAME:
		var branchIdProductIdIndex = indexKeyBranchId + "_1_" + indexKeyProductId + "_1_" + indexKeyDeletedAt + "_1"
		return []GetCollectionIndexesResult{
			createdAtIndex,
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    bson.D{{Key: indexKeyBranchId, Value: 1}, {Key: indexKeyProductId, Value: 1}, {Key: indexKeyDeletedAt, Value: 1}},
					Options: options.Index().SetUnique(true).SetName(branchIdProductIdIndex),
				},
				Inserted: branchIdProductIdIndex,
			},
		}
	case enums.CART_MODEL_NAME:
		var userIdIsActiveDeletedAtIndex = indexKeyUserId + "_1_" + indexKeyIsActive + "_1_" + indexKeyDeletedAt + "_1"
		return []GetCollectionIndexesResult{
			createdAtIndex,
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    idDeletedAtIndex,
					Options: options.Index().SetName(idDeletedAtIndexName),
				},
				Inserted: idDeletedAtIndexName,
			},
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    bson.D{{Key: indexKeyUserId, Value: 1}, {Key: indexKeyIsActive, Value: 1}, {Key: indexKeyDeletedAt, Value: 1}},
					Options: options.Index().SetName(userIdIsActiveDeletedAtIndex),
				},
				Inserted: userIdIsActiveDeletedAtIndex,
			},
		}
	case enums.CART_ITEM_MODEL_NAME:
		var cartItUserIdIsActiveDeletedAtIndex = indexKeyCartId + "_1_" + indexKeyUserId + "_1_" + indexKeyIsActive + "_1_" + indexKeyDeletedAt + "_1"
		return []GetCollectionIndexesResult{
			createdAtIndex,
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    idDeletedAtIndex,
					Options: options.Index().SetName(idDeletedAtIndexName),
				},
				Inserted: idDeletedAtIndexName,
			},
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    bson.D{{Key: indexKeyCartId, Value: 1}, {Key: indexKeyUserId, Value: 1}, {Key: indexKeyIsActive, Value: 1}, {Key: indexKeyDeletedAt, Value: 1}},
					Options: options.Index().SetName(cartItUserIdIsActiveDeletedAtIndex),
				},
				Inserted: cartItUserIdIsActiveDeletedAtIndex,
			},
		}
	// TODO: add coupon index
	// case enums.COUPON_MODEL_NAME:
	// 	var codeIndex = indexKeyCode + "_1"
	// 	return []GetCollectionIndexesResult{
	// 		createdAtIndex,
	// 		{
	// 			ToBeInserted: mongo.IndexModel{
	// 				Keys:    bson.D{{Key: indexKeyCode, Value: 1}},
	// 				Options: options.Index().SetUnique(true).SetName(codeIndex),
	// 			},
	// 			Inserted: codeIndex,
	// 		},
	// 	}

	case enums.CURRENCY_MODEL_NAME:
		var codeIndex = indexKeyCode + "_1"
		return []GetCollectionIndexesResult{
			createdAtIndex,
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    bson.D{{Key: indexKeyCode, Value: 1}},
					Options: options.Index().SetUnique(true).SetName(codeIndex),
				},
				Inserted: codeIndex,
			},
		}
	case enums.USER_COUPON_MODEL_NAME:
		var userIdCodeIndex = indexKeyUserId + "_1_" + indexKeyCouponCode + "_1"
		return []GetCollectionIndexesResult{
			createdAtIndex,
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    bson.D{{Key: indexKeyUserId, Value: 1}, {Key: indexKeyCouponCode, Value: 1}},
					Options: options.Index().SetUnique(true).SetName(userIdCodeIndex),
				},
				Inserted: userIdCodeIndex,
			},
		}
	case enums.PROFILE_MODEL_NAME:
		return []GetCollectionIndexesResult{
			createdAtIndex,
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    idDeletedAtIndex,
					Options: options.Index().SetName(idDeletedAtIndexName),
				},
				Inserted: idDeletedAtIndexName,
			},
		}

	case enums.AREA_MODEL_NAME:
		return []GetCollectionIndexesResult{
			createdAtIndex,
		}
	case enums.CITY_MODEL_NAME:
		return []GetCollectionIndexesResult{
			createdAtIndex,
		}
	case enums.COUNTRY_MODEL_NAME:
		return []GetCollectionIndexesResult{
			createdAtIndex,
		}
	case enums.UNSERVED_DELIVERY_REQUEST_MODEL_NAME:
		var countryCityAreaCreatedAtIndex = "country_1_city_1_area_1_createdAt_1"
		return []GetCollectionIndexesResult{
			createdAtIndex,
			{
				ToBeInserted: mongo.IndexModel{
					Keys:    bson.D{{Key: "country", Value: 1}, {Key: "city", Value: 1}, {Key: "area", Value: 1}, {Key: indexKeyCreatedAt, Value: 1}},
					Options: options.Index().SetName(countryCityAreaCreatedAtIndex),
				},
				Inserted: countryCityAreaCreatedAtIndex,
			},
		}
	}

	return []GetCollectionIndexesResult{
		createdAtIndex,
	}
}
