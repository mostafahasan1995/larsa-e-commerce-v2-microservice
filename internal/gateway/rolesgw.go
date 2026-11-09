package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"slices"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleList struct {
	Roles []Role `json:"roles"`
}

type Role struct {
	Id   primitive.ObjectID `json:"_id"`
	Name string             `json:"name"`
}

type RolesGw struct{}

func NewRolesGw(i *do.Injector) (*RolesGw, error) {
	return &RolesGw{}, nil
}

func (rgw *RolesGw) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("roles?query={\"name\":{\"$eq\":\"%s\"}}", name)
	reqOp := &common.RequestParams{
		Service: "users",
		Path:    path,
		Method:  "GET",
		Data:    map[string]interface{}{},
		Header:  cfg.Hp,
	}

	res, err := common.CallService(reqOp)
	if err != nil {
		return nil, err
	}

	var result RoleList
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Roles) == 0 {
		return nil, types.ErrNotFound
	}

	return &result.Roles[0], nil
}

func (rgw *RolesGw) GetRelevantRoles(ctx context.Context, roles []string) ([]Role, error) {
	serviceToken, err := common.GetServiceToken("ecommerce", &common.HeaderParams{})

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	hp := common.HeaderParams{
		Client:       cfg.Hp.Client,
		ServiceToken: serviceToken,
	}
	if err != nil {
		return nil, err
	}

	reqOp := &common.RequestParams{
		Service: "users",
		Path:    "roles/all",
		Method:  "GET",
		Data:    map[string]interface{}{},
		Header:  &hp,
	}

	res, err := common.CallService(reqOp)
	if err != nil {
		return nil, err
	}

	var result RoleList

	json.NewDecoder(res.Body).Decode(&result)

	var matchRoles []Role = make([]Role, 0)

	for _, r := range result.Roles {
		for _, role := range roles {
			if r.Name == role {
				matchRoles = append(matchRoles, r)
				break
			}
		}
	}

	if len(matchRoles) > 0 {
		return matchRoles, nil
	}

	return nil, errors.New("role not found")

}

func (rgw *RolesGw) GetRolesByIds(ctx context.Context, roleIds []primitive.ObjectID) ([]Role, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	// Convert ObjectIDs to strings for the query
	var idStrings []string
	for _, id := range roleIds {
		idStrings = append(idStrings, id.Hex())
	}

	// Create query with $in operator to match multiple IDs
	query := fmt.Sprintf("roles?query={\"_id\":{\"$in\":[%s]}}", formatIdsForQuery(idStrings))

	reqOp := &common.RequestParams{
		Service: "users",
		Path:    query,
		Method:  "GET",
		Data:    map[string]interface{}{},
		Header:  cfg.Hp,
	}

	res, err := common.CallService(reqOp)
	if err != nil {
		return nil, err
	}

	var result RoleList
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Roles) == 0 {
		return nil, types.ErrNotFound
	}

	return result.Roles, nil
}

// formatIdsForQuery formats a slice of ID strings into the proper format for MongoDB $in query
func formatIdsForQuery(ids []string) string {
	if len(ids) == 0 {
		return ""
	}

	result := ""
	for i, id := range ids {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("\"%s\"", id)
	}
	return result
}

// IsUserBranchManager checks if a user has the "branchManager" role
// func (rgw *RolesGw) IsUserBranchManager(ctx context.Context) (bool, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return false, err
// 	}

// 	// Get user's role IDs
// 	userRoleIds := cfg.User.UserData.Roles
// 	if len(userRoleIds) == 0 {
// 		return false, nil
// 	}

// 	// Convert string IDs to ObjectIDs
// 	var objectIds []primitive.ObjectID
// 	for _, roleIdStr := range userRoleIds {
// 		objectId, err := primitive.ObjectIDFromHex(roleIdStr)
// 		if err != nil {
// 			continue // Skip invalid IDs
// 		}
// 		objectIds = append(objectIds, objectId)
// 	}

// 	if len(objectIds) == 0 {
// 		return false, nil
// 	}

// 	fmt.Println("objectIds", objectIds)
// 	// Get roles by IDs
// 	roles, err := rgw.GetRolesByIds(ctx, objectIds)
// 	if err != nil {
// 		return false, err
// 	}

// 	// Check if any role is "branchManager"
// 	for _, role := range roles {
// 		if role.Name == "branchManager" {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }

func (rgw *RolesGw) IsUserHaveRole(ctx context.Context, roleName string) (bool, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return false, err
	}

	// Get user's role IDs
	userRoleIds := cfg.User.UserData.Roles
	if len(userRoleIds) == 0 {
		return false, nil
	}

	// Convert string IDs to ObjectIDs
	var objectIds []primitive.ObjectID
	for _, roleIdStr := range userRoleIds {
		objectId, err := primitive.ObjectIDFromHex(roleIdStr)
		if err != nil {
			continue // Skip invalid IDs
		}
		objectIds = append(objectIds, objectId)
	}

	if len(objectIds) == 0 {
		return false, nil
	}

	fmt.Println("objectIds", objectIds)
	// Get roles by IDs
	roles, err := rgw.GetRolesByIds(ctx, objectIds)
	if err != nil {
		return false, err
	}

	// Check if any role is "branchManager"
	for _, role := range roles {
		if role.Name == roleName {
			return true, nil
		}
	}

	return false, nil
}

func (rgw *RolesGw) IsUserHaveAnyRole(ctx context.Context, roleNames ...string) (bool, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return false, err
	}

	// Get user's role IDs
	userRoleIds := cfg.User.UserData.Roles
	if len(userRoleIds) == 0 {
		return false, nil
	}

	// Convert string IDs to ObjectIDs
	var objectIds []primitive.ObjectID
	for _, roleIdStr := range userRoleIds {
		objectId, err := primitive.ObjectIDFromHex(roleIdStr)
		if err != nil {
			continue // Skip invalid IDs
		}
		objectIds = append(objectIds, objectId)
	}

	if len(objectIds) == 0 {
		return false, nil
	}

	fmt.Println("objectIds", objectIds)
	// Get roles by IDs
	roles, err := rgw.GetRolesByIds(ctx, objectIds)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		found := slices.Contains(roleNames, role.Name)
		if found {
			return true, nil
		}
	}

	return false, nil
}
