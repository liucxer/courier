package gengo

import (
	"bytes"
	"fmt"
	"go/types"
	"reflect"
	"sort"
	"strconv"

	"github.com/liucxer/courier/gengo/pkg/namer"
	gengotypes "github.com/liucxer/courier/gengo/pkg/types"
	"github.com/liucxer/courier/reflectx"
	"github.com/liucxer/courier/reflectx/typesutil"
)

func NewDumper(pkg string, imports namer.ImportTracker) *Dumper {
	return &Dumper{
		namer: namer.NewRawNamer(pkg, imports),
	}
}

type Dumper struct {
	namer namer.Namer
}

func (d *Dumper) Name(named gengotypes.TypeName) string {
	return d.namer.Name(named)
}

func (d *Dumper) ReflectTypeLit(tpe reflect.Type) string {
	return d.TypeLit(typesutil.FromRType(tpe))
}

func (d *Dumper) TypesTypeLit(tpe types.Type) string {
	return d.TypeLit(typesutil.FromTType(tpe))
}

func (d *Dumper) TypeLit(tpe typesutil.Type) string {
	if tpe.PkgPath() != "" {
		return d.Name(gengotypes.Ref(tpe.PkgPath(), tpe.Name()))
	}

	switch tpe.Kind() {
	case reflect.Ptr:
		return "*" + d.TypeLit(tpe.Elem())
	case reflect.Chan:
		return "chan " + d.TypeLit(tpe.Elem())
	case reflect.Struct:
		b := bytes.NewBufferString("struct {")

		for i := 0; i < tpe.NumField(); i++ {
			f := tpe.Field(i)

			if !f.Anonymous() {
				_, _ = fmt.Fprintf(b, "%s ", f.Name())
			}

			b.WriteString(d.TypeLit(f.Type()))

			if tag := f.Tag(); tag != "" {
				_, _ = fmt.Fprintf(b, " `%s`", tag)
			}

			b.WriteString("\n")
		}

		b.WriteString("}")

		return b.String()
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", tpe.Len(), d.TypeLit(tpe.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", d.TypeLit(tpe.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", d.TypeLit(tpe.Key()), d.TypeLit(tpe.Elem()))
	default:
		return tpe.String()
	}
}

type ValueLitOpt struct {
	SubValue    bool
	OnInterface func(v interface{}) string
}

type ValueLitOptFn func(o *ValueLitOpt)

func OnInterface(onUnknown func(v interface{}) string) ValueLitOptFn {
	return func(o *ValueLitOpt) {
		o.OnInterface = onUnknown
	}
}

func SubValue(sub bool) ValueLitOptFn {
	return func(o *ValueLitOpt) {
		o.SubValue = sub
	}
}

var basicKinds = map[reflect.Kind]bool{
	reflect.Bool:       true,
	reflect.Int:        true,
	reflect.Int8:       true,
	reflect.Int16:      true,
	reflect.Int32:      true,
	reflect.Int64:      true,
	reflect.Uint:       true,
	reflect.Uint8:      true,
	reflect.Uint16:     true,
	reflect.Uint32:     true,
	reflect.Uint64:     true,
	reflect.Uintptr:    true,
	reflect.Float32:    true,
	reflect.Float64:    true,
	reflect.Complex64:  true,
	reflect.Complex128: true,
}

func (d *Dumper) ValueLit(in interface{}, optFns ...ValueLitOptFn) string {
	rv, ok := in.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(in)
	}

	o := &ValueLitOpt{}

	for i := range optFns {
		optFns[i](o)
	}

	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return "nil"
	}

	tpe := rv.Type()

	switch tpe.Kind() {
	case reflect.Ptr:
		kind := rv.Elem().Kind()
		if _, ok := basicKinds[kind]; ok {
			return fmt.Sprintf("func(v %s) *%s { return &v }(%s)", kind, kind, d.ValueLit(rv.Elem(), optFns...))
		}
		return fmt.Sprintf("&(%s)", d.ValueLit(rv.Elem(), optFns...))
	case reflect.Struct:
		buf := bytes.NewBufferString(d.ReflectTypeLit(tpe))
		buf.WriteString(`{`)

		c := 0

		for i := 0; i < rv.NumField(); i++ {
			f := rv.Field(i)
			ft := tpe.Field(i)

			if !reflectx.IsEmptyValue(f) {
				v := d.ValueLit(f, append(optFns, SubValue(true))...)

				if v == "" {
					continue
				}

				if c == 0 {
					buf.WriteString("\n")
				}

				buf.WriteString(ft.Name)
				buf.WriteString(":")
				buf.WriteString(v)
				buf.WriteString(",")
				buf.WriteString("\n")

				c++
			}
		}

		// no field
		if o.SubValue && c == 0 {
			return ""
		}

		buf.WriteString(`}`)

		return buf.String()
	case reflect.Map:
		buf := bytes.NewBufferString(d.ReflectTypeLit(tpe))
		buf.WriteString(`{`)

		keys := make([]string, 0)
		keyValues := map[string]reflect.Value{}

		for _, key := range rv.MapKeys() {
			k := key.String()
			keys = append(keys, k)
			keyValues[k] = rv.MapIndex(key)
		}

		sort.Strings(keys)

		for i, k := range keys {
			if i == 0 {
				buf.WriteString("\n")
			}

			buf.WriteString(strconv.Quote(k))
			buf.WriteString(":")
			buf.WriteString(d.ValueLit(keyValues[k], optFns...))
			buf.WriteString(",")
			buf.WriteString("\n")
		}

		buf.WriteString(`}`)
		return buf.String()
	case reflect.Slice, reflect.Array:
		buf := bytes.NewBufferString(d.ReflectTypeLit(tpe))
		buf.WriteString(`{`)

		for i := 0; i < rv.Len(); i++ {
			if i == 0 {
				buf.WriteString("\n")
			}

			buf.WriteString(d.ValueLit(rv.Index(i), append(optFns, SubValue(false))...))
			buf.WriteString(",")
			buf.WriteString("\n")
		}

		buf.WriteString(`}`)

		return buf.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
		return fmt.Sprintf("%d", rv.Int())
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return fmt.Sprintf("%d", rv.Uint())
	case reflect.Int32:
		if b, ok := rv.Interface().(rune); ok {
			r := strconv.QuoteRune(b)
			if len(r) == 3 {
				return r
			}
		}
		return fmt.Sprintf("%d", rv.Int())
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64)
	case reflect.String:
		return strconv.Quote(rv.String())
	case reflect.Interface:
		if o.OnInterface != nil {
			return o.OnInterface(rv.Interface())
		}
		return "nil"
	case reflect.Invalid:
		return "nil"
	default:
		panic(fmt.Errorf("%s is an unsupported type", tpe.String()))
	}
}
