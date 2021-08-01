package raw

import (
	"fmt"
)

func ValueOf(in interface{}) Value {
	switch v := in.(type) {
	case float64:
		return FloatValue(v)
	case float32:
		return FloatValue(v)
	case int:
		return IntValue(v)
	case int8:
		return IntValue(v)
	case int16:
		return IntValue(v)
	case int32:
		return IntValue(v)
	case int64:
		return IntValue(v)
	case uint:
		return UintValue(v)
	case uint8:
		return UintValue(v)
	case uint16:
		return UintValue(v)
	case uint32:
		return UintValue(v)
	case uint64:
		return UintValue(v)
	case string:
		return StringValue(v)
	case bool:
		return BoolValue(v)
	}
	return nil
}

type Kind int

const (
	Invalid Kind = iota
	Uint
	Int
	Float
	String
	Bool
)

type Value interface {
	Kind() Kind

	Int() int64
	Uint() uint64
	Float() float64
	Bool() bool
	String() string
}

type FloatValue float64

func (FloatValue) Kind() Kind {
	return Float
}

func (v FloatValue) Int() int64 {
	return int64(v)
}

func (v FloatValue) Uint() uint64 {
	return uint64(v)
}

func (v FloatValue) Float() float64 {
	return float64(v)
}

func (v FloatValue) Bool() bool {
	return v > 0
}

func (v FloatValue) String() string {
	return fmt.Sprintf("%v", float64(v))
}

type IntValue int64

func (IntValue) Kind() Kind {
	return Int
}

func (v IntValue) Int() int64 {
	return int64(v)
}

func (v IntValue) Uint() uint64 {
	return uint64(v)
}

func (v IntValue) Float() float64 {
	return float64(v)
}

func (v IntValue) Bool() bool {
	return v > 0
}

func (v IntValue) String() string {
	return fmt.Sprintf("%v", float64(v))
}

type UintValue uint64

func (UintValue) Kind() Kind {
	return Uint
}

func (v UintValue) Int() int64 {
	return int64(v)
}

func (v UintValue) Uint() uint64 {
	return uint64(v)
}

func (v UintValue) Float() float64 {
	return float64(v)
}

func (v UintValue) Bool() bool {
	return v > 0
}

func (v UintValue) String() string {
	return fmt.Sprintf("%v", float64(v))
}

type StringValue string

func (StringValue) Kind() Kind {
	return String
}

func (v StringValue) Int() int64 {
	return 0
}

func (v StringValue) Uint() uint64 {
	return 0
}

func (v StringValue) Float() float64 {
	return 0
}

func (v StringValue) Bool() bool {
	return v != ""
}

func (v StringValue) String() string {
	return string(v)
}

type BoolValue bool

func (BoolValue) Kind() Kind {
	return Bool
}

func (v BoolValue) Int() int64 {
	if v {
		return 1
	}
	return 0
}

func (v BoolValue) Uint() uint64 {
	if v {
		return 1
	}
	return 0
}

func (v BoolValue) Float() float64 {
	if v {
		return 1
	}
	return 0
}

func (v BoolValue) Bool() bool {
	return bool(v)
}

func (v BoolValue) String() string {
	if v {
		return "true"
	}
	return "false"
}
