package reflectx

import (
	"errors"
	"reflect"
)

func IsNil(obj interface{}) bool {
	val := getReflectValue(obj)
	return obj == nil || !val.IsValid() || ((val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface) && val.IsNil())
}

func IsZero(obj interface{}) bool {
	v := getReflectValue(obj)
	if IsNil(v) {
		return true
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func IsNillable(obj interface{}) bool {
	switch getReflectValue(obj).Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.Chan, reflect.Ptr, reflect.Interface, reflect.Func, reflect.UnsafePointer:
		return true
	}
	return false
}

// GetValue dynamically retrieves the value of a field by name. Returns the zero Value if the field is not found
func GetValue(obj interface{}, field string) reflect.Value {
	if field == "" {
		return reflect.Value{}
	}
	val := getReflectValue(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.FieldByName(field)
}

// SetValue dynamically sets a value to a field by the given name
func SetValue(obj interface{}, field string, value interface{}) {
	if field == "" {
		return
	}
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	val.FieldByName(field).Set(reflect.ValueOf(value))
}

// GetNonPointerValue iterates over the V.Elem() of the object until the kind is not a pointer
// or an interface and returns its value.
func GetNonPointerValue(obj interface{}) (reflect.Value, error) {
	ret, _, err := getNonPointerValue(obj)
	return ret, err
}

func getReflectValue(obj interface{}) reflect.Value {
	if v, ok := obj.(reflect.Value); ok {
		return v
	}
	return reflect.ValueOf(obj)
}
func getNonPointerValue(obj interface{}) (reflect.Value, []reflect.Kind, error) {
	val := getReflectValue(obj)
	var typeTracker []reflect.Kind
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		typeTracker = append(typeTracker, val.Kind())
		val = val.Elem()
	}
	if !val.IsValid() {
		return reflect.Value{}, typeTracker, errors.New("value is invalid")
	}
	if IsNillable(val) && val.IsNil() {
		return reflect.Value{}, typeTracker, errors.New("value is nil")
	}
	return val, typeTracker, nil
}
