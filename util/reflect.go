package util

import (
	"fmt"
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

func SetValueByName(name string, value interface{}, result interface{}) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("failed to set value: %v", rec)
		}
	}()

	v := reflect.ValueOf(result)
	if v.Kind() != reflect.Ptr {
		err = fmt.Errorf("result must be a pointer")
		return
	}

	v = v.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		var (
			fieldValue = v.Field(i)
			fieldType  = t.Field(i)
		)

		if fieldType.Name != name {
			continue
		}

		// will panic if the value type does not match the field type
		// panic is handled by the defer function on the top of this function
		fieldValue.Set(reflect.ValueOf(value))
		return
	}

	err = fmt.Errorf("field %s not found", name)
	return
}
