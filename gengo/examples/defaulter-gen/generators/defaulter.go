package generators

import (
	"go/types"
	"io"

	"github.com/liucxer/courier/gengo/pkg/gengo"
	"github.com/liucxer/courier/gengo/pkg/namer"
)

func init() {
	gengo.Register(&defaulterGen{})
}

type defaulterGen struct {
	imports namer.ImportTracker
}

func (defaulterGen) Name() string {
	return "defaulter"
}

func (defaulterGen) New() gengo.Generator {
	return &defaulterGen{
		imports: namer.NewDefaultImportTracker(),
	}
}

func (d defaulterGen) Imports(c *gengo.Context) map[string]string {
	return d.imports.Imports()
}

func (d defaulterGen) Init(c *gengo.Context, w io.Writer) error {
	return nil
}

func (d defaulterGen) GenerateType(c *gengo.Context, t *types.Named, w io.Writer) error {
	sw := gengo.NewSnippetWriter(w, namer.NameSystems{
		"raw": namer.NewRawNamer(c.Package.Pkg().Path(), d.imports),
	})

	args := map[string]interface{}{
		"type": t.Obj(),
	}

	if err := sw.Do(`
func(v *{{ .type | raw }}) SetDefault() {
}
`, args); err != nil {
		return err
	}

	return nil
}
