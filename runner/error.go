package runner

import (
	"fmt"
)

type Error struct {
	Type ErrType
	Name string
}

func (e *Error) Error() string {
	return fmt.Sprintf("task %s %s", e.Name, e.Type)
}

type ErrType int

const (
	ErrTypeNormal ErrType = iota
	ErrTypeTimeout
	ErrTypeInterrupt
)

func (errType ErrType) String() string {
	switch errType {
	case ErrTypeTimeout:
		return "timeout"
	case ErrTypeInterrupt:
		return "interrupt"
	default:
		return "failed"
	}
}

func IsErrTimeout(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrTypeTimeout
	}
	return false
}

func IsErrInterrupt(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrTypeInterrupt
	}
	return false
}
