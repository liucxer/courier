package expressions

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
)

type Expr interface {
	// Names expr name and aliases
	Names() []string
	// Len length of args
	Len() int
	// Exec final execute
	Exec(ctx context.Context, args ...interface{}) (interface{}, error)
}

type WithLitArg interface {
	LitArg() bool
}

type Expression = []interface{}

func resolveExpression(in interface{}) (string, []interface{}, bool) {
	if expr, ok := in.(Expression); ok {
		if len(expr) > 1 {
			if name, ok := expr[0].(string); ok {
				return name, expr[1:], ok
			}
		}
	}
	return "", nil, false
}

const LITERAL_EXPRESSION = "literal"

func StringifyValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return fmt.Sprintf("%v", strconv.Quote(v))
	case []interface{}:
		buf := bytes.NewBufferString("[")

		for i := range v {
			if i != 0 {
				buf.WriteString(",")
			}
			buf.WriteString(StringifyValue(v[i]))
		}

		buf.WriteString("]")

		return buf.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

func StringifyArgs(args []interface{}) string {
	buf := bytes.NewBuffer(nil)

	for i := range args {
		arg := args[i]

		if i != 0 {
			buf.WriteString(",")
		}

		if name, exprArgs, ok := resolveExpression(arg); ok {
			if name == LITERAL_EXPRESSION {
				arg = exprArgs[0]
			} else {
				buf.WriteString(StringifyExpression(arg.(Expression)))
				continue
			}
		}

		buf.WriteString(StringifyValue(arg))
	}

	return buf.String()
}

func StringifyExpression(expression Expression) string {
	return expression[0].(string) + "(" + StringifyArgs(expression[1:]) + ")"
}
