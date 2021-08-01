package reflectx

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

func MarshalText(v interface{}) ([]byte, error) {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(&v).Elem()
	}

	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return nil, nil
	}

	if textMarshaler, ok := rv.Interface().(encoding.TextMarshaler); ok {
		return textMarshaler.MarshalText()
	}

	switch rv.Kind() {
	case reflect.Interface, reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		}
		return MarshalText(rv.Elem())
	case reflect.String:
		return []byte(rv.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(fmt.Sprintf("%d", rv.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(fmt.Sprintf("%d", rv.Uint())), nil
	case reflect.Bool:
		return []byte(strconv.FormatBool(rv.Bool())), nil
	case reflect.Float32, reflect.Float64:
		return []byte(strconv.FormatFloat(rv.Float(), 'f', -1, 64)), nil
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

func UnmarshalText(v interface{}, data []byte) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	errorf := func(err error) error {
		return fmt.Errorf("cannot set value `%s`: %s", data, err)
	}

	irv := Indirect(rv)
	if irv.CanAddr() {
		if textUnmarshaler, ok := irv.Addr().Interface().(encoding.TextUnmarshaler); ok {
			if err := textUnmarshaler.UnmarshalText(data); err != nil {
				return errorf(err)
			}
			return nil
		}
	}

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			if rv.CanSet() {
				rv.Set(New(rv.Type()))
			}
		}
		return UnmarshalText(rv.Elem(), data)
	case reflect.String:
		rv.SetString(string(data))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intV, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return errorf(err)
		}
		rv.Set(reflect.ValueOf(intV).Convert(rv.Type()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintV, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return errorf(err)
		}
		rv.Set(reflect.ValueOf(uintV).Convert(rv.Type()))
	case reflect.Float32, reflect.Float64:
		floatV, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return errorf(err)
		}
		rv.Set(reflect.ValueOf(floatV).Convert(rv.Type()))
	case reflect.Bool:
		boolV, err := strconv.ParseBool(string(data))
		if err != nil {
			return errorf(err)
		}
		rv.SetBool(boolV)
	}
	return nil
}
