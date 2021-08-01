package raw

import "fmt"

// Compare returns an integer comparing two Values .
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func Compare(a Value, b Value) (int, error) {
	switch a.Kind() {
	case Float:
		switch b.Kind() {
		case Float, Int, Uint:
			r := a.Float() - b.Float()
			if r > 0 {
				return 1, nil
			}
			if r < 0 {
				return -1, nil
			}
			return 0, nil
		}
	case Int:
		switch b.Kind() {
		case Float:
			r := a.Float() - b.Float()
			if r > 0 {
				return 1, nil
			}
			if r < 0 {
				return -1, nil
			}
			return 0, nil
		case Int, Uint:
			r := a.Int() - b.Int()
			if r > 0 {
				return 1, nil
			}
			if r < 0 {
				return -1, nil
			}
			return 0, nil
		}
	case Uint:
		switch b.Kind() {
		case Float:
			r := a.Float() - b.Float()
			if r > 0 {
				return 1, nil
			}
			if r < 0 {
				return -1, nil
			}
			return 0, nil
		case Int:
			r := a.Int() - b.Int()
			if r > 0 {
				return 1, nil
			}
			if r < 0 {
				return -1, nil
			}
			return 0, nil
		case Uint:
			r := a.Uint() - b.Uint()
			if r > 0 {
				return 1, nil
			}
			if r < 0 {
				return -1, nil
			}
			return 0, nil
		}
	case String:
		if b.Kind() == String {
			as := a.String()
			bs := b.String()
			if as < bs {
				return -1, nil
			}
			if as > bs {
				return 1, nil
			}
			return 0, nil
		}
	}
	return 0, fmt.Errorf("not comparable %T, %T", a, b)
}
