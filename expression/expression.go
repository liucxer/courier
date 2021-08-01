package expression

import (
	"errors"
	"fmt"
)

var DefaultExFactory = ExSet{}

func Register(key string, expr ExBuilder) {
	DefaultExFactory.Register(key, expr)
}

func New(args []interface{}) (ExDo, error) {
	return DefaultExFactory.New(args)
}

func Ex(args ...interface{}) []interface{} {
	return args
}

type ExDo func(params ...interface{}) (interface{}, error)
type ExBuilder func(f ExFactory, args []interface{}) (ExDo, error)

type ExFactory interface {
	New(args []interface{}) (ExDo, error)
}

type ExSet map[string]ExBuilder

func (e ExSet) Register(key string, expr ExBuilder) {
	e[key] = expr
}

func (e ExSet) New(args []interface{}) (ExDo, error) {
	if len(args) > 0 {
		switch name := args[0].(type) {
		case string:
			if builder, ok := e[name]; ok {
				return builder(e, args[1:])
			}
		}

		return nil, errors.New("invalid expression")
	}
	return func(params ...interface{}) (i interface{}, e error) {
		return true, nil
	}, nil
}

func BuildEx(rule string, minArgN int, fn func(exes []ExDo) (ExDo, error)) ExBuilder {
	return func(f ExFactory, args []interface{}) (ExDo, error) {
		if len(args) >= minArgN {
			exes, err := withSubEx(f, args)
			if err != nil {
				return nil, err
			}
			return fn(exes)
		}
		return nil, NewInvalidExpression(fmt.Sprintf(`should be %s`, rule))
	}
}

func Exec(exes []ExDo, params ...interface{}) ([]interface{}, error) {
	values := make([]interface{}, len(exes))
	for i := range values {
		ex := exes[i]

		v, err := ex(params...)
		if err != nil {
			return nil, err
		}
		values[i] = v
	}
	return values, nil
}

func withSubEx(f ExFactory, args []interface{}) ([]ExDo, error) {
	expressions := make([]ExDo, len(args))

	for i := range expressions {
		arg := args[i]

		if args, ok := arg.([]interface{}); ok {
			ex, err := f.New(args)
			if err != nil {
				return nil, err
			}
			expressions[i] = ex
		} else {
			expressions[i] = func(params ...interface{}) (interface{}, error) {
				return arg, nil
			}
		}
	}

	return expressions, nil
}
