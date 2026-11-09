package product_utils

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetByBranchIdAndProductIdArgs struct {
	ProductId primitive.ObjectID
	BranchId  primitive.ObjectID
}

type GetWithPaginationArgs struct {
	Skip         int64
	Limit        int64
	Query        string
	BranchId     *primitive.ObjectID
	UnitId       *primitive.ObjectID
	CurrencyCode *string
}

type GetByIdArgs struct {
	Id           primitive.ObjectID
	Query        string
	BranchId     *primitive.ObjectID
	UnitId       *primitive.ObjectID
	CurrencyCode *string
}
