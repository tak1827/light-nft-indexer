package util

import (
	"reflect"
)

func IsStruct(item interface{}) bool {
	v := reflect.ValueOf(item)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Kind() == reflect.Struct
}

func IsPointer(item interface{}) bool {
	v := reflect.ValueOf(item)
	return v.Kind() == reflect.Ptr
}

func IsPointerOfPointer(item interface{}) bool {
	rt := reflect.TypeOf(item)
	return rt.String()[:2] == "**"
}

func IsArrayOrSlice(item interface{}) bool {
	v := reflect.ValueOf(item)
	return v.Kind() == reflect.Slice || v.Kind() == reflect.Array
}

func Clone(inter interface{}) interface{} {
	nInter := reflect.New(reflect.TypeOf(inter).Elem())

	val := reflect.ValueOf(inter).Elem()
	nVal := nInter.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		nvField.Set(val.Field(i))
	}

	return nInter.Interface()
}
