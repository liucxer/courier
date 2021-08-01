package expressions

import (
	"context"

	"github.com/liucxer/courier/schema/pkg/expressions/raw"
	"github.com/pkg/errors"
)

func init() {
	DefaultFactory.Register(Add{})
	DefaultFactory.Register(Sub{})
	DefaultFactory.Register(Mul{})
	DefaultFactory.Register(Div{})
	DefaultFactory.Register(Mod{})
	DefaultFactory.Register(Pow{})
	DefaultFactory.Register(Abs{})
}

// Add
//
// Syntax
//     ["add", value, value]: int
type Add struct{}

func (Add) Names() []string {
	return []string{"add", "+"}
}

func (Add) Len() int {
	return 2
}

func (Add) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have 2 arguments")
	}
	return raw.Add(raw.ValueOf(args[1]), raw.ValueOf(args[0]))
}

// Sub
//
// Syntax
//     ["sub", value, value]: int
type Sub struct{}

func (Sub) Names() []string {
	return []string{"sub", "-"}
}

func (Sub) Len() int {
	return 2
}

func (Sub) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have 2 arguments")
	}
	return raw.Sub(raw.ValueOf(args[1]), raw.ValueOf(args[0]))
}

// Mul
//
// Syntax
//     ["mul", value, value]: int
type Mul struct{}

func (Mul) Names() []string {
	return []string{"mul", "*"}
}

func (Mul) Len() int {
	return 2
}

func (Mul) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have 2 arguments")
	}
	return raw.Mul(raw.ValueOf(args[1]), raw.ValueOf(args[0]))
}

// Div
//
// Syntax
//     ["div", value, value]: int
type Div struct{}

func (Div) Names() []string {
	return []string{"div", "/"}
}

func (Div) Len() int {
	return 2
}

func (Div) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have 2 arguments")
	}
	return raw.Div(raw.ValueOf(args[1]), raw.ValueOf(args[0]))
}

// Mod
//
// Syntax
//     ["mod", value, value]: int
type Mod struct{}

func (Mod) Names() []string {
	return []string{"mod", "%"}
}

func (Mod) Len() int {
	return 2
}

func (Mod) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have 2 arguments")
	}
	return raw.Mod(raw.ValueOf(args[1]), raw.ValueOf(args[0]))
}

// Pow
//
// Syntax
//     ["pow", value, value]: int
type Pow struct{}

func (Pow) Names() []string {
	return []string{"pow"}
}

func (Pow) Len() int {
	return 2
}

func (Pow) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("should have 2 arguments")
	}
	return raw.Pow(raw.ValueOf(args[1]), raw.ValueOf(args[0]))
}

// Abs
//
// Syntax
//     ["abs", value]: int
type Abs struct{}

func (Abs) Names() []string {
	return []string{"abs"}
}

func (Abs) Len() int {
	return 1
}

func (Abs) Exec(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("`abs` should have at least 1 arguments")
	}

	//math.Abs()

	return nil, errors.New("`abs` only support use for number value")
}
