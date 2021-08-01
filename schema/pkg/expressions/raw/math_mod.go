package raw

import (
	"fmt"
	"math"
)

// Mod
// v % x
func Mod(x Value, v Value) (interface{}, error) {
	switch x.Kind() {
	case Float:
		switch v.Kind() {
		case Int, Uint, Float:
			return fixDecimal(math.Mod(v.Float(), x.Float())), nil
		}
	case Int:
		switch v.Kind() {
		case Int, Uint:
			return v.Int() % x.Int(), nil
		case Float:
			return fixDecimal(math.Mod(v.Float(), x.Float())), nil
		}
	case Uint:
		switch v.Kind() {
		case Uint:
			return v.Uint() % x.Uint(), nil
		case Int:
			return v.Int() % x.Int(), nil
		case Float:
			return fixDecimal(math.Mod(v.Float(), x.Float())), nil
		}
	}

	return nil, fmt.Errorf("%T can't mod %T", v, x)
}
