package jsonschema

import (
	"context"
	"github.com/liucxer/courier/gengo/pkg/gengo"
	gengotypes "github.com/liucxer/courier/gengo/pkg/types"
	"github.com/liucxer/courier/schema/pkg/generators/enum"
	"github.com/liucxer/courier/schema/pkg/jsonschema"
	"github.com/liucxer/courier/schema/pkg/jsonschema/extractors"
	"go/token"
	"go/types"
	"strings"
)

type Scanner struct {
	enumTypes   enum.EnumTypes
	definitions map[string]*jsonschema.Schema
}

func (g *Scanner) Init(c *gengo.Context) {
	g.enumTypes = enum.EnumTypes{}
	g.enumTypes.Walk(c, c.Package.Pkg().Path())

	g.definitions = map[string]*jsonschema.Schema{}
}

func (g *Scanner) SchemaFromType(ctx context.Context, pos token.Pos, tpe types.Type, def bool) (s *jsonschema.Schema) {
	defer func() {
		if s != nil {
			if def {
				if e, ok := g.enumTypes.ResolveEnumType(tpe); ok {
					values := make([]interface{}, len(e.Constants))
					labels := make([]string, len(e.Constants))

					for i := range e.Constants {
						values[i] = e.Value(e.Constants[i])
						labels[i] = e.Label(e.Constants[i])
					}

					if len(values) > 0 {
						if _, ok := values[0].(string); ok {
							s.Type = jsonschema.StringOrArray{"string"}
						}
					}

					s.Enum = values
					s.AddExtension(jsonschema.XEnumLabels, labels)
				}
			}

			if docGetter := extractors.DocGetterFromContext(ctx); docGetter != nil {
				_, doc := gengotypes.ExtractCommentTags("+", docGetter.Doc(pos))
				s.Description = strings.Join(doc, "\n")
			}
		}
	}()

	return extractors.OpenAPISchemaGetterFromGoType(ctx, tpe, def).OpenAPISchema(func(t string) jsonschema.Refer {
		return extractors.TypeName(t)
	})
}
