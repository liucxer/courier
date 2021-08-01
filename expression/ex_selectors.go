package expression

import (
	"fmt"
	"reflect"
)

func init() {
	Register("len", func(f ExFactory, args []interface{}) (ExDo, error) {
		return BuildEx(`["len", any]`, 1, func(exes []ExDo) (ex ExDo, e error) {
			return func(params ...interface{}) (interface{}, error) {
				values, err := Exec(exes, params...)
				if err != nil {
					return nil, err
				}

				rv := reflect.ValueOf(values[0])

				switch rv.Kind() {
				case reflect.Slice, reflect.Array, reflect.Map:
					return rv.Len(), nil
				}

				return nil, NewRuntimeError(fmt.Sprintf("unexpected type for len, got type %T", values[0]))
			}, nil
		})(f, args)
	})

	Register("get", func(f ExFactory, args []interface{}) (ExDo, error) {
		return BuildEx(`["get", key: string]`, 1, func(exes []ExDo) (ex ExDo, e error) {
			return func(params ...interface{}) (interface{}, error) {
				values, err := Exec(exes, params...)
				if err != nil {
					return nil, err
				}

				key, ok := values[0].(string)
				if !ok {
					return nil, NewRuntimeError("key must be a string value")
				}

				if len(params) > 0 {
					if values, ok := params[0].(map[string]interface{}); ok {
						if v, ok := values[key]; ok {
							return v, nil
						}
					}
				}

				return nil, NewRuntimeError("missing target")
			}, nil
		})(f, args)
	})

	Register("has", func(f ExFactory, args []interface{}) (ex ExDo, e error) {
		return BuildEx(`["has", key: string]`, 1, func(exes []ExDo) (ex ExDo, e error) {
			return func(params ...interface{}) (interface{}, error) {
				values, err := Exec(exes, params...)
				if err != nil {
					return nil, err
				}

				key, ok := values[0].(string)
				if !ok {
					return nil, NewRuntimeError("key must be a string value")
				}

				if len(params) > 0 {
					if values, ok := params[0].(map[string]interface{}); ok {
						_, ok := values[key]
						return ok, nil
					}
				}

				return nil, NewRuntimeError("missing target")
			}, nil
		})(f, args)
	})
}
