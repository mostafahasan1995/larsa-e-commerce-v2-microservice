package helpers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type BuildPaginationPipelineProps struct {
	BaseFilters    *base_filters.BaseFilters
	FilterPipeline []bson.M
	PostPipeline   []bson.M
	Skip           int64
	Limit          int64
}

// BuildPaginationPipeline builds a MongoDB aggregation pipeline for pagination
// The pipeline will always use "data" as the key for the paginated results
func BuildPaginationPipeline(collectionName enums.CollectionName, props BuildPaginationPipelineProps) []bson.M {
	dataStage := make([]bson.M, 0)
	dataStage = append(dataStage, bson.M{
		"$skip": props.Skip,
	})
	dataStage = append(dataStage, bson.M{
		"$limit": props.Limit,
	})

	dataStage = base_filters.BuildPopulatePipeline(props.BaseFilters.PopFields, collectionName, dataStage)

	dataStage = base_filters.BuildSortingPipeline(props.BaseFilters.Sorts, dataStage)

	if len(props.PostPipeline) > 0 {
		dataStage = append(dataStage, props.PostPipeline...)
	}

	return append(props.FilterPipeline, bson.M{
		"$facet": bson.M{
			"pagination": []bson.M{
				{"$count": "totalCount"},
				{"$addFields": bson.M{
					"perPage": props.Limit,
					"totalPages": bson.M{
						"$ceil": bson.M{
							"$divide": []interface{}{"$totalCount", props.Limit},
						},
					},
				}},
			},
			"data": dataStage,
		},
	}, bson.M{
		"$unwind": bson.M{
			"path":                       "$pagination",
			"preserveNullAndEmptyArrays": true,
		},
	})
}

// BuildOptimizedPaginationPipeline builds separate count and data pipelines to avoid $facet document size limits
func BuildOptimizedPaginationPipeline(collectionName enums.CollectionName, props BuildPaginationPipelineProps) ([]bson.M, []bson.M) {
	// Count pipeline - lightweight, no population or sorting
	countPipeline := make([]bson.M, 0)
	countPipeline = append(countPipeline, props.FilterPipeline...)
	countPipeline = append(countPipeline, bson.M{"$count": "totalCount"})

	// Data pipeline - includes population and sorting
	dataPipeline := make([]bson.M, 0)
	dataPipeline = append(dataPipeline, props.FilterPipeline...)

	// Add sorting before skip/limit for better performance
	dataPipeline = base_filters.BuildSortingPipeline(props.BaseFilters.Sorts, dataPipeline)

	dataPipeline = append(dataPipeline, bson.M{"$skip": props.Skip})
	dataPipeline = append(dataPipeline, bson.M{"$limit": props.Limit})

	// Add population after pagination for efficiency
	dataPipeline = base_filters.BuildPopulatePipeline(props.BaseFilters.PopFields, collectionName, dataPipeline)

	if len(props.PostPipeline) > 0 {
		dataPipeline = append(dataPipeline, props.PostPipeline...)
	}

	return countPipeline, dataPipeline
}
