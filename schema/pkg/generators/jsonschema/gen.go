package jsonschema

import (
	"context"
	"fmt"
	"go/types"
	"io"
	"strconv"

	"github.com/liucxer/courier/gengo/pkg/gengo"
	"github.com/liucxer/courier/gengo/pkg/namer"
	gengotypes "github.com/liucxer/courier/gengo/pkg/types"
	"github.com/liucxer/courier/schema/pkg/jsonschema/extractors"
)

func init() {
	gengo.Register(&jsonSchemaGen{})
}

type jsonSchemaGen struct {
	imports namer.ImportTracker
	scanner Scanner
}

func (g *jsonSchemaGen) Name() string {
	return "jsonschema"
}

func (jsonSchemaGen) New() gengo.Generator {
	return &jsonSchemaGen{
		imports: namer.NewDefaultImportTracker(),
	}
}

func (g *jsonSchemaGen) Init(c *gengo.Context, w io.Writer) error {
	g.scanner.Init(c)
	return nil
}

func (g *jsonSchemaGen) Imports(c *gengo.Context) map[string]string {
	return g.imports.Imports()
}

func (g *jsonSchemaGen) GenerateType(c *gengo.Context, named *types.Named, w io.Writer) error {
	sw := gengo.NewSnippetWriter(w, namer.NameSystems{
		"raw": namer.NewRawNamer(c.Package.Pkg().Path(), g.imports),
	})

	args := gengo.Args{
		"typeName":    named.Obj().Name(),
		"typePkgPath": named.Obj().Pkg().Path(),

		"FnRegister": gengotypes.Ref("github.com/liucxer/courier/schema/pkg/jsonschema/extractors", "Register"),
		"TSchema":    gengotypes.Ref("github.com/liucxer/courier/schema/pkg/jsonschema", "Schema"),
		"TRefer":     gengotypes.Ref("github.com/liucxer/courier/schema/pkg/jsonschema", "Refer"),
	}

	docGetter := c.Universe.Package(named.Obj().Pkg().Path())
	ctx := extractors.WithDocGetter(context.Background(), docGetter)

	s := g.scanner.SchemaFromType(ctx, named.Obj().Pos(), named, true)

	if s == nil {
		return nil
	}

	d := gengo.NewDumper(c.Package.Pkg().Path(), g.imports)

	sw.Do(`
func init() {
	[[ .FnRegister | raw ]]("[[ .typePkgPath ]].[[ .typeName ]]", new([[ .typeName ]]))
}

func([[ .typeName ]]) OpenAPISchema(ref func(t string) [[ .TRefer | raw ]]) *[[ .TSchema | raw ]] {
	return [[ .openAPISchema ]]
}

`, args, map[string]interface{}{
		"openAPISchema": d.ValueLit(s, gengo.OnInterface(func(v interface{}) string {
			switch t := v.(type) {
			case extractors.TypeName:
				return fmt.Sprintf(`ref(%s)`, strconv.Quote(string(t)))
			}
			return d.ValueLit(v)
		})),
	})

	return nil
}
