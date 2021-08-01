package operator

import (
	"context"
	"fmt"
	"github.com/liucxer/courier/courier"
	"github.com/liucxer/courier/gengo/pkg/gengo"
	"github.com/liucxer/courier/gengo/pkg/namer"
	gengotypes "github.com/liucxer/courier/gengo/pkg/types"
	"github.com/liucxer/courier/httptransport/transformers"
	"github.com/liucxer/courier/reflectx/typesutil"
	"github.com/liucxer/courier/schema/pkg/generators/jsonschema"
	"github.com/liucxer/courier/schema/pkg/jsonschema/extractors"
	"github.com/liucxer/courier/schema/pkg/openapi"
	"github.com/pkg/errors"
	"go/types"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	gengo.Register(&operatorGen{})
}

type operatorGen struct {
	imports           namer.ImportTracker
	jsonschemaScanner jsonschema.Scanner
}

func (g *operatorGen) Name() string {
	return "operator"
}

func (operatorGen) New() gengo.Generator {
	return &operatorGen{
		imports: namer.NewDefaultImportTracker(),
	}
}

func (g *operatorGen) Imports(context *gengo.Context) map[string]string {
	return g.imports.Imports()
}

func (g *operatorGen) Init(c *gengo.Context, writer io.Writer) error {
	g.jsonschemaScanner.Init(c)
	return nil
}

var typOperator = reflect.TypeOf((*courier.Operator)(nil)).Elem()

func isCourierOperator(tpe typesutil.Type, lookup func(importPath string) *types.Package) bool {
	switch tpe.(type) {
	case *typesutil.RType:
		return tpe.Implements(typesutil.FromRType(typOperator))
	case *typesutil.TType:
		pkg := lookup(typOperator.PkgPath())
		if pkg == nil {
			return false
		}
		t := pkg.Scope().Lookup(typOperator.Name())
		if t == nil {
			return false
		}
		return types.Implements(tpe.(*typesutil.TType).Type, t.Type().Underlying().(*types.Interface))
	}
	return false
}

func (g *operatorGen) GenerateType(c *gengo.Context, named *types.Named, w io.Writer) error {
	sw := gengo.NewSnippetWriter(w, namer.NameSystems{
		"raw": namer.NewRawNamer(c.Package.Pkg().Path(), g.imports),
	})

	ptrTyp := typesutil.FromTType(types.NewPointer(named))

	if isCourierOperator(ptrTyp, func(importPath string) *types.Package { return c.Universe.Package(importPath).Pkg() }) {
		g.generateOpenAPIOperation(c, sw, named)
	}

	return nil
}

func (g *operatorGen) generateOpenAPIOperation(c *gengo.Context, sw *gengo.SnippetWriter, named *types.Named) {
	d := gengo.NewDumper(c.Package.Pkg().Path(), g.imports)

	o := &openapi.Operation{}
	o.OperationId = named.Obj().Name()

	docGetter := c.Universe.Package(named.Obj().Pkg().Path())
	ctx := extractors.WithDocGetter(context.Background(), docGetter)

	g.scanOperationInputsFromStruct(c, ctx, o, typesutil.FromTType(named.Underlying().(*types.Struct)))

	sw.Do(`
func([[ .typeName ]]) New() [[ .TOperator | raw ]] {
	return &[[ .typeName ]]{}
}

func([[ .typeName ]]) OpenAPIOperation(ref func(t string) [[ .TRefer | raw ]]) *[[ .TOperation | raw ]] {
	return [[ .operation ]]
}
`, gengo.Args{
		"TOperator":  gengotypes.Ref("github.com/liucxer/courier/courier", "Operator"),
		"TOperation": gengotypes.Ref("github.com/liucxer/courier/schema/pkg/openapi", "Operation"),
		"TRefer":     gengotypes.Ref("github.com/liucxer/courier/schema/pkg/jsonschema", "Refer"),

		"typeName": named.Obj().Name(),
		"operation": d.ValueLit(o, gengo.OnInterface(func(v interface{}) string {
			fmt.Printf("%T", v)

			switch t := v.(type) {
			case extractors.TypeName:
				return fmt.Sprintf(`ref(%s)`, strconv.Quote(string(t)))
			}
			return d.ValueLit(v)
		})),
	})
}

func (g *operatorGen) scanOperationInputsFromStruct(c *gengo.Context, ctx context.Context, op *openapi.Operation, s typesutil.Type) {
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		structTag := field.Tag()

		nameTag, hasNameTag := structTag.Lookup("name")

		if field.Anonymous() && !hasNameTag {
			if ft := field.Type(); ft.Kind() == reflect.Struct {
				g.scanOperationInputsFromStruct(c, ctx, op, ft)
				continue
			}
		}

		fieldDisplayName := field.Name()
		omitempty := false

		if nameTag != "" {
			if n := strings.Split(nameTag, ",")[0]; n != "" {
				fieldDisplayName = n
			}
			omitempty = strings.Contains(nameTag, "omitempty")
		}

		location := structTag.Get("in")

		if location == "" {
			panic(errors.Errorf("missing tag `in` for %s of %s", field.Name(), op.OperationId))
		}

		schema := extractors.PropSchemaFromStructField(
			ctx,
			field,
			!omitempty,
		)

		transformer, err := transformers.TransformerMgrDefault.NewTransformer(context.Background(), field.Type(), transformers.TransformerOption{
			MIME: structTag.Get("mime"),
		})

		if err != nil {
			panic(err)
		}

		switch location {
		case "body":
			reqBody := openapi.NewRequestBody("", true)
			reqBody.AddContent(transformer.Names()[0], openapi.NewMediaTypeWithSchema(schema))
			op.SetRequestBody(reqBody)
		case "query":
			op.AddParameter(openapi.QueryParameter(fieldDisplayName, schema, !omitempty))
		case "cookie":
			op.AddParameter(openapi.CookieParameter(fieldDisplayName, schema, !omitempty))
		case "header":
			op.AddParameter(openapi.HeaderParameter(fieldDisplayName, schema, !omitempty))
		case "path":
			op.AddParameter(openapi.PathParameter(fieldDisplayName, schema))
		}
	}
}
