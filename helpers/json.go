package helpers

import (
	"reflect"
	"strings"
)

func ListJsonName(in any) []string {
	ret := make([]string, 0)
	t := reflect.TypeOf(in)
	// Handle pointer to struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ret
	}
	return append(ret, ListJsonNameForType(t)...)
}

func ListJsonNameForType(t reflect.Type) []string {
	ret := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type

		// Dereference pointer types
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// Recurse into nested structs
		if fieldType.Kind() == reflect.Struct {
			ret = append(ret, ListJsonNameForType(fieldType)...)
		} else {
			jsonName := strings.Split(field.Tag.Get("json"), ",")[0] // using split to ignore options like omitempty
			// Skip fields with no json tag or explicitly excluded
			if jsonName != "" && jsonName != "-" {
				ret = append(ret, jsonName)
			}
		}
	}
	return ret
}
