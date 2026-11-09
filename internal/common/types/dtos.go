package types

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RegisterValidators registers custom validators
func RegisterValidators(v *validator.Validate) {
	v.RegisterValidation("objectid", ValidateObjectID)
	v.RegisterValidation("localizable", ValidateLocalizable)
	v.RegisterValidation("objectids", ValidateObjectIDs)
	v.RegisterValidation("nonemptysliceptr", ValidateNonEmptySlicePtr)
}

// ValidateObjectID validates that a string is a valid MongoDB ObjectID
func ValidateObjectID(fl validator.FieldLevel) bool {
	_, err := primitive.ObjectIDFromHex(fl.Field().String())
	return err == nil
}

func ValidateLocalizable(fl validator.FieldLevel) bool {
	localizable := fl.Field().Interface().(Localizable[string])
	delete(localizable, "defaultTranslation")

	for key, _ := range localizable {
		languageExists := contains(enums.LanguageList, enums.LanguageKey(key))

		if !languageExists {
			return false
		}
	}

	return true
}

func contains[T comparable](arr []T, key T) bool {
	for _, value := range arr {
		if value == key {
			return true
		}
	}
	return false
}

// ValidateObjectIDs validates that all strings in a slice are valid MongoDB ObjectIDs
func ValidateObjectIDs(fl validator.FieldLevel) bool {
	ids, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	for _, id := range ids {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return false
		}
	}
	return true
}

// ValidateNonEmptySlicePtr validates that a pointer to a slice is not empty
func ValidateNonEmptySlicePtr(fl validator.FieldLevel) bool {
	field := fl.Field()

	// If it's not a pointer, skip
	if field.Kind() != reflect.Ptr {
		return true
	}

	// If pointer is nil, it's valid (omitempty should handle this)
	if field.IsNil() {
		return true
	}

	// Get the slice the pointer points to
	slice := field.Elem()

	// Check if it's a slice
	if slice.Kind() != reflect.Slice {
		return true
	}

	// Validate that slice has at least 1 element
	return slice.Len() >= 1
}

type BulkDeleteDto struct {
	Ids []primitive.ObjectID `json:"ids" validate:"required,min=1,dive,required"`
}

type GetByIdDto struct {
	Id string `param:"id" validate:"required,objectid"`
}

type FileFieldDto struct {
	Id           string     `bson:"_id" json:"_id" validate:"required"`
	Path         string     `bson:"path" json:"path" validate:"required,url"`
	Variants     *[]string  `bson:"variants,omitempty" json:"variants,omitempty"`
	Expire       *time.Time `bson:"expire,omitempty" json:"expire,omitempty"`
	OriginalName *string    `bson:"originalName,omitempty" json:"originalName,omitempty"`
	Service      *string    `bson:"service,omitempty" json:"service,omitempty"`
}
