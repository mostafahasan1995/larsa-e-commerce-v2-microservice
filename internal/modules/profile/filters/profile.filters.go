package profile_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	profile_enums "larsa-e-commerce-microservice-v2/internal/modules/profile/enums"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DateRange represents a date range with greater than or equal and less than or equal
type DateRange struct {
	Gte *time.Time `json:"gte"`
	Lte *time.Time `json:"lte"`
}

// ProfileFilters represents the filters for a profile
type ProfileFilters struct {
	base_filters.BaseFilters
	CompanyName  *string                    `json:"companyName,omitempty"`
	Country      *string                    `json:"country,omitempty"`
	City         *string                    `json:"city,omitempty"`
	Industry     *[]string                  `json:"industry,omitempty"`
	IncludeGuest *bool                      `json:"includeGuest,omitempty"`
	ProfileType  *profile_enums.ProfileType `json:"profileType,omitempty"`
}

// NewProfileFilters creates a new instance of ProfileFilters from a query string
func NewProfileFilters(query string) (*ProfileFilters, error) {
	profileFilters := ProfileFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &profileFilters); err != nil {
			return nil, err
		}
	}

	return &profileFilters, nil
}

// BuildPipeline builds a MongoDB aggregation pipeline based on the provided filters
func (profileFilters *ProfileFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if profileFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &profileFilters.BaseFilters

	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.PROFILE)

	// Filter by profile type if provided
	if profileFilters.ProfileType != nil {
		matchStage["profileType"] = profileFilters.ProfileType
	}

	// Filter by company name if provided
	if profileFilters.CompanyName != nil {
		// For each language in the localizable field
		matchStage["companyName"] = bson.M{
			"$regex":   *profileFilters.CompanyName,
			"$options": "i",
		}
	}

	// Filter by country if provided
	if profileFilters.Country != nil && *profileFilters.Country != "" {
		matchStage["country"] = *profileFilters.Country
	}

	// Filter by city if provided
	if profileFilters.City != nil && *profileFilters.City != "" {
		matchStage["city"] = *profileFilters.City
	}

	// Filter by industry if provided
	if profileFilters.Industry != nil && len(*profileFilters.Industry) > 0 {
		matchStage["industry"] = bson.M{"$in": *profileFilters.Industry}
	}

	matchStage["firstName"] = bson.M{"$ne": "Guest"}
	// Filter by isGuest if provided
	if profileFilters.IncludeGuest != nil {
		matchStage["firstName"] = "Guest"
	}

	// Filter profiles based on user type
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg != nil && cfg.User != nil {
		// If user is Vendor or Member, they can only see their own profiles
		if cfg.User.Type == enums.Vendor || cfg.User.Type == enums.Member {
			matchStage["_id"] = cfg.User.Id // Since profile ID is now the same as user ID
		}
		// SuperAdmin can see all profiles, so no filter needed
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(profileFilters.PopFields, enums.PROFILE, pipeline)
		pipeline = base_filters.BuildSortingPipeline(profileFilters.Sorts, pipeline)
	}

	return pipeline
}
