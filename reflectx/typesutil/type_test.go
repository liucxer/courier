package typesutil

import (
	"context"
	"encoding"
	"go/types"
	"reflect"
	"testing"
	"unsafe"

	"github.com/liucxer/courier/ptr"
	"github.com/liucxer/courier/reflectx/typesutil/__fixtures__/typ"
	typ2 "github.com/liucxer/courier/reflectx/typesutil/__fixtures__/typ/typ"
	. "github.com/onsi/gomega"
)

func TestType(t *testing.T) {
	fn := func(a, b string) bool {
		return true
	}

	values := []interface{}{
		typ.DeepCompose{},

		func() *typ.Enum { v := typ.ENUM__ONE; return &v }(),
		typ.ENUM__ONE,

		reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem(),
		reflect.TypeOf((*interface {
			encoding.TextMarshaler
			Stringify(ctx context.Context, vs ...interface{}) string
			Add(a, b string) string
			Bytes() []byte
			s() string
		})(nil)).Elem(),

		unsafe.Pointer(t),

		make(typ.Chan),
		make(chan string, 100),

		typ.F,
		typ.Func(fn),
		fn,

		typ.String(""),
		"",
		typ.Bool(true),
		true,
		typ.Int(0),
		ptr.Int(1),
		int(0),
		typ.Int8(0),
		int8(0),
		typ.Int16(0),
		int16(0),
		typ.Int32(0),
		int32(0),
		typ.Int64(0),
		int64(0),
		typ.Uint(0),
		uint(0),
		typ.Uintptr(0),
		uintptr(0),
		typ.Uint8(0),
		uint8(0),
		typ.Uint16(0),
		uint16(0),
		typ.Uint32(0),
		uint32(0),
		typ.Uint64(0),
		uint64(0),
		typ.Float32(0),
		float32(0),
		typ.Float64(0),
		float64(0),
		typ.Complex64(0),
		complex64(0),
		typ.Complex128(0),
		complex128(0),
		typ.Array{},
		[1]string{},
		typ.Slice{},
		[]string{},
		typ.Map{},
		map[string]string{},
		typ.Struct{},
		struct{}{},
		struct {
			typ.Part
			Part2  typ2.Part
			a      string
			A      string `json:"a"`
			Struct struct {
				B string
			}
		}{},
	}

	for i := range values {
		check(t, values[i])
	}
}

func check(t *testing.T, v interface{}) {
	rtype, ok := v.(reflect.Type)
	if !ok {
		rtype = reflect.TypeOf(v)
	}
	ttype := NewTypesTypeFromReflectType(rtype)

	rt := FromRType(rtype)
	tt := FromTType(ttype)

	t.Run(FullTypeName(rt), func(t *testing.T) {
		NewWithT(t).Expect(rt.String()).To(Equal(tt.String()))
		NewWithT(t).Expect(rt.Kind().String()).To(Equal(tt.Kind().String()))
		NewWithT(t).Expect(rt.Name()).To(Equal(tt.Name()))
		NewWithT(t).Expect(rt.PkgPath()).To(Equal(tt.PkgPath()))
		NewWithT(t).Expect(rt.Comparable()).To(Equal(tt.Comparable()))
		NewWithT(t).Expect(rt.AssignableTo(FromRType(reflect.TypeOf("")))).To(Equal(tt.AssignableTo(FromTType(types.Typ[types.String]))))
		NewWithT(t).Expect(rt.ConvertibleTo(FromRType(reflect.TypeOf("")))).To(Equal(tt.ConvertibleTo(FromTType(types.Typ[types.String]))))

		NewWithT(t).Expect(rt.NumMethod()).To(Equal(tt.NumMethod()))

		for i := 0; i < rt.NumMethod(); i++ {
			rMethod := rt.Method(i)
			tMethod, ok := tt.MethodByName(rMethod.Name())
			NewWithT(t).Expect(ok).To(BeTrue())

			NewWithT(t).Expect(rMethod.Name()).To(Equal(tMethod.Name()))
			NewWithT(t).Expect(rMethod.PkgPath()).To(Equal(tMethod.PkgPath()))
			NewWithT(t).Expect(rMethod.Type().String()).To(Equal(tMethod.Type().String()))
		}

		{
			_, rOk := rt.MethodByName("String")
			_, tOk := tt.MethodByName("String")
			NewWithT(t).Expect(rOk).To(Equal(tOk))
		}

		{
			rReplacer, rIs := EncodingTextMarshalerTypeReplacer(rt)
			tReplacer, tIs := EncodingTextMarshalerTypeReplacer(tt)
			NewWithT(t).Expect(rIs).To(Equal(tIs))
			NewWithT(t).Expect(rReplacer.String()).To(Equal(tReplacer.String()))
		}

		if rt.Kind() == reflect.Array {
			NewWithT(t).Expect(rt.Len()).To(Equal(tt.Len()))
		}

		if rt.Kind() == reflect.Map {
			NewWithT(t).Expect(FullTypeName(rt.Key())).To(Equal(FullTypeName(tt.Key())))
		}

		if rt.Kind() == reflect.Array || rt.Kind() == reflect.Slice || rt.Kind() == reflect.Map {
			NewWithT(t).Expect(FullTypeName(rt.Elem())).To(Equal(FullTypeName(tt.Elem())))
		}

		if rt.Kind() == reflect.Struct {
			NewWithT(t).Expect(rt.NumField()).To(Equal(tt.NumField()))

			for i := 0; i < rt.NumField(); i++ {
				rsf := rt.Field(i)
				tsf := tt.Field(i)

				NewWithT(t).Expect(rsf.Anonymous()).To(Equal(tsf.Anonymous()))
				NewWithT(t).Expect(rsf.Tag()).To(Equal(tsf.Tag()))
				NewWithT(t).Expect(rsf.Name()).To(Equal(tsf.Name()))
				NewWithT(t).Expect(rsf.PkgPath()).To(Equal(tsf.PkgPath()))
				NewWithT(t).Expect(FullTypeName(rsf.Type())).To(Equal(FullTypeName(tsf.Type())))
			}

			if rt.NumField() > 0 {
				{
					rsf, _ := rt.FieldByName("A")
					tsf, _ := tt.FieldByName("A")

					NewWithT(t).Expect(rsf.Anonymous()).To(Equal(tsf.Anonymous()))
					NewWithT(t).Expect(rsf.Tag()).To(Equal(tsf.Tag()))
					NewWithT(t).Expect(rsf.Name()).To(Equal(tsf.Name()))
					NewWithT(t).Expect(rsf.PkgPath()).To(Equal(tsf.PkgPath()))
					NewWithT(t).Expect(FullTypeName(rsf.Type())).To(Equal(FullTypeName(tsf.Type())))

					{
						_, ok := rt.FieldByName("_")
						NewWithT(t).Expect(ok).To(BeFalse())
					}
					{
						_, ok := tt.FieldByName("_")
						NewWithT(t).Expect(ok).To(BeFalse())
					}
				}

				{
					rsf, _ := rt.FieldByNameFunc(func(s string) bool {
						return s == "A"
					})
					tsf, _ := tt.FieldByNameFunc(func(s string) bool {
						return s == "A"
					})

					NewWithT(t).Expect(rsf.Anonymous()).To(Equal(tsf.Anonymous()))
					NewWithT(t).Expect(rsf.Tag()).To(Equal(tsf.Tag()))
					NewWithT(t).Expect(rsf.Name()).To(Equal(tsf.Name()))
					NewWithT(t).Expect(rsf.PkgPath()).To(Equal(tsf.PkgPath()))
					NewWithT(t).Expect(FullTypeName(rsf.Type())).To(Equal(FullTypeName(tsf.Type())))

					{
						_, ok := rt.FieldByNameFunc(func(s string) bool {
							return false
						})
						NewWithT(t).Expect(ok).To(BeFalse())
					}
					{
						_, ok := tt.FieldByNameFunc(func(s string) bool {
							return false
						})
						NewWithT(t).Expect(ok).To(BeFalse())
					}
				}
			}
		}

		if rt.Kind() == reflect.Func {
			NewWithT(t).Expect(rt.NumIn()).To(Equal(tt.NumIn()))
			NewWithT(t).Expect(rt.NumOut()).To(Equal(tt.NumOut()))

			for i := 0; i < rt.NumIn(); i++ {
				rParam := rt.In(i)
				tParam := tt.In(i)
				NewWithT(t).Expect(rParam.String()).To(Equal(tParam.String()))
			}

			for i := 0; i < rt.NumOut(); i++ {
				rResult := rt.Out(i)
				tResult := tt.Out(i)
				NewWithT(t).Expect(rResult.String()).To(Equal(tResult.String()))
			}
		}

		if rt.Kind() == reflect.Ptr {
			rt = Deref(rt).(*RType)
			tt = Deref(tt).(*TType)

			NewWithT(t).Expect(rt.String()).To(Equal(tt.String()))
		}
	})
}

func TestTryNew(t *testing.T) {
	{
		_, ok := TryNew(FromRType(reflect.TypeOf(typ.Struct{})))
		NewWithT(t).Expect(ok).To(BeTrue())
	}
	{
		_, ok := TryNew(FromTType(NewTypesTypeFromReflectType(reflect.TypeOf(typ.Struct{}))))
		NewWithT(t).Expect(ok).To(BeFalse())
	}
}

func TestEachField(t *testing.T) {
	expect := []string{
		"a", "b", "bool", "c", "Part2",
	}

	{
		rtype := FromRType(reflect.TypeOf(typ.Struct{}))
		names := make([]string, 0)
		EachField(rtype, "json", func(field StructField, fieldDisplayName string, omitempty bool) bool {
			names = append(names, fieldDisplayName)
			return true
		})
		NewWithT(t).Expect(expect).To(Equal(names))
	}

	{
		ttype := FromTType(NewTypesTypeFromReflectType(reflect.TypeOf(typ.Struct{})))
		names := make([]string, 0)
		EachField(ttype, "json", func(field StructField, fieldDisplayName string, omitempty bool) bool {
			names = append(names, fieldDisplayName)
			return true
		})
		NewWithT(t).Expect(expect).To(Equal(names))
	}
}
