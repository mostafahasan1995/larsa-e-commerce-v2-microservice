package base_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

// RefField represents a reference field configuration
type RefField struct {
	Let              *bson.M
	Collection       string
	IsNeedAuth       bool
	IsArray          bool
	ForeignField     string
	LocalField       string
	As               string
	Pipeline         []bson.M
	ExternalPipeline []bson.M
}

func projectToRemoveMetaData(keys []string) bson.M {
	project := bson.M{"createdAt": 0, "updatedAt": 0, "createdBy": 0, "updatedBy": 0, "_id": 0}
	for _, key := range keys {
		project[key] = 0
	}
	return bson.M{"$project": project}
}

// GetRefFields returns the reference fields configuration for a given module
func GetRefField(moduleName enums.CollectionName) map[string]RefField {
	switch moduleName {
	case enums.CATEGORY:
		return map[string]RefField{
			"parentId": {
				Collection:   enums.CATEGORY.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "parentId",
				As:           "parent",
			},
		}
	case enums.ATTRIBUTE:
		return map[string]RefField{
			"groupId": {
				Collection:   enums.GROUP.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "groupId",
				As:           "group",
			},
		}
	case enums.CURRENCY:
		return map[string]RefField{
			"countryIso3": {
				Collection:   enums.COUNTRY.GetCollectionName(),
				IsArray:      false,
				ForeignField: "code.iso3",
				LocalField:   "countryIso3",
				As:           "country",
			},
		}
	case enums.COUNTRY:
		return map[string]RefField{
			"currencyCodes": {
				Collection:   enums.CURRENCY.GetCollectionName(),
				IsArray:      true,
				ForeignField: "code",
				LocalField:   "currencyCodes",
				As:           "currencies",
			},
		}
	case enums.CITY:
		return map[string]RefField{
			"countryId": {
				Collection:   enums.COUNTRY.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "countryId",
				As:           "country",
				Pipeline: []bson.M{
					{"$project": bson.M{"multiPolygon": 0}},
				},
			},
			"branches": {
				Collection:   enums.BRANCH.GetCollectionName(),
				IsArray:      true,
				ForeignField: "city",
				LocalField:   "name",
				As:           "branches",
				Pipeline: []bson.M{
					{
						"$match": AddSoftDeletedFilter(bson.M{
							"isActive": true,
						}),
					},
				},
				ExternalPipeline: []bson.M{
					{
						"$addFields": bson.M{
							"branches": bson.M{"$ifNull": []interface{}{"$branches", []interface{}{}}},
						},
					},
				},
			},
		}
	case enums.AREA:
		return map[string]RefField{
			"cityId": {
				Collection:   enums.CITY.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "cityId",
				As:           "city",
				Pipeline: []bson.M{
					{"$project": bson.M{"multiPolygon": 0}},
				},
			},
		}

	case enums.BRANCH:
		return map[string]RefField{
			"products": {
				Let:          &bson.M{"showProductOnIsActive": "$showProductOnIsActive"},
				Collection:   enums.BRANCH_PRODUCT.GetCollectionName(),
				IsArray:      true,
				LocalField:   "_id",
				ForeignField: "branchId",
				As:           "products",
				Pipeline: []bson.M{{"$match": bson.M{
					"$or": []bson.M{
						{"$expr": bson.M{"$eq": []interface{}{bson.M{"$ifNull": []interface{}{"$$showProductOnIsActive", nil}}, nil}}}, // If isActive is null, include all items
						{"$expr": bson.M{"$eq": []interface{}{"$isActive", "$$showProductOnIsActive"}}},                                // If isActive has value, match it
					},
				}}, {"$match": AddSoftDeletedFilter(bson.M{})}, projectToRemoveMetaData([]string{"branchId"})},
			},
		}
	case enums.PRODUCT:
		return map[string]RefField{
			"unitIds": {
				Collection:   enums.UNIT.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "unitIds",
				As:           "units",
			},
			"inventoryProductIds": {
				Collection:   enums.PRODUCT.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "inventoryProductIds",
				As:           "inventoryProducts",
			},
			"productionLineId": {
				Collection:   enums.PRODUCTION_LINE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "productionLineId",
				As:           "productionLine",
			},
			"branchIds": {
				Collection:   enums.BRANCH.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "branchIds",
				As:           "branches",
			},
			"categoryIds": {
				Collection:   enums.CATEGORY.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "categoryIds",
				As:           "categories",
			},
			"attributes.value": {
				Collection:   enums.BRANCH_PRODUCT.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "attributes.value",
				As:           "values",
			},
			"isInWatchListed": {
				Collection:   enums.WATCH_LIST.GetCollectionName(),
				IsArray:      true,
				IsNeedAuth:   true,
				ForeignField: "itemId",
				LocalField:   "_id",
				Let:          &bson.M{"userId": "$userId"},
				As:           "watchlist",
				Pipeline: []bson.M{
					{
						"$match": bson.M{
							"itemType":  "PRODUCT",
							"$expr":     bson.M{"$eq": []interface{}{"$userId", "$$userId"}},
							"deletedAt": bson.M{"$exists": false},
						},
					},
				},
				ExternalPipeline: []bson.M{
					{
						"$addFields": bson.M{
							"isInWatchListed": bson.M{
								"$gt": bson.A{bson.M{"$size": "$watchlist"}, 0},
							},
						},
					},
				},
			},

			"reviews": {
				Collection:   enums.REVIEW.GetCollectionName(),
				IsArray:      true,
				ForeignField: "productId",
				LocalField:   "_id",
				As:           "reviews",
				Pipeline: []bson.M{
					{"$match": AddSoftDeletedFilter(bson.M{})},
				},
			},
			"reviewStats": {
				Collection:   enums.REVIEW.GetCollectionName(),
				IsArray:      false,
				ForeignField: "productId",
				LocalField:   "_id",
				As:           "reviewStats",
				Pipeline: []bson.M{
					{"$match": AddSoftDeletedFilter(bson.M{})},
					{
						"$group": bson.M{
							"_id": "$productId",
							"averageRating": bson.M{
								"$avg": "$rate",
							},
							"totalReviews": bson.M{
								"$sum": 1,
							},
						},
					},
					{
						"$project": bson.M{
							"_id": 0,
							"averageRating": bson.M{
								"$round": bson.A{"$averageRating", 2},
							},
							"totalReviews": 1,
						},
					},
				},
				ExternalPipeline: []bson.M{
					{
						"$addFields": bson.M{
							"averageRating": bson.M{
								"$ifNull": bson.A{
									bson.M{"$arrayElemAt": bson.A{"$reviewStats.averageRating", 0}},
									nil,
								},
							},
							"totalReviews": bson.M{
								"$ifNull": bson.A{
									bson.M{"$arrayElemAt": bson.A{"$reviewStats.totalReviews", 0}},
									nil,
								},
							},
						},
					},
				},
			},
		}
	case enums.PRODUCTION_LINE:
		return map[string]RefField{
			"productIds": {
				Collection:   enums.PRODUCT.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "productIds",
				As:           "products",
			},
			"categoryIds": {
				Collection:   enums.CATEGORY.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "categoryIds",
				As:           "categories",
			},
		}
	case enums.UNIT:
		return map[string]RefField{
			"categoryId": {
				Collection:   enums.CATEGORY.GetCollectionName(),
				IsArray:      false,
				LocalField:   "categoryId",
				ForeignField: "_id",
				As:           "category",
			},
		}
	case enums.QUOTATION:
		return map[string]RefField{
			"vendorProfileId": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "vendorProfileId",
				As:           "vendorProfile",
			},
			"quotationBuyerInfo.buyerProfileId": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "quotationBuyerInfo.buyerProfileId",
				As:           "quotationBuyerInfo.buyerProfile",
			},
			//NOTE: I remove the populate products because it doesn't work due to nested array populating. So I add it manually in the quotation.constants.go
			// "products.productId": {
			// 	Collection:   enums.PRODUCT.GetCollectionName(),
			// 	IsArray:      false,
			// 	ForeignField: "_id",
			// 	LocalField:   "products.productId",
			// 	As:           "products.product",
			// },
		}
	case enums.INVOICE:
		return map[string]RefField{
			"transactionId": {
				Collection:   enums.TRANSACTION.GetCollectionName(),
				IsArray:      false,
				LocalField:   "transactionId",
				ForeignField: "_id",
				As:           "transaction",
			},
		}
	case enums.PAYMENT:
		return map[string]RefField{
			"invoiceId": {
				Collection:   enums.INVOICE.GetCollectionName(),
				IsArray:      false,
				LocalField:   "invoiceId",
				ForeignField: "_id",
				As:           "invoice",
			},
		}
	case enums.CUSTOMER:
		return map[string]RefField{
			"transactions": {
				Collection:   enums.TRANSACTION.GetCollectionName(),
				IsArray:      true,
				LocalField:   "transactions",
				ForeignField: "_id",
				As:           "transactions",
			},
			"invoices": {
				Collection:   enums.INVOICE.GetCollectionName(),
				IsArray:      true,
				LocalField:   "invoices",
				ForeignField: "_id",
				As:           "invoices",
			},
			"payments": {
				Collection:   enums.PAYMENT.GetCollectionName(),
				IsArray:      true,
				LocalField:   "payments",
				ForeignField: "_id",
				As:           "payments",
			},
		}
	case enums.FEATURE:
		return map[string]RefField{
			"groupId": {
				Collection:   enums.FEATURE_GROUP.GetCollectionName(),
				IsArray:      false,
				LocalField:   "groupId",
				ForeignField: "_id",
				As:           "group",
			},
		}
	case enums.FEATURE_GROUP:
		return map[string]RefField{
			"featureIds": {
				Collection:   enums.FEATURE.GetCollectionName(),
				IsArray:      true,
				LocalField:   "featureIds",
				ForeignField: "_id",
				As:           "features",
			},
		}
	case enums.PROMOTION:
		return map[string]RefField{
			"branchIds": {
				Collection:   enums.BRANCH.GetCollectionName(),
				IsArray:      true,
				LocalField:   "branchIds",
				ForeignField: "_id",
				As:           "branches",
			},
			"productIds": {
				Collection:   enums.PRODUCT.GetCollectionName(),
				IsArray:      true,
				LocalField:   "productIds",
				ForeignField: "_id",
				As:           "products",
			},
		}
	case enums.APPOINTMENT:
		return map[string]RefField{
			"vendorId": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				LocalField:   "vendorId",
				ForeignField: "_id",
				As:           "vendor",
			},
			"buyerId": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				LocalField:   "buyerId",
				ForeignField: "_id",
				As:           "buyer",
			},
			"translation.translateFromId": {
				//TODO: Fix from
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				LocalField:   "translation.translateFromId",
				ForeignField: "_id",
				As:           "translation.translateFrom",
			},
			"translation.translateToId": {
				//TODO: Fix from
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				LocalField:   "translation.translateToId",
				ForeignField: "_id",
				As:           "translation.translateTo",
			},
		}
	case enums.ORDER:
		return map[string]RefField{
			"quotationId": {
				Collection:   enums.QUOTATION.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "quotationId",
				As:           "quotation",
			},
			"buyerId": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "buyerInfo.id",
				As:           "buyerInfo.buyer",
			},
			"vendorId": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "vendorInfo.id",
				As:           "vendorInfo.vendor",
			},
			"pickupInfo.pickupBranchId": {
				Collection:   enums.BRANCH.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "pickupInfo.pickupBranchId",
				As:           "pickupInfo.pickupBranch",
			},
			"deliveryInfo.b2cInfo.branchId": {
				Collection:   enums.BRANCH.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "deliveryInfo.b2cInfo.branchId",
				As:           "deliveryInfo.b2cInfo.branch",
			},
		}

	case enums.COUPON:
		return map[string]RefField{
			"userCoupon": {
				Let:          &bson.M{"userId": "$userId"},
				Collection:   enums.USER_COUPON.GetCollectionName(),
				IsArray:      false,
				LocalField:   "code",
				ForeignField: "couponCode",
				As:           "userCoupon",
				Pipeline: []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": []interface{}{"$userId", "$$userId"},
							},
						},
					},
				},
			},
		}
	case enums.CART_ITEM:
		return map[string]RefField{
			/*
			* NOTE: We need to handle all types of items in the cart item.
			 */
			"itemId": {
				Collection:   enums.PRODUCT.GetCollectionName(),
				IsArray:      false,
				LocalField:   "itemId",
				ForeignField: "_id",
				As:           "product",
			},
			"addons": {
				Collection:   enums.ADDON.GetCollectionName(),
				IsArray:      true,
				LocalField:   "addons.id",
				ForeignField: "_id",
				As:           "addonList",
			},
		}
	case enums.BANNER:
		/*
		* NOTE: We need to handle all types of items in the cart item.
		 */
		return map[string]RefField{
			"item": {
				Collection:   enums.PRODUCT.GetCollectionName(),
				IsArray:      false,
				LocalField:   "item.id",
				ForeignField: "_id",
				As:           "product",
			},
		}
	case enums.CART:
		return map[string]RefField{
			"cartItemIds": {
				Let:          &bson.M{"vendorId": "$vendorId", "branchId": "$branchId"},
				Collection:   enums.CART_ITEM.GetCollectionName(),
				IsArray:      true,
				LocalField:   "cartItemIds",
				ForeignField: "_id",
				As:           "cartItems",
				Pipeline: []bson.M{
					{
						"$match": bson.M{
							"$and": []bson.M{
								{
									"isActive": true,
								},
								{
									"$or": []bson.M{
										{"$expr": bson.M{"$eq": []interface{}{bson.M{"$ifNull": []interface{}{"$$vendorId", nil}}, nil}}}, // If vendorId is null, include all items
										{"$expr": bson.M{"$eq": []interface{}{"$vendorId", "$$vendorId"}}},                                // If vendorId has value, match it
									},
								},
								{
									"$or": []bson.M{
										{"$expr": bson.M{"$eq": []interface{}{bson.M{"$ifNull": []interface{}{"$$branchId", nil}}, nil}}}, // If branchId is null, include all items
										{"$expr": bson.M{"$eq": []interface{}{"$branchId", "$$branchId"}}},                                // If branchId has value, match it
									},
								},
							},
						},
					},
					{
						"$lookup": bson.M{
							"from":         enums.ADDON.GetCollectionName(),
							"localField":   "addons.id",
							"foreignField": "_id",
							"as":           "addonList",
						},
					},
					{
						"$lookup": bson.M{
							"from":         enums.PRODUCT.GetCollectionName(),
							"localField":   "itemId",
							"foreignField": "_id",
							"as":           "product",
							"pipeline": []bson.M{
								{
									"$match": AddSoftDeletedFilter(bson.M{}),
								},
							},
						},
					},
					{
						"$unwind": bson.M{
							"path":                       "$product",
							"preserveNullAndEmptyArrays": true,
						},
					},
				},
				ExternalPipeline: []bson.M{
					{
						"$addFields": bson.M{
							"cartItemIds": bson.M{
								"$map": bson.M{
									"input": "$cartItems",
									"as":    "item",
									"in":    "$$item._id",
								},
							},
						},
					},
				},
			},
		}
	case enums.WAITER:
		return map[string]RefField{
			"branchId": {
				Collection:   enums.BRANCH.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "branchId",
				As:           "branch",
			},
		}
	case enums.QUESTION:
		return map[string]RefField{
			"branchIds": {
				Collection:   enums.BRANCH.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "branchIds",
				As:           "branches",
			},
		}
	case enums.FEEDBACK:
		return map[string]RefField{
			"branchId": {
				Collection:   enums.BRANCH.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "branchId",
				As:           "branch",
			},
			"waiterId": {
				Collection:   enums.WAITER.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "waiterId",
				As:           "waiter",
			},
			"answers": {
				Collection:   enums.QUESTION.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "answers.questionId",
				As:           "questions",
			},
		}
	case enums.TICKET_TYPE:
		return map[string]RefField{
			"teamNotifyIds": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "teamNotifyIds",
				As:           "teamNotify",
			},
		}
	case enums.TICKET:
		return map[string]RefField{
			"typeId": {
				Collection:   enums.TICKET_TYPE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "typeId",
				As:           "ticketType",
			},
			"requesterId": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "requesterId",
				As:           "requester",
			},
			"vendorId": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "vendorId",
				As:           "vendor",
			},
			"buyerId": {
				Collection:   enums.PROFILE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "buyerId",
				As:           "buyer",
			},
			"quotationId": {
				Collection:   enums.QUOTATION.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "quotationId",
				As:           "quotation",
			},
			"invoiceId": {
				Collection:   enums.INVOICE.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "invoiceId",
				As:           "invoice",
			},
			"orderId": {
				Collection:   enums.ORDER.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "orderId",
				As:           "order",
			},
		}
	case enums.WATCH_LIST:
		return map[string]RefField{
			"itemId": {
				Collection:   enums.PRODUCT.GetCollectionName(),
				IsArray:      false,
				ForeignField: "_id",
				LocalField:   "itemId",
				As:           "product",
				Pipeline: []bson.M{
					{"$match": AddSoftDeletedFilter(bson.M{})},
				},
			},
		}
	case enums.RECIPE:
		return map[string]RefField{
			"productIds": {
				Collection:   enums.PRODUCT.GetCollectionName(),
				IsArray:      true,
				ForeignField: "_id",
				LocalField:   "productIds",
				As:           "products",
				Pipeline: []bson.M{
					{"$match": AddSoftDeletedFilter(bson.M{})},
				},
			},
		}

	default:
		return nil
	}
}
