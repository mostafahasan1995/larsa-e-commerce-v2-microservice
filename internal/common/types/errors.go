package types

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrNotFound = errors.New("item not found")
var ErrDupliateEmail = errors.New("duplicated email")
var ErrUnKnowen = errors.New("unknowen error")
var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")
var ErrBadRequest = errors.New("bad request")

var ErrErrors = []error{
	ErrUserNotFound,
	ErrDupliateEmail,
	ErrUnKnowen,
	ErrUnauthorized,
	ErrForbidden,
	ErrBadRequest,
}

func IsErrorInArray(err error, errs []error) bool {
	for _, e := range errs {
		if e == err {
			return true
		}
	}
	return false
}

func GetStatusCode(err error) int {
	switch err {
	case ErrUserNotFound, ErrNotFound:
		return 404
	case ErrDupliateEmail:
		return 409
	case ErrUnauthorized:
		return 401
	case ErrForbidden:
		return 403
	default:
		return 500
	}
}

// ValidationError represents a single validation error with details about
// the field that failed validation, its value, the validation tag that failed,
// and a human-readable error message.
type ValidationError struct {
	// Field is the name of the struct field that failed validation
	Field string `json:"field,omitempty"`
	// Value is the actual value that failed validation
	Value interface{} `json:"value,omitempty"`
	// Tag is the validation tag that failed (e.g., "required", "email", "min")
	Tag string `json:"tag,omitempty"`
	// Message is a human-readable error message explaining the validation failure
	Message string `json:"message,omitempty"`
}

func (v *ValidationError) Error() string {
	return v.Message
}
