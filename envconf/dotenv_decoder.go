package envconf

import (
	"go/ast"
	"reflect"

	"github.com/liucxer/courier/reflectx"
)

func NewDotEnvDecoder(envVars *EnvVars) *DotEnvDecoder {
	return &DotEnvDecoder{
		envVars: envVars,
	}
}

type DotEnvDecoder struct {
	envVars *EnvVars
}

func (d *DotEnvDecoder) Decode(v interface{}) error {
	walker := NewPathWalker()
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	return d.scanAndSetValue(walker, rv)
}

func (d *DotEnvDecoder) scanAndSetValue(walker *PathWalker, rv reflect.Value) error {
	kind := rv.Kind()

	if kind != reflect.Ptr && rv.CanAddr() {
		if defaultsSetter, ok := rv.Addr().Interface().(interface{ SetDefaults() }); ok {
			defaultsSetter.SetDefaults()
		}
	}

	switch kind {
	case reflect.Ptr:
		if rv.IsNil() {
			rv.Set(reflectx.New(rv.Type()))
		}
		return d.scanAndSetValue(walker, rv.Elem())
	case reflect.Func, reflect.Interface, reflect.Chan, reflect.Map:
		// skip
	default:
		typ := rv.Type()
		if typ.Implements(interfaceTextUnmarshaller) || reflect.PtrTo(typ).Implements(interfaceTextUnmarshaller) {
			v := d.envVars.Get(walker.String())
			if v != nil {
				if err := reflectx.UnmarshalText(rv, []byte(v.Value)); err != nil {
					return err
				}
			}
			return nil
		}

		switch kind {
		case reflect.Array, reflect.Slice:
			n := d.envVars.Len(walker.String())

			if kind == reflect.Slice && rv.IsNil() {
				rv.Set(reflect.MakeSlice(rv.Type(), n, n))
			}

			for i := 0; i < rv.Len(); i++ {
				walker.Enter(i)
				if err := d.scanAndSetValue(walker, rv.Index(i)); err != nil {
					return err
				}
				walker.Exit()
			}

		case reflect.Struct:
			tpe := rv.Type()
			for i := 0; i < rv.NumField(); i++ {
				field := tpe.Field(i)

				flags := (map[string]bool)(nil)
				name := field.Name

				if !ast.IsExported(name) {
					continue
				}

				if tag, ok := field.Tag.Lookup("env"); ok {
					n, fs := tagValueAndFlags(tag)
					if n == "-" {
						continue
					}
					if n != "" {
						name = n
					}
					flags = fs
				}

				inline := flags == nil && reflectx.Deref(field.Type).Kind() == reflect.Struct && field.Anonymous

				if !inline {
					walker.Enter(name)
				}

				if err := d.scanAndSetValue(walker, rv.Field(i)); err != nil {
					return err
				}

				if !inline {
					walker.Exit()
				}
			}
		default:
			v := d.envVars.Get(walker.String())
			if v != nil {
				if err := reflectx.UnmarshalText(rv, []byte(v.Value)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
