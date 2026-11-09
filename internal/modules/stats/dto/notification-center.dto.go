package stats_dto

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Your filter DTO
type GenericFilter struct {
	Collection   string   `json:"collection"`
	Field        string   `json:"field"`
	Operation    string   `json:"operation"`
	Logic        string   `json:"logic"`
	Value        any      `json:"value"`
	FilterDays   int      `json:"filterDays"`
	ReturnFields []string `json:"returnFields"`
	GroupBy      string   `json:"groupBy"`
	PreMatch     []bson.M `json:"preMatch"`
}

func BuildGenericMongoPipeline(f GenericFilter) bson.A {
	pipe := bson.A{}

	// 1. filter by days
	if f.FilterDays > 0 {
		pipe = append(pipe, bson.D{{
			Key: "$match",
			Value: bson.D{{
				Key: "createdAt",
				Value: bson.D{{
					Key:   "$gte",
					Value: time.Now().AddDate(0, 0, -f.FilterDays),
				}},
			}},
		}})
	}

	// 2. apply preMatch if given
	if len(f.PreMatch) > 0 {
		pipe = append(pipe, bson.D{{
			Key:   "$match",
			Value: bson.D{{Key: "$and", Value: f.PreMatch}},
		}})
	}

	// 3. build group stage
	var groupBy interface{}
	if strings.TrimSpace(f.GroupBy) == "" {
		groupBy = nil
	} else if strings.HasPrefix(f.GroupBy, "$") {
		groupBy = f.GroupBy
	} else {
		groupBy = "$" + f.GroupBy
	}
	groupStage := bson.D{
		{Key: "_id", Value: groupBy},
	}

	// capture some buyerInfo to return
	groupStage = append(groupStage, bson.E{Key: "buyerInfo", Value: bson.D{{Key: "$first", Value: "$buyerInfo"}}})

	if f.Operation == "$sum" && f.Field == "" {
		// count mode
		groupStage = append(groupStage, bson.E{Key: "result", Value: bson.D{{Key: "$sum", Value: 1}}})
	} else if f.Field != "" {
		groupStage = append(groupStage, bson.E{Key: "result", Value: bson.D{{Key: f.Operation, Value: "$" + f.Field}}})
	}

	pipe = append(pipe, bson.D{{Key: "$group", Value: groupStage}})

	// 4. apply logic filter on result (except when logic is $in → handle separately)
	if f.Logic != "" && f.Logic != "$in" {
		pipe = append(pipe, bson.D{{
			Key: "$match",
			Value: bson.D{{
				Key:   "result",
				Value: bson.D{{Key: f.Logic, Value: f.Value}},
			}},
		}})
	}

	// 5. projection: if ReturnFields is empty, return all fields
	if len(f.ReturnFields) == 0 {
		// keep all fields from the group, just add buyerId
		pipe = append(pipe, bson.D{{Key: "$addFields", Value: bson.D{{Key: "buyerId", Value: "$_id"}}}})
	} else {
		project := bson.D{{Key: "_id", Value: 0}, {Key: "result", Value: 1}}
		project = append(project, bson.E{Key: "buyerId", Value: "$_id"})
		for _, field := range f.ReturnFields {
			project = append(project, bson.E{Key: field, Value: 1})
		}
		pipe = append(pipe, bson.D{{Key: "$project", Value: project}})
	}

	return pipe
}
