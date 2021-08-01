package raw

import (
	"fmt"
	"math"
)

// Pow
// v ** x
func Pow(x Value, v Value) (interface{}, error) {
	switch x.Kind() {
	case Float:
		switch v.Kind() {
		case Int, Uint, Float:
			return fixDecimal(math.Pow(v.Float(), x.Float())), nil
		}
	case Int:
		switch v.Kind() {
		case Int, Uint:
			return int64(math.Pow(v.Float(), x.Float())), nil
		case Float:
			return fixDecimal(math.Pow(v.Float(), x.Float())), nil
		}
	case Uint:
		switch v.Kind() {
		case Uint:
			return uint64(math.Pow(v.Float(), x.Float())), nil
		case Int:
			return int64(math.Pow(v.Float(), x.Float())), nil
		case Float:
			return fixDecimal(math.Pow(v.Float(), x.Float())), nil
		}
	}

	return nil, fmt.Errorf("%T can't pow %T", v, x)
}
