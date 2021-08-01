package kvcondition

import (
	"bytes"
	"encoding"
	"errors"
	"strconv"
)

func ParseKVCondition(r []byte) (*KVCondition, error) {
	node, err := newNodeScanner(r).ScanNode()
	if err != nil {
		return nil, err
	}
	kvCondition := &KVCondition{Node: node}
	return kvCondition, nil
}

// openapi:strfmt kv-condition
type KVCondition struct {
	Node
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*KVCondition)(nil)

func (v KVCondition) IsZero() bool  {
	return v.Node == nil
}


func (v KVCondition) MarshalText() ([]byte, error) {
	if v.IsZero() {
		return nil, nil
	}
	return []byte(v.String()), nil
}

func (v *KVCondition) UnmarshalText(data []byte) error {
	kvc, err := ParseKVCondition(data)
	if err != nil {
		return err
	}
	*v = *kvc
	return nil
}

type visitor func(visit visitor, node Node)

func (v *KVCondition) Range(cb func(condition *Rule)) {
	visit := func(next visitor, node Node) {
		if c, ok := node.(*Condition); ok {
			next(next, c.Left)
			next(next, c.Right)
			return
		}
		if label, ok := node.(*Rule); ok {
			cb(label)
		}
	}
	visit(visit, v.Node)
}

type Node interface {
	String() string
}

type Rule struct {
	Operator Operator
	Key      []byte
	Value    []byte
}

func (l *Rule) String() string {
	buf := &bytes.Buffer{}

	buf.WriteString(string(l.Key))

	if l.Operator != 0 {
		buf.WriteByte(' ')
		buf.WriteString(l.Operator.String())
	}

	if len(l.Value) != 0 {
		buf.WriteByte(' ')
		buf.WriteString(strconv.Quote(string(l.Value)))
	}

	return buf.String()
}

type Operator int

const (
	OperatorExists Operator = iota
	OperatorEqual
	OperatorNotEqual
	OperatorContains
	OperatorStartsWith
	OperatorEndsWith
)

func (op Operator) Of(key string, value string) *Rule {
	return &Rule{
		Operator: op,
		Key:      []byte(key),
		Value: func() []byte {
			if len(value) == 0 {
				return nil
			}
			return []byte(value)
		}(),
	}
}

var (
	UnknownOperator = errors.New("unknown operator")
)

func ParseOperator(op string) (Operator, error) {
	switch op {
	case "=":
		return OperatorEqual, nil
	case "!=":
		return OperatorNotEqual, nil
	case "*=":
		return OperatorContains, nil
	case "^=":
		return OperatorStartsWith, nil
	case "$=":
		return OperatorEndsWith, nil
	case "":
		return OperatorExists, nil
	}
	return OperatorExists, UnknownOperator
}

func (op Operator) String() string {
	switch op {
	case OperatorEqual:
		return "="
	case OperatorNotEqual:
		return "!="
	case OperatorContains:
		return "*="
	case OperatorStartsWith:
		return "^="
	case OperatorEndsWith:
		return "$="
	default:
		return ""
	}
}

func And(left Node, right Node) *Condition {
	return &Condition{
		Operator: ConditionOperatorAND,
		Left:     left,
		Right:    right,
	}
}

func Or(left Node, right Node) *Condition {
	return &Condition{
		Operator: ConditionOperatorOR,
		Left:     left,
		Right:    right,
	}
}

type ConditionOperator int

const (
	ConditionOperatorAND ConditionOperator = iota + 1
	ConditionOperatorOR
)

type Condition struct {
	Operator ConditionOperator
	Left     Node
	Right    Node
}

func (c *Condition) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteRune('(')
	buf.WriteRune(' ')

	if c.Left != nil {
		buf.WriteString(c.Left.String())
	}

	switch c.Operator {
	case ConditionOperatorAND:
		buf.WriteString(" & ")
	case ConditionOperatorOR:
		buf.WriteString(" | ")
	}

	if c.Right != nil {
		buf.WriteString(c.Right.String())
	}

	buf.WriteRune(' ')
	buf.WriteRune(')')

	return buf.String()
}
