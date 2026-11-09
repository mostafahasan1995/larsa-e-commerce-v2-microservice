package middleware

import (
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"reflect"
)

// validateNestedStructure recursively validates nested structures and their custom validators.
// It handles different types of nested structures including:
// - Pointer types
// - Arrays and slices
// - Structs that implement CustomValidator interface
// - Nested struct fields
//
// The function returns a slice of CustomValidatorContext and an optional ApiError.
// Each CustomValidatorContext contains validation data that needs to be added to the request context.
//
// How it works:
// 1. For pointer types: Dereferences the pointer and continues validation on the underlying value
// 2. For arrays/slices: Iterates through each element and validates them recursively
// 3. For structs:
//   - Checks if the struct implements CustomValidator interface
//   - If yes, calls CustomValidator method and collects validation contexts
//   - Recursively validates all fields within the struct
//
// When to use:
// - When you need to validate complex nested data structures
// - When different parts of your data structure need different validation rules
// - When you need to collect validation contexts from multiple levels of nesting
//
// Why it works this way:
// - Recursive approach allows handling arbitrary depth of nesting
// - Interface-based validation enables custom validation logic per type
// - Context collection enables passing validated data through the request chain
//
// Example usage with different types of structures:
//
//  1. Simple struct with custom validation:
//     type User struct {
//     Name string `json:"name"`
//     }
//     func (u *User) CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
//     // Custom validation logic
//     return []CustomValidatorContext{{Key: "user", Value: u}}, nil
//     }
//
//  2. Nested struct with array:
//     type Order struct {
//     Items []OrderItem `json:"items"`
//     }
//     type OrderItem struct {
//     ProductID string `json:"productId"`
//     }
//     func (o *OrderItem) CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
//     // Validate each order item
//     return []CustomValidatorContext{{Key: "item", Value: o}}, nil
//     }
//
//  3. Pointer to struct with nested validation:
//     type Address struct {
//     Street string `json:"street"`
//     }
//     type Customer struct {
//     Address *Address `json:"address"`
//     }
//     func (a *Address) CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
//     // Validate address
//     return []CustomValidatorContext{{Key: "address", Value: a}}, nil
//     }
//
//  4. Complex nested structure with multiple levels:
//     type Product struct {
//     ID          string    `json:"id"`
//     Categories  []Category `json:"categories"`
//     Attributes  *ProductAttributes `json:"attributes"`
//     }
//     type Category struct {
//     ID   string `json:"id"`
//     Name string `json:"name"`
//     }
//     type ProductAttributes struct {
//     Color string `json:"color"`
//     Size  string `json:"size"`
//     }
//     func (p *Product) CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
//     // Validate product
//     return []CustomValidatorContext{{Key: "product", Value: p}}, nil
//     }
//     func (c *Category) CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
//     // Validate category
//     return []CustomValidatorContext{{Key: "category", Value: c}}, nil
//     }
//
//  5. Optional fields with nil handling:
//     type Profile struct {
//     User     *User      `json:"user"`
//     Settings *Settings  `json:"settings"`
//     }
//     type Settings struct {
//     Theme string `json:"theme"`
//     }
//     func (p *Profile) CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
//     // Validate profile, handles nil User and Settings
//     return []CustomValidatorContext{{Key: "profile", Value: p}}, nil
//     }
//
//  6. Array of pointers with custom validation:
//     type Team struct {
//     Members []*Member `json:"members"`
//     }
//     type Member struct {
//     ID   string `json:"id"`
//     Role string `json:"role"`
//     }
//     func (m *Member) CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
//     // Validate each team member
//     return []CustomValidatorContext{{Key: "member", Value: m}}, nil
//     }
//
//  7. Deeply nested structure example (showing unlimited nesting capability):
//     type Organization struct {
//     Departments []Department
//     }
//     type Department struct {
//     Teams []Team
//     }
//     type Team struct {
//     Projects []Project
//     }
//     type Project struct {
//     Tasks []Task
//     }
//     type Task struct {
//     Subtasks []Subtask
//     }
//     type Subtask struct {
//     Assignments []Assignment
//     }
//     type Assignment struct {
//     User *User
//     }
//
//     // Each level can implement CustomValidator
//     func (o *Organization) CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
//     return []CustomValidatorContext{{Key: "organization", Value: o}}, nil
//     }
//     func (d *Department) CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
//     return []CustomValidatorContext{{Key: "department", Value: d}}, nil
//     }
//     // ... and so on for each level
//
// The function will validate all levels of nesting, collecting validation contexts
// from each level that implements CustomValidator. There is no practical limit to
// the depth of nesting, though extremely deep structures should be avoided for
// performance and maintainability reasons.
func validateNestedStructure(value reflect.Value, args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError) {
	var contexts []CustomValidatorContext

	// Handle pointer types
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, nil
		}
		value = value.Elem()
	}

	// Handle arrays and slices
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			nestedContexts, err := validateNestedStructure(value.Index(i), args)
			if err != nil {
				return nil, err
			}
			contexts = append(contexts, nestedContexts...)
		}
		return contexts, nil
	}

	// Handle structs
	if value.Kind() == reflect.Struct {
		// Check if the struct implements CustomValidator
		if validator, ok := value.Addr().Interface().(CustomValidator); ok {
			nestedContexts, err := validator.CustomValidator(args)
			if err != nil {
				return nil, err
			}
			contexts = append(contexts, nestedContexts...)
		}

		// Recursively validate nested struct fields
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			nestedContexts, err := validateNestedStructure(field, args)
			if err != nil {
				return nil, err
			}
			contexts = append(contexts, nestedContexts...)
		}
	}

	return contexts, nil
}
