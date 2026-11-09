package types

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataWithPagination[T any] struct {
	Data       []T        `bson:"data" json:"data"`
	Pagination Pagination `bson:"pagination" json:"pagination"`
}

type Pagination struct {
	TotalPages float64 `json:"totalPages"`
	PerPage    int64   `json:"perPage"`
	TotalCount int64   `json:"totalCount"`
}

type User struct {
	Id       primitive.ObjectID
	UserData common.User
	Type     enums.UserType
}

type ServiceToken struct {
	Token string `bson:"token" json:"token"`
}

type UserList struct {
	Users []common.User `bson:"users" json:"users"`
}

type AdditionalOptionsMode string

const (
	CREATE AdditionalOptionsMode = "create"
	UPDATE AdditionalOptionsMode = "update"
)

type AdditionalOptions struct {
	UserId primitive.ObjectID    `bson:"userId" json:"userId"`
	Mode   AdditionalOptionsMode `bson:"mode,omitempty" json:"mode,omitempty"`
}

// MigrationResponse represents the response for migration operations
type MigrationResponse struct {
	TotalUsersProcessed int64    `json:"totalUsersProcessed"`
	ProfilesCreated     int64    `json:"profilesCreated"`
	ProfilesSkipped     int64    `json:"profilesSkipped"`
	Errors              []string `json:"errors,omitempty"`
	Message             string   `json:"message"`
}
