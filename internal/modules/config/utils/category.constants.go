package config_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CategoryCalculateProductAndChildrenCount = []bson.M{
	{
		"$lookup": bson.M{
			"from":         enums.PRODUCT.GetCollectionName(),
			"localField":   "_id",
			"foreignField": "categoryIds",
			"as":           "products",
			"pipeline": []bson.M{
				{
					"$match": base_utils.AddSoftDeletedFilter(bson.M{}),
				},
				{
					"$match": bson.M{
						"productType": bson.M{"$ne": product_enums.PRODUCT_TYPE_INVENTORY},
					},
				},
			},
		},
	},
	{
		"$addFields": bson.M{
			"productsCount": bson.M{"$size": "$products"},
		},
	},
	{
		"$lookup": bson.M{
			"from":         enums.CATEGORY.GetCollectionName(),
			"localField":   "_id",
			"foreignField": "parentId",
			"as":           "children",
			"pipeline": []bson.M{
				{
					"$match": base_utils.AddSoftDeletedFilter(bson.M{}),
				},
			},
		},
	},
	{
		"$addFields": bson.M{
			"childrenCount": bson.M{"$size": "$children"},
		},
	},
	{
		"$project": bson.M{
			"products": 0,
			"children": 0,
		},
	},
}

// CategoryParentsPipe is a MongoDB aggregation pipeline that retrieves all parent categories
// for a given category in a hierarchical structure. The pipeline:
// 1. Uses $graphLookup to traverse the category hierarchy upwards from the current category
//   - Starts with the current category's parentId
//   - Connects parent categories through parentId -> _id relationships
//   - Tracks the depth/level of each parent in the hierarchy
//
// 2. Unwinds the parents array to process each parent individually
// 3. Sorts parents by level in descending order (root parent first)
// 4. Groups the results back together, preserving the original category and its parents
// 5. Restructures the output to merge the category with its parents array
//
// The resulting documents will have the original category fields plus a "parents" array
// containing all ancestor categories, ordered from root to immediate parent.
//
// Example:
// Given the following category hierarchy:
// Electronics (root)
// └── Computers
//
//	└── Laptops
//	    └── Gaming Laptops (current category)
//
// Pipeline Steps:
// 1. $graphLookup stage:
//   - Input: Gaming Laptops category
//   - Output: Adds a "parents" array with all ancestor categories and their levels
//     {
//     "_id": "gaming-laptops-id",
//     "name": "Gaming Laptops",
//     "parentId": "laptops-id",
//     "parents": [
//     { "_id": "electronics-id", "name": "Electronics", "parentId": null, "level": 0 },
//     { "_id": "computers-id", "name": "Computers", "parentId": "electronics-id", "level": 1 },
//     { "_id": "laptops-id", "name": "Laptops", "parentId": "computers-id", "level": 2 }
//     ]
//     }
//
// 2. $unwind stage:
//   - Input: Document with parents array
//   - Output: One document per parent, preserving the original category
//     [
//     { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id", "parents": { "_id": "electronics-id", "level": 0 } },
//     { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id", "parents": { "_id": "computers-id", "level": 1 } },
//     { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id", "parents": { "_id": "laptops-id", "level": 2 } }
//     ]
//
// 3. $sort stage:
//   - Input: Unwound documents
//   - Output: Same documents sorted by parents.level in descending order
//     [
//     { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id", "parents": { "_id": "electronics-id", "level": 0 } },
//     { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id", "parents": { "_id": "computers-id", "level": 1 } },
//     { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id", "parents": { "_id": "laptops-id", "level": 2 } }
//     ]
//
// 4. $group stage:
//   - Input: Sorted documents
//   - Output: Groups by _id, preserves original category, and collects parents
//     {
//     "_id": "gaming-laptops-id",
//     "category": { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id" },
//     "parents": [
//     { "_id": "electronics-id", "name": "Electronics", "parentId": null, "level": 0 },
//     { "_id": "computers-id", "name": "Computers", "parentId": "electronics-id", "level": 1 },
//     { "_id": "laptops-id", "name": "Laptops", "parentId": "computers-id", "level": 2 }
//     ]
//     }
//
// 5. $replaceRoot stage:
//   - Input: Grouped document
//   - Output: Merges category and parents into a single document
//     {
//     "_id": "gaming-laptops-id",
//     "name": "Gaming Laptops",
//     "parentId": "laptops-id",
//     "parents": [
//     { "_id": "electronics-id", "name": "Electronics", "parentId": null, "level": 0 },
//     { "_id": "computers-id", "name": "Computers", "parentId": "electronics-id", "level": 1 },
//     { "_id": "laptops-id", "name": "Laptops", "parentId": "computers-id", "level": 2 }
//     ]
//     }
var CategoryParentsPipe = []bson.M{
	{
		"$graphLookup": bson.M{
			"from":             enums.CATEGORY.GetCollectionName(),
			"startWith":        "$parentId",
			"connectFromField": "parentId",
			"connectToField":   "_id",
			"as":               "parents",
			"depthField":       "level",
			"restrictSearchWithMatch": bson.M{
				"deletedAt": bson.M{"$exists": false},
			},
		},
	},
	{
		"$unwind": bson.M{
			"path":                       "$parents",
			"preserveNullAndEmptyArrays": true,
		},
	},
	{
		"$sort": bson.M{
			"parents.level": -1,
		},
	},
	{
		"$group": bson.M{
			"_id":      "$_id",
			"category": bson.M{"$first": "$$ROOT"},
			"parents":  bson.M{"$push": "$parents"},
		},
	},
	{
		"$replaceRoot": bson.M{
			"newRoot": bson.M{
				"$mergeObjects": bson.A{
					"$category",
					bson.M{"parents": "$parents"},
				},
			},
		},
	},
}

// CategoryChildrenPipe is a MongoDB aggregation pipeline that retrieves all child categories
// for a given category in a hierarchical structure. The pipeline:
// 1. Uses $graphLookup to traverse the category hierarchy downwards from the current category
//   - Starts with the current category's _id
//   - Connects child categories through _id -> parentId relationships
//   - Tracks the depth/level of each child in the hierarchy
//
// 2. Unwinds the children array to process each child individually
// 3. Sorts children by level in ascending order (immediate children first)
// 4. Groups the results back together, preserving the original category and its children
// 5. Restructures the output to merge the category with its children array
//
// The resulting documents will have the original category fields plus a "children" array
// containing all descendant categories, ordered from immediate children to deepest descendants.
//
// Example:
// Given the following category hierarchy:
// Electronics (current category)
// └── Computers
//
//	└── Laptops
//	    └── Gaming Laptops
//
// Pipeline Steps:
// 1. $graphLookup stage:
//   - Input: Electronics category
//   - Output: Adds a "children" array with all descendant categories and their levels
//     {
//     "_id": "electronics-id",
//     "name": "Electronics",
//     "parentId": null,
//     "children": [
//     { "_id": "computers-id", "name": "Computers", "parentId": "electronics-id", "level": 1 },
//     { "_id": "laptops-id", "name": "Laptops", "parentId": "computers-id", "level": 2 },
//     { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id", "level": 3 }
//     ]
//     }
//
// 2. $unwind stage:
//   - Input: Document with children array
//   - Output: One document per child, preserving the original category
//     [
//     { "_id": "electronics-id", "name": "Electronics", "parentId": null, "children": { "_id": "computers-id", "level": 1 } },
//     { "_id": "electronics-id", "name": "Electronics", "parentId": null, "children": { "_id": "laptops-id", "level": 2 } },
//     { "_id": "electronics-id", "name": "Electronics", "parentId": null, "children": { "_id": "gaming-laptops-id", "level": 3 } }
//     ]
//
// 3. $sort stage:
//   - Input: Unwound documents
//   - Output: Same documents sorted by children.level in ascending order
//     [
//     { "_id": "electronics-id", "name": "Electronics", "parentId": null, "children": { "_id": "computers-id", "level": 1 } },
//     { "_id": "electronics-id", "name": "Electronics", "parentId": null, "children": { "_id": "laptops-id", "level": 2 } },
//     { "_id": "electronics-id", "name": "Electronics", "parentId": null, "children": { "_id": "gaming-laptops-id", "level": 3 } }
//     ]
//
// 4. $group stage:
//   - Input: Sorted documents
//   - Output: Groups by _id, preserves original category, and collects children
//     {
//     "_id": "electronics-id",
//     "category": { "_id": "electronics-id", "name": "Electronics", "parentId": null },
//     "children": [
//     { "_id": "computers-id", "name": "Computers", "parentId": "electronics-id", "level": 1 },
//     { "_id": "laptops-id", "name": "Laptops", "parentId": "computers-id", "level": 2 },
//     { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id", "level": 3 }
//     ]
//     }
//
// 5. $replaceRoot stage:
//   - Input: Grouped document
//   - Output: Merges category and children into a single document
//     {
//     "_id": "electronics-id",
//     "name": "Electronics",
//     "parentId": null,
//     "children": [
//     { "_id": "computers-id", "name": "Computers", "parentId": "electronics-id", "level": 1 },
//     { "_id": "laptops-id", "name": "Laptops", "parentId": "computers-id", "level": 2 },
//     { "_id": "gaming-laptops-id", "name": "Gaming Laptops", "parentId": "laptops-id", "level": 3 }
//     ]
//     }
var CategoryChildrenPipe = []bson.M{
	{
		"$graphLookup": bson.M{
			"from":             enums.CATEGORY.GetCollectionName(),
			"startWith":        "$_id",
			"connectFromField": "_id",
			"connectToField":   "parentId",
			"as":               "children",
			"depthField":       "level",
			"restrictSearchWithMatch": bson.M{
				"deletedAt": bson.M{"$exists": false},
			},
		},
	},
	{
		"$unwind": bson.M{
			"path":                       "$children",
			"preserveNullAndEmptyArrays": true,
		},
	},
	{
		"$sort": bson.M{
			"children.level": 1,
		},
	},
	{
		"$group": bson.M{
			"_id":      "$_id",
			"category": bson.M{"$first": "$$ROOT"},
			"children": bson.M{"$push": "$children"},
		},
	},
	{
		"$replaceRoot": bson.M{
			"newRoot": bson.M{
				"$mergeObjects": bson.A{
					"$category",
					bson.M{"children": "$children"},
				},
			},
		},
	},
}

// GetCategoryChildrenIdsPipeline returns a pipeline that gets all child category IDs for given category IDs
// This is used in product filtering to include products from child categories
func GetCategoryChildrenIdsPipeline(categoryIds []primitive.ObjectID) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"_id": bson.M{"$in": categoryIds},
			},
		},
		{
			"$graphLookup": bson.M{
				"from":             enums.CATEGORY.GetCollectionName(),
				"startWith":        "$_id",
				"connectFromField": "_id",
				"connectToField":   "parentId",
				"as":               "children",
				"depthField":       "level",
				"restrictSearchWithMatch": bson.M{
					"deletedAt": bson.M{"$exists": false},
				},
			},
		},
		{
			"$addFields": bson.M{
				"allCategoryIds": bson.M{
					"$concatArrays": bson.A{
						"$_id",
						bson.M{
							"$map": bson.M{
								"input": "$children",
								"as":    "child",
								"in":    "$$child._id",
							},
						},
					},
				},
			},
		},
		{
			"$unwind": "$allCategoryIds",
		},
		{
			"$group": bson.M{
				"_id": nil,
				"categoryIds": bson.M{
					"$addToSet": "$allCategoryIds",
				},
			},
		},
	}
}
