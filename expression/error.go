package expression

func NewInvalidExpression(msg string) *InvalidExpression {
	return &InvalidExpression{msg: msg}
}

type InvalidExpression struct {
	msg string
}

func (e *InvalidExpression) Error() string {
	return "invalid expression: " + e.msg
}

func NewRuntimeError(msg string) *RuntimeError {
	return &RuntimeError{msg: msg}
}

type RuntimeError struct {
	msg string
}

func (e *RuntimeError) Error() string {
	return "runtime error: " + e.msg
}
