package reflectx

import (
	"reflect"
	"testing"

	"github.com/liucxer/courier/ptr"
	. "github.com/onsi/gomega"
)

func TestIndirect(t *testing.T) {
	NewWithT(t).Expect(reflect.ValueOf(1).Interface()).To(Equal(Indirect(reflect.ValueOf(ptr.Int(1))).Interface()))
	NewWithT(t).Expect(reflect.ValueOf(0).Interface()).To(Equal(Indirect(reflect.New(reflect.TypeOf(0))).Interface()))

	rv := New(reflect.PtrTo(reflect.PtrTo(reflect.PtrTo(reflect.TypeOf("")))))
	NewWithT(t).Expect(reflect.ValueOf("").Interface()).To(Equal(Indirect(rv).Interface()))
}

type Zero string

func (Zero) IsZero() bool {
	return true
}

func TestIsEmptyValue(t *testing.T) {
	type S struct {
		V interface{}
	}

	emptyValues := []interface{}{
		Zero(""),
		(*string)(nil),
		(interface{})(nil),
		(S{}).V,
		"",
		0,
		uint(0),
		float32(0),
		false,
		reflect.ValueOf(S{}).FieldByName("V"),
		nil,
	}
	for _, v := range emptyValues {
		if rv, ok := v.(reflect.Value); ok {
			NewWithT(t).Expect(IsEmptyValue(rv)).To(BeTrue())
		} else {
			NewWithT(t).Expect(IsEmptyValue(v)).To(BeTrue())
			NewWithT(t).Expect(IsEmptyValue(reflect.ValueOf(v))).To(BeTrue())
		}

	}
}
