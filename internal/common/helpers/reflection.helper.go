package helpers

import "reflect"

// GetType returns the type of the provided interface
func GetType(i interface{}) reflect.Type {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

// CreateInstance creates a new instance of the provided type
func CreateInstance(t reflect.Type) interface{} {
	return reflect.New(t).Interface()
}
