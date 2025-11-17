package helpers

import (
	"reflect"
	"strings"
)

func ListJsonName(in any) []string {
	ret := make([]string, 0)
	t := reflect.TypeOf(in)
	if t.Kind() != reflect.Struct {
		return ret
	}
	return append(ret, ListJsonNameForType(t)...)
}

func ListJsonNameForType(t reflect.Type) []string {
	ret := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			ret = append(ret, ListJsonNameForType(field.Type)...)
		} else if field.Type.Kind() == reflect.Ptr {
			ret = append(ret, ListJsonNameForType(field.Type.Elem())...)
		} else {
			jsonName := strings.Split(field.Tag.Get("json"), ",")[0] //using split to ignore options like omit and string
			ret = append(ret, jsonName)
		}
	}
	return ret
}
