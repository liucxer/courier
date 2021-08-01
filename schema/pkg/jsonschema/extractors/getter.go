package extractors

import (
	"context"
	"go/types"
	"reflect"

	"github.com/liucxer/courier/reflectx/typesutil"
	"github.com/liucxer/courier/schema/pkg/jsonschema"
)

type Definitions map[string]OpenAPISchemaGetter

func (Definitions) Ref(key string) jsonschema.Refer {
	return jsonschema.RefSchema("#/definitions/" + key)
}

var definitions = Definitions{}

func Register(key string, getter OpenAPISchemaGetter) {
	definitions[key] = getter
}

type OpenAPISchemaGetter interface {
	OpenAPISchema(ref func(t string) jsonschema.Refer) *jsonschema.Schema
}

type OpenAPISchemaTypeGetter interface {
	OpenAPISchemaType() []string
}

type OpenAPISchemaFormatGetter interface {
	OpenAPISchemaFormat() string
}

func OpenAPISchemaGetterFromReflectType(ctx context.Context, t reflect.Type) OpenAPISchemaGetter {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return &valueOpenAPISchemaGetter{
		ctx: ctx,
		t:   typesutil.FromRType(t),
	}
}

func OpenAPISchemaGetterFromGoType(ctx context.Context, types types.Type, def bool) OpenAPISchemaGetter {
	return &valueOpenAPISchemaGetter{
		ctx: ctx,
		t:   typesutil.FromTType(types),
		def: def,
	}
}

type valueOpenAPISchemaGetter struct {
	ctx context.Context
	t   typesutil.Type
	def bool
}

func (v *valueOpenAPISchemaGetter) OpenAPISchema(ref func(t string) jsonschema.Refer) *jsonschema.Schema {
	return SchemaFromType(v.ctx, v.t, v.def)
}

type opt struct {
	Ref func(t string) jsonschema.Refer
}

type SchemaFromOptFn = func(o *opt)

func SchemaFrom(v interface{}, optionFns ...SchemaFromOptFn) (s *jsonschema.Schema) {
	o := opt{}

	for i := range optionFns {
		optionFns[i](&o)
	}

	if o.Ref == nil {
		o.Ref = func(t string) jsonschema.Refer {
			return TypeName(t)
		}
	}

	defer func() {
		if g, ok := v.(OpenAPISchemaTypeGetter); ok {
			s.Type = g.OpenAPISchemaType()
			s.Format = ""
		}

		if g, ok := v.(OpenAPISchemaFormatGetter); ok {
			s.Format = g.OpenAPISchemaFormat()
		}
	}()

	if g, ok := v.(OpenAPISchemaGetter); ok {
		return g.OpenAPISchema(o.Ref)
	}

	return OpenAPISchemaGetterFromReflectType(context.Background(), reflect.TypeOf(v)).OpenAPISchema(o.Ref)
}
