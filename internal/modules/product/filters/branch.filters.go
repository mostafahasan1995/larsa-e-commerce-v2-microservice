package product_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductBranchFilters struct {
	IsActive *bool `json:"-"`
}

type BranchFilters struct {
	base_filters.BaseFilters
	VendorId             *primitive.ObjectID   `json:"vendorId,omitempty"`
	Lat                  float64               `json:"lat,omitempty"`
	Lng                  float64               `json:"lng,omitempty"`
	City                 string                `json:"city,omitempty"`
	ZoneIds              []primitive.ObjectID  `json:"zoneIds,omitempty"`
	ProductBranchFilters *ProductBranchFilters `json:"productBranchFilters,omitempty"`
	ManagerId            *primitive.ObjectID   `json:"managerId,omitempty"`
}

func NewBranchFilters(query string) (*BranchFilters, error) {
	branchFilters := BranchFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &branchFilters); err != nil {
			return nil, err
		}
	}

	return &branchFilters, nil
}

func (branchFilters *BranchFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if branchFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	baseFilters := &branchFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.BRANCH)

	if branchFilters.City != "" {
		matchStage["city"] = branchFilters.City
	}

	if len(branchFilters.ZoneIds) > 0 {
		matchStage["zones"] = bson.M{
			"$elemMatch": bson.M{
				"zoneId": bson.M{
					"$in": branchFilters.ZoneIds,
				},
			},
		}
	}

	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if branchFilters.VendorId != nil {
		matchStage["vendorId"] = branchFilters.VendorId
	}

	if branchFilters.ManagerId != nil {
		matchStage["managersIds"] = bson.M{
			"$in": []primitive.ObjectID{*branchFilters.ManagerId},
		}
	}

	if addExtraOperations {

		if branchFilters.ProductBranchFilters != nil {
			pipeline = append(pipeline, bson.M{"$addFields": bson.M{"showProductOnIsActive": branchFilters.ProductBranchFilters.IsActive}})
		}

		pipeline = base_filters.BuildPopulatePipeline(branchFilters.PopFields, enums.BRANCH, pipeline)
		pipeline = base_filters.BuildSortingPipeline(branchFilters.Sorts, pipeline)
	}

	return pipeline
}
