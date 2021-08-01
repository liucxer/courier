package raw

import (
	"fmt"

	"github.com/pkg/errors"
)

// Div
// v / x
func Div(x Value, v Value) (interface{}, error) {
	switch x.Kind() {
	case Float:
		if x.Float() == 0 {
			return nil, errors.New("can't divide 0")
		}
		switch v.Kind() {
		case Int, Uint, Float:
			return v.Float() / x.Float(), nil
		}
	case Int:
		if x.Int() == 0 {
			return nil, errors.New("can't divide 0")
		}

		switch v.Kind() {
		case Int, Uint:
			if v.Int()%x.Int() == 0 {
				return v.Int() / x.Int(), nil
			}
			return v.Float() / x.Float(), nil
		case Float:
			return v.Float() / x.Float(), nil
		}
	case Uint:
		if x.Uint() == 0 {
			return nil, errors.New("can't divide 0")
		}
		switch v.Kind() {
		case Uint:
			if v.Uint()%x.Uint() == 0 {
				return v.Uint() / x.Uint(), nil
			}
			return v.Float() / x.Float(), nil
		case Int:
			if v.Int()%x.Int() == 0 {
				return v.Int() / x.Int(), nil
			}
			return v.Float() / x.Float(), nil
		case Float:
			return v.Float() / x.Float(), nil
		}
	}

	return nil, fmt.Errorf("%T can't divide %T", v, x)
}
