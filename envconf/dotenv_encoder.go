package envconf

import (
	"github.com/liucxer/courier/reflectx"
	"go/ast"
	"reflect"
)

func NewDotEnvEncoder(envVars *EnvVars) *DotEnvEncoder {
	return &DotEnvEncoder{
		envVars: envVars,
	}
}

type DotEnvEncoder struct {
	envVars  *EnvVars
	flagsSet map[string]map[string]bool
}

func (d *DotEnvEncoder) setFlags(k string, flags map[string]bool) {
	if d.flagsSet == nil {
		d.flagsSet = map[string]map[string]bool{}
	}
	d.flagsSet[k] = flags
}

func (d *DotEnvEncoder) SecurityEncode(v interface{}) ([]byte, error) {
	walker := NewPathWalker()

	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	if err := d.scan(walker, rv); err != nil {
		return nil, err
	}

	return d.envVars.MaskBytes(), nil
}

func (d *DotEnvEncoder) Encode(v interface{}) ([]byte, error) {
	walker := NewPathWalker()

	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	if err := d.scan(walker, rv); err != nil {
		return nil, err
	}

	return d.envVars.Bytes(), nil
}

func (d *DotEnvEncoder) scan(walker *PathWalker, rv reflect.Value) error {
	kind := rv.Kind()

	setValue := func(rv reflect.Value) error {
		key := walker.String()

		envVar := &EnvVar{
			KeyPath:  key,
			Optional: true,
		}

		if d.flagsSet != nil {
			if flags, ok := d.flagsSet[key]; ok {
				envVar.metaFromFlags(flags)
			} else {
				if len(walker.path) > 1 {
					k := StringifyPath(walker.path[0 : len(walker.path)-1]...)
					if flags, ok := d.flagsSet[k]; ok {
						envVar.metaFromFlags(flags)
					}
				}
			}
		}

		if securityStringer, ok := rv.Interface().(SecurityStringer); ok {
			envVar.Mask = securityStringer.SecurityString()
		}

		text, err := reflectx.MarshalText(rv)
		if err != nil {
			return err
		}

		envVar.Value = string(text)

		d.envVars.Set(envVar)
		return nil
	}

	switch kind {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil
		}
		return d.scan(walker, rv.Elem())
	case reflect.Func, reflect.Interface, reflect.Chan, reflect.Map:
		// skip
	default:
		typ := rv.Type()
		if typ.Implements(interfaceTextMarshaller) {
			if err := setValue(rv); err != nil {
				return err
			}
			return nil
		}

		switch kind {
		case reflect.Array, reflect.Slice:
			for i := 0; i < rv.Len(); i++ {
				walker.Enter(i)
				if err := d.scan(walker, rv.Index(i)); err != nil {
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

				if flags != nil {
					d.setFlags(walker.String(), flags)
				}

				if err := d.scan(walker, rv.Field(i)); err != nil {
					return err
				}

				if !inline {
					walker.Exit()
				}
			}
		default:
			if err := setValue(rv); err != nil {
				return err
			}
		}
	}
	return nil
}
