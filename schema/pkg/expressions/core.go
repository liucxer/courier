package expressions

import (
	"context"

	"github.com/pkg/errors"
)

type ExecFunc = func(ctx context.Context, inputs ...interface{}) (interface{}, error)

var DefaultFactory = Factory{}

type Factory map[string]Expr

func (f Factory) Register(expressions ...Expr) {
	for i := range expressions {
		expr := expressions[i]
		if expr != nil {
			for _, n := range expr.Names() {
				f[n] = expr
			}
		}
	}
}

func (f Factory) From(expression Expression) (ExecFunc, error) {
	name, args, ok := resolveExpression(expression)
	if !ok {
		return nil, errors.New("invalid expression, should be [string, ...any]")
	}
	return f.from(name, args)
}

func (f Factory) from(name string, args []interface{}) (ExecFunc, error) {

	expr, ok := f[name]
	if !ok {
		return nil, errors.Errorf("`%s` is not registered expression", name)
	}

	definedN := len(args)
	n := definedN
	inputNeeds := 0

	if l := expr.Len(); l != -1 {
		n = l
		inputNeeds = l - definedN
	}

	finalInputs := make([]interface{}, n)

	for i := range args {
		arg := args[i]

		if name, exprArgs, ok := resolveExpression(arg); ok {
			if name == LITERAL_EXPRESSION {
				finalInputs[i] = exprArgs[0]
				continue
			}
			fn, err := f.from(name, exprArgs)
			if err != nil {
				return nil, err
			}
			finalInputs[i] = fn
			continue
		}

		finalInputs[i] = arg
	}

	return func(ctx context.Context, inputs ...interface{}) (interface{}, error) {
		if ctx.Value(ctxKeyInputs{}) == nil {
			// only need to inject to ctx once
			ctx = ContextWithInputs(ctx, inputs...)
		}

		for i := 0; i < inputNeeds; i++ {
			finalInputs[definedN+i] = inputs[i]
		}

		return expr.Exec(ctx, finalInputs...)
	}, nil
}

type ctxKeyInputs struct {
}

func ContextWithInputs(ctx context.Context, inputs ...interface{}) context.Context {
	return context.WithValue(ctx, ctxKeyInputs{}, inputs)
}

func InputsFromContext(ctx context.Context) []interface{} {
	if inputs, ok := ctx.Value(ctxKeyInputs{}).([]interface{}); ok {
		return inputs
	}
	return []interface{}{}
}

func ResolveArg(ctx context.Context, arg interface{}) (interface{}, error) {
	if fn, ok := arg.(ExecFunc); ok {
		return fn(ctx, InputsFromContext(ctx)...)
	}
	return arg, nil
}

func ResolveArgs(ctx context.Context, args ...interface{}) ([]interface{}, error) {
	finals := make([]interface{}, len(args))

	for i := range finals {
		arg, err := ResolveArg(ctx, args[i])
		if err != nil {
			return nil, err
		}
		finals[i] = arg
	}

	return finals, nil
}
