package reflectx

import (
	"reflect"
)

func Indirect(rv reflect.Value) reflect.Value {
	if rv.Kind() == reflect.Ptr {
		return Indirect(rv.Elem())
	}
	return rv
}

func New(tpe reflect.Type) reflect.Value {
	rv := reflect.New(tpe).Elem()
	if tpe.Kind() == reflect.Ptr {
		rv.Set(New(tpe.Elem()).Addr())
		return rv
	}
	return rv
}

func IsEmptyValue(v interface{}) bool {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(&v).Elem()
	}

	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return true
	}

	if rv.IsValid() && rv.CanInterface() {
		if canZero, ok := rv.Interface().(interface{ IsZero() bool }); ok {
			return canZero.IsZero()
		}
	}

	switch rv.Kind() {
	case reflect.Interface:
		if rv.IsNil() {
			return true
		}
		return IsEmptyValue(rv.Elem())
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Ptr:
		return rv.IsNil()
	case reflect.Invalid:
		return true
	}
	return false
}
