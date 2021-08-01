package expressions

import (
	"context"
	"errors"
	"regexp"
	"sync"

	"github.com/liucxer/courier/schema/pkg/expressions/raw"
)

func init() {
	DefaultFactory.Register(Not{})
	DefaultFactory.Register(Neq{})
	DefaultFactory.Register(Lt{})
	DefaultFactory.Register(Lte{})
	DefaultFactory.Register(Eq{})
	DefaultFactory.Register(Gt{})
	DefaultFactory.Register(Gte{})
	DefaultFactory.Register(Match{})
	DefaultFactory.Register(All{})
	DefaultFactory.Register(Any{})
	DefaultFactory.Register(Case{})
	DefaultFactory.Register(Coalesce{})
}

// Not
//
// Syntax
//     ["not", value]: bool
type Not struct{}

func (Not) Names() []string {
	return []string{
		"not", "!",
	}
}

func (Not) Len() int {
	return 1
}

func (Not) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("should have at least 1 arguments")
	}

	arg, err := ResolveArg(ctx, args[0])
	if err != nil {
		return nil, err
	}
	if b, ok := arg.(bool); ok {
		return !b, nil
	}
	return nil, errors.New("should have bool argument")
}

// Neq
//
// Syntax
//     ["neq", compare_value, value]: bool
type Neq struct{}

func (Neq) Names() []string {
	return []string{
		"neq", "!=",
	}
}

func (Neq) Len() int {
	return 2
}

func (Neq) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have at least 2 arguments")
	}
	finalArgs, err := ResolveArgs(ctx, args...)
	if err != nil {
		return nil, err
	}
	ret, err := raw.Compare(raw.ValueOf(finalArgs[0]), raw.ValueOf(finalArgs[1]))
	if err != nil {
		return nil, err
	}
	return ret != 0, nil
}

// Lt
//
// Syntax
//     ["lt", compare_value, value]: bool
type Lt struct{}

func (Lt) Names() []string {
	return []string{
		"lt", "<",
	}
}

func (Lt) Len() int {
	return 2
}

func (Lt) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have at least 2 arguments")
	}
	finalArgs, err := ResolveArgs(ctx, args...)
	if err != nil {
		return nil, err
	}
	ret, err := raw.Compare(raw.ValueOf(finalArgs[1]), raw.ValueOf(finalArgs[0]))
	if err != nil {
		return nil, err
	}
	return ret < 0, nil
}

// Lte
//
// Syntax
//     ["lte", compare_value, value]: bool
type Lte struct{}

func (Lte) Names() []string {
	return []string{
		"lte", "<=",
	}
}

func (Lte) Len() int {
	return 2
}

func (Lte) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have at least 2 arguments")
	}
	finalArgs, err := ResolveArgs(ctx, args...)
	if err != nil {
		return nil, err
	}
	ret, err := raw.Compare(raw.ValueOf(finalArgs[1]), raw.ValueOf(finalArgs[0]))
	if err != nil {
		return nil, err
	}
	return ret <= 0, nil
}

// Eq
//
// Syntax
//     ["eq", compare_value, value]: bool
type Eq struct{}

func (Eq) Names() []string {
	return []string{
		"eq", "==",
	}
}

func (Eq) Len() int {
	return 2
}

func (Eq) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have at least 2 arguments")
	}
	finalArgs, err := ResolveArgs(ctx, args...)
	if err != nil {
		return nil, err
	}
	ret, err := raw.Compare(raw.ValueOf(finalArgs[1]), raw.ValueOf(finalArgs[0]))
	if err != nil {
		return nil, err
	}
	return ret == 0, nil
}

// Gt
//
// Syntax
//     ["gt", compare_value, value]: bool
type Gt struct{}

func (Gt) Names() []string {
	return []string{
		"gt", ">",
	}
}

func (Gt) Len() int {
	return 2
}

func (Gt) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have at least 2 arguments")
	}
	finalArgs, err := ResolveArgs(ctx, args...)
	if err != nil {
		return nil, err
	}
	ret, err := raw.Compare(raw.ValueOf(finalArgs[1]), raw.ValueOf(finalArgs[0]))
	if err != nil {
		return nil, err
	}
	return ret > 0, nil
}

// Gte
//
// Syntax
//     ["gte", compare_value, value]: bool
type Gte struct{}

func (Gte) Names() []string {
	return []string{
		"gte", ">=",
	}
}

func (Gte) Len() int {
	return 2
}

func (Gte) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("`gte` should have at least 2 arguments")
	}
	finalArgs, err := ResolveArgs(ctx, args...)
	if err != nil {
		return nil, err
	}
	ret, err := raw.Compare(raw.ValueOf(finalArgs[1]), raw.ValueOf(finalArgs[0]))
	if err != nil {
		return nil, err
	}
	return ret >= 0, nil
}

// Match
//
// Syntax
//     ["match", "pattern", value]: bool
type Match struct{}

func (Match) Names() []string {
	return []string{
		"match",
	}
}

func (Match) Len() int {
	return 2
}

var parsedPattern = sync.Map{}

func (Match) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have at least 2 arguments")
	}

	finalArgs, err := ResolveArgs(ctx, args...)
	if err != nil {
		return nil, err
	}

	pattern, ok := finalArgs[0].(string)
	if !ok {
		return nil, errors.New("first arg should be a valid pattern string")
	}

	re, ok := parsedPattern.Load(pattern)

	if !ok {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return nil, errors.New("`match` first arg should be a valid pattern string")
		}

		re = compiled
		parsedPattern.Store(pattern, compiled)
	}

	v, ok := finalArgs[1].(string)
	if !ok {
		return nil, errors.New("`match` second arg should be a valid string")
	}

	return re.(*regexp.Regexp).MatchString(v), nil
}

// All
//
// Syntax
//     ["all", ...]: bool
type All struct{}

func (All) Names() []string {
	return []string{"all"}
}

func (All) Len() int {
	return -1
}

func (All) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	for i := range args {
		arg, err := ResolveArg(ctx, args[i])
		if err != nil {
			return nil, err
		}
		if b, ok := arg.(bool); ok {
			if !b {
				return false, nil
			}
		}
	}
	return true, nil
}

// Any
//
// Syntax
//     ["any", ...]: bool
type Any struct{}

func (Any) Names() []string {
	return []string{"any"}
}

func (Any) Len() int {
	return -1
}

func (Any) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	for i := range args {
		arg, err := ResolveArg(ctx, args[i])
		if err != nil {
			return nil, err
		}
		if b, ok := arg.(bool); ok {
			if b {
				return true, nil
			}
		}
	}
	return false, nil
}

// Case
//
// Syntax
//     ["case",
//       condition: boolean, output: OutputType,
//       condition: boolean, output: OutputType,
//       ...,
//       fallback: OutputType
//      ]: OutputType
type Case struct{}

func (Case) Names() []string {
	return []string{"case"}
}

func (Case) Len() int {
	return -1
}

func (Case) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	n := len(args)
	partN := n / 2

	for i := 0; i < partN; i++ {
		condition, err := ResolveArg(ctx, args[2*i])
		if err != nil {
			return nil, err
		}

		if b, ok := condition.(bool); ok {
			if b {
				ret, err := ResolveArg(ctx, args[2*i+1])
				if err != nil {
					return nil, err
				}
				return ret, nil
			}
		} else {
			return nil, errors.New("`case` condition arg should be bool")
		}
	}

	if partN*2 != n {
		ret, err := ResolveArg(ctx, args[n-1])
		if err != nil {
			return nil, err
		}
		return ret, nil
	}

	return false, nil
}

// Coalesce
//
// Evaluates each expression in turn until the first non-null value is obtained, and returns that value.
//
// Syntax
//     ["coalesce", OutputType, OutputType, ...]: OutputType
type Coalesce struct{}

func (Coalesce) Names() []string {
	return []string{"coalesce"}
}

func (Coalesce) Len() int {
	return -1
}

func (Coalesce) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	for i := range args {
		arg, err := ResolveArg(ctx, args[i])
		if err != nil {
			return nil, err
		}
		if arg != nil {
			return arg, nil
		}
	}
	return false, nil
}
