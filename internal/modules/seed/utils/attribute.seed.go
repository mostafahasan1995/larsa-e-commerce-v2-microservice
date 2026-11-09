package seed_utils

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	seeds_data "larsa-e-commerce-microservice-v2/internal/modules/seed/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedAttributes seeds the initial attributes into the database
func SeedAttributes(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.ATTRIBUTE.GetCollectionName())
	groupCollection := db.Collection(enums.GROUP.GetCollectionName())
	filter := bson.M{
		"$and": []bson.M{
			{
				"tag": product_enums.SYSTEM,
				"$or": []bson.M{
					{"code": product_enums.SYSTEM_BRANCH},
					{"code": product_enums.SYSTEM_ADDON},
				},
			},
		},
	}
	storedAttributes, findErr := collection.Find(ctx, filter)
	if findErr != nil {
		return helpers.DefaultError(findErr)
	}

	storedAttributesMap := make(map[string]product_models.Attribute)
	for storedAttributes.Next(ctx) {
		var attribute product_models.Attribute
		if err := storedAttributes.Decode(&attribute); err != nil {
			return helpers.DefaultError(err)
		}
		storedAttributesMap[attribute.Code] = attribute
	}

	result := groupCollection.FindOne(ctx, bson.M{"code": product_enums.EXTRA})

	var group *product_models.Group
	var groupId *primitive.ObjectID
	if result.Err() != nil {
		// Document not found, create new group
		groupModel := product_utils.ConvertToGroupModel(&seeds_data.AttributesGroup, types.AdditionalOptions{
			UserId: defaultUserId,
			Mode:   types.CREATE,
		})
		groupId = &groupModel.Id
		_, err := groupCollection.InsertOne(ctx, groupModel)
		if err != nil {
			return helpers.DefaultError(err)
		}
		group = groupModel
	} else {
		// Document found, decode it
		if err := result.Decode(&group); err != nil {
			return helpers.DefaultError(err)
		}
		groupModel := product_utils.ConvertToGroupModel(&seeds_data.AttributesGroup, types.AdditionalOptions{
			UserId: defaultUserId,
			Mode:   types.UPDATE,
		})
		groupId = &group.Id
		_, err := groupCollection.UpdateOne(ctx, bson.M{"_id": group.Id}, bson.M{"$set": groupModel})
		if err != nil {
			return helpers.DefaultError(err)
		}
	}

	// Delete existing attributes that match either code or system tag
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return helpers.DefaultError(err)
	}

	// Convert seed data to attribute models
	attributes := make([]product_models.Attribute, 0, len(seeds_data.Attributes))
	for _, seed := range seeds_data.Attributes {
		attribute := product_utils.ConvertToAttributeModel(&seed, types.AdditionalOptions{
			UserId: defaultUserId,
			Mode:   types.CREATE,
		})
		attribute.Tag = product_enums.SYSTEM
		attributes = append(attributes, *attribute)
	}

	// Convert attributes to []interface{}
	docs := make([]interface{}, len(attributes))
	for i := range attributes {
		storedAttributeMap := storedAttributesMap[attributes[i].Code]
		if storedAttributeMap.Id.IsZero() {
			attributes[i].GroupId = groupId
			docs[i] = attributes[i]
		} else {
			attributes[i].Id = storedAttributeMap.Id
			attributes[i].GroupId = groupId
			attributes[i].CreatedAt = storedAttributeMap.CreatedAt
			attributes[i].CreatedBy = storedAttributeMap.CreatedBy
			attributes[i].BelongsTo = storedAttributeMap.BelongsTo

			uniqueBelongsTo := make(map[product_enums.BelongTo]bool)

			for _, belongTo := range attributes[i].BelongsTo {
				uniqueBelongsTo[belongTo] = true
			}
			for _, belongTo := range storedAttributeMap.BelongsTo {
				uniqueBelongsTo[belongTo] = true
			}
			attributes[i].BelongsTo = []product_enums.BelongTo{}
			for belongTo := range uniqueBelongsTo {
				attributes[i].BelongsTo = append(attributes[i].BelongsTo, belongTo)
			}

			docs[i] = attributes[i]
		}
	}

	_, insertManyErr := collection.InsertMany(ctx, docs)
	if insertManyErr != nil {
		return insertManyErr
	}

	return nil
}
