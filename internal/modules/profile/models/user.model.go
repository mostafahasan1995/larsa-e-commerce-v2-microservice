package profile_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the database for migration purposes
type User struct {
	Id           primitive.ObjectID     `bson:"_id" json:"id"`
	FirstName    string                 `bson:"firstName" json:"firstName"`
	LastName     string                 `bson:"lastName" json:"lastName"`
	Email        string                 `bson:"email" json:"email"`
	Password     string                 `bson:"password" json:"password"`
	Gender       string                 `bson:"gender" json:"gender"`
	Birthday     *time.Time             `bson:"birthday" json:"birthday"`
	Phone        string                 `bson:"phone" json:"phone"`
	IsActive     bool                   `bson:"isActive" json:"isActive"`
	IsVerified   bool                   `bson:"isVerified" json:"isVerified"`
	UserType     string                 `bson:"userType" json:"userType"`
	Roles        []primitive.ObjectID   `bson:"roles" json:"roles"`
	Capabilities []primitive.ObjectID   `bson:"capabilities" json:"capabilities"`
	CreatedAt    time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time              `bson:"updatedAt" json:"updatedAt"`
	LastLoginAt  *time.Time             `bson:"lastLoginAt" json:"lastLoginAt"`
	Avatar       string                 `bson:"avatar" json:"avatar"`
	Address      string                 `bson:"address" json:"address"`
	City         string                 `bson:"city" json:"city"`
	Country      string                 `bson:"country" json:"country"`
	PostalCode   string                 `bson:"postalCode" json:"postalCode"`
	Language     string                 `bson:"language" json:"language"`
	Currency     string                 `bson:"currency" json:"currency"`
	Timezone     string                 `bson:"timezone" json:"timezone"`
	Preferences  map[string]interface{} `bson:"preferences" json:"preferences"`
	Metadata     map[string]interface{} `bson:"metadata" json:"metadata"`
	IsDeleted    bool                   `bson:"isDeleted" json:"isDeleted"`
	DeletedAt    *time.Time             `bson:"deletedAt" json:"deletedAt"`
}

// UserResponse represents the response structure for user data
type UserResponse struct {
	Id           primitive.ObjectID     `bson:"_id" json:"id"`
	FirstName    string                 `bson:"firstName" json:"firstName"`
	LastName     string                 `bson:"lastName" json:"lastName"`
	Email        string                 `bson:"email" json:"email"`
	Gender       string                 `bson:"gender" json:"gender"`
	Birthday     *time.Time             `bson:"birthday" json:"birthday"`
	Phone        string                 `bson:"phone" json:"phone"`
	IsActive     bool                   `bson:"isActive" json:"isActive"`
	IsVerified   bool                   `bson:"isVerified" json:"isVerified"`
	UserType     string                 `bson:"userType" json:"userType"`
	Roles        []primitive.ObjectID   `bson:"roles" json:"roles"`
	Capabilities []primitive.ObjectID   `bson:"capabilities" json:"capabilities"`
	CreatedAt    time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time              `bson:"updatedAt" json:"updatedAt"`
	LastLoginAt  *time.Time             `bson:"lastLoginAt" json:"lastLoginAt"`
	Avatar       string                 `bson:"avatar" json:"avatar"`
	Address      string                 `bson:"address" json:"address"`
	City         string                 `bson:"city" json:"city"`
	Country      string                 `bson:"country" json:"country"`
	PostalCode   string                 `bson:"postalCode" json:"postalCode"`
	Language     string                 `bson:"language" json:"language"`
	Currency     string                 `bson:"currency" json:"currency"`
	Timezone     string                 `bson:"timezone" json:"timezone"`
	Preferences  map[string]interface{} `bson:"preferences" json:"preferences"`
	Metadata     map[string]interface{} `bson:"metadata" json:"metadata"`
}
