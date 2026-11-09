package order_utils

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/gateway"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BranchAccessResult struct {
	LimitAccess bool
	BranchId    *primitive.ObjectID
}

func CheckUserBranchAccess(ctx context.Context, injector *do.Injector) (*BranchAccessResult, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	db := do.MustInvoke[*mongo.Client](injector)
	collection := db.Database(cfg.Db).Collection(enums.BRANCH.GetCollectionName())

	rolesGw := do.MustInvoke[*gateway.RolesGw](injector)

	isBranchManager, err := rolesGw.IsUserHaveRole(ctx, "branchManager")
	if err != nil {
		return nil, err
	}

	if isBranchManager {
		branchesIds, err := getUserAssignedBranches(ctx, cfg.User.Id, collection)
		if err != nil {
			return nil, err
		}

		var branchId *primitive.ObjectID
		if len(branchesIds) > 0 {
			branchId = &branchesIds[0]
		}

		return &BranchAccessResult{
			LimitAccess: true,
			BranchId:    branchId,
		}, nil
	}

	return &BranchAccessResult{
		LimitAccess: false,
		BranchId:    nil,
	}, nil
}

func getUserAssignedBranches(ctx context.Context, userId primitive.ObjectID, collection *mongo.Collection) ([]primitive.ObjectID, error) {
	filter := bson.M{
		"managersIds": bson.M{"$in": []primitive.ObjectID{userId}},
		"deletedAt":   bson.M{"$exists": false},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var branches []struct {
		Id primitive.ObjectID `bson:"_id"`
	}
	if err := cursor.All(ctx, &branches); err != nil {
		return nil, err
	}

	branchIds := make([]primitive.ObjectID, len(branches))
	for i, branch := range branches {
		branchIds[i] = branch.Id
	}

	return branchIds, nil
}
