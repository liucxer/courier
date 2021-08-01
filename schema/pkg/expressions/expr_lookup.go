package expressions

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/liucxer/courier/schema/pkg/expressions/raw"
)

func init() {
	DefaultFactory.Register(In{})
}

// In
//
// Syntax
//     ["in", "", value]: bool
//     ["in", [...], value]: bool
type In struct{}

func (In) Names() []string {
	return []string{
		"in",
	}
}

func (In) Len() int {
	return 2
}

func (In) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("`in` should have at least 2 arguments")
	}

	value := args[1]

	switch rng := args[0].(type) {
	case []interface{}:
		for i := range rng {
			if ret, _ := raw.Compare(raw.ValueOf(rng[i]), raw.ValueOf(value)); ret == 0 {
				return true, nil
			}
		}
	case string:
		if v, ok := value.(string); ok {
			if strings.Contains(rng, v) {
				return true, nil
			}
		} else {
			return nil, errors.New("`in` second value should be string value")
		}
	default:
		return nil, fmt.Errorf("unsupported range %#v", rng)
	}

	return false, nil
}
