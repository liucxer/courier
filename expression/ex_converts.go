package expression

import (
	"fmt"
	"strconv"
)

func init() {
	Register("toNumber", createConvert(func(v interface{}) (interface{}, error) {
		v, err := toNumber(v)
		return v, err
	}))

	Register("toString", createConvert(func(v interface{}) (interface{}, error) {
		v, err := toString(v)
		return v, err
	}))

	Register("toBoolean", createConvert(func(v interface{}) (interface{}, error) {
		v, err := toBoolean(v)
		return v, err
	}))
}

func createConvert(convert func(v interface{}) (interface{}, error)) ExBuilder {
	return BuildEx(`["toNumber", any]`, 1, func(exes []ExDo) (ex ExDo, e error) {
		return func(params ...interface{}) (interface{}, error) {
			values, err := Exec(exes, params...)
			if err != nil {
				return nil, err
			}
			return convert(values[0])
		}, nil
	})
}

func toBoolean(v interface{}) (bool, error) {
	switch val := v.(type) {
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return val != "", nil
		}
		return b, nil
	case []byte:
		b, err := strconv.ParseBool(string(val))
		if err != nil {
			return len(val) > 0, nil
		}
		return b, nil
	case float64:
		return val != 0, nil
	case float32:
		return val != 0, nil
	case int:
		return val != 0, nil
	case int8:
		return val != 0, nil
	case int16:
		return val != 0, nil
	case int32:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case bool:
		return val, nil
	}
	return false, fmt.Errorf("unexpected type for toBoolean, got type %T", v)
}

func toString(v interface{}) (string, error) {
	return fmt.Sprintf("%v", v), nil
}

func toNumber(v interface{}) (float64, error) {
	switch val := v.(type) {
	case string:
		n, err := strconv.ParseFloat(val, 64)
		return n, err
	case []byte:
		n, err := strconv.ParseFloat(string(val), 64)
		return n, err
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	}

	return 0, fmt.Errorf("unexpected type for toNumber, got type %T", v)
}
