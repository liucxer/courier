package enum

import (
	"fmt"
	"go/types"
	"io"
	"strconv"

	"github.com/liucxer/courier/gengo/pkg/gengo"
	"github.com/liucxer/courier/gengo/pkg/namer"
	gengotypes "github.com/liucxer/courier/gengo/pkg/types"
)

func init() {
	gengo.Register(&enumGen{})
}

type enumGen struct {
	imports   namer.ImportTracker
	enumTypes EnumTypes
}

func (enumGen) Name() string {
	return "enum"
}

func (enumGen) New() gengo.Generator {
	return &enumGen{
		imports: namer.NewDefaultImportTracker(),
	}
}

func (g *enumGen) Init(c *gengo.Context, w io.Writer) error {
	if g.enumTypes == nil {
		g.enumTypes = EnumTypes{}
		g.enumTypes.Walk(c, c.Package.Pkg().Path())
	}
	return nil
}

func (g *enumGen) Imports(c *gengo.Context) map[string]string {
	return g.imports.Imports()
}

func (g *enumGen) GenerateType(c *gengo.Context, named *types.Named, w io.Writer) error {
	if enum, ok := g.enumTypes.ResolveEnumType(named); ok {

		if enum.IsIntStringer() {
			sw := gengo.NewSnippetWriter(w, namer.NameSystems{
				"raw": namer.NewRawNamer(c.Package.Pkg().Path(), g.imports),
			})

			g.genIntStringerMethods(sw, named, enum)
		}
	}
	return nil
}

func (g *enumGen) genIntStringerMethods(sw *gengo.SnippetWriter, tpe types.Type, enum *EnumType) {
	options := make([]struct {
		Name        string
		QuotedValue string
		QuotedLabel string
	}, len(enum.Constants))

	tpeObj := tpe.(*types.Named).Obj()

	for i := range enum.Constants {
		options[i].Name = enum.Constants[i].Name()
		options[i].QuotedValue = strconv.Quote(fmt.Sprintf("%v", enum.Value(enum.Constants[i])))
		options[i].QuotedLabel = strconv.Quote(enum.Label(enum.Constants[i]))
	}

	a := gengo.Args{
		"typeName":    tpeObj.Name(),
		"typePkgPath": tpeObj.Pkg().Path(),

		"constUnknown": enum.ConstUnknown,
		"options":      options,

		"ToUpper":        gengotypes.Ref("bytes", "ToUpper"),
		"NewError":       gengotypes.Ref("github.com/pkg/errors", "New"),
		"SqlDriverValue": gengotypes.Ref("database/sql/driver", "Value"),

		"IntStringerEnum":     gengotypes.Ref("github.com/liucxer/courier/schema/pkg/enumeration", "IntStringerEnum"),
		"ScanIntEnumStringer": gengotypes.Ref("github.com/liucxer/courier/schema/pkg/enumeration", "ScanIntEnumStringer"),
		"DriverValueOffset":   gengotypes.Ref("github.com/liucxer/courier/schema/pkg/enumeration", "DriverValueOffset"),
	}

	sw.Do(`
var Invalid[[ .typeName ]] = [[ .NewError | raw ]]("invalid [[ .typeName ]]")

func Parse[[ .typeName ]]FromString(s string) ([[ .typeName ]], error) {
	switch s {
	[[ range .options ]] case [[ .QuotedValue ]]:
		return [[ .Name ]], nil 
	[[ end ]] }
	return [[ .constUnknown | raw ]], Invalid[[ .typeName ]]
}

func Parse[[ .typeName ]]FromLabelString(s string) ([[ .typeName ]], error) {
	switch s {
	[[ range .options ]] case [[ .QuotedLabel ]]:
		return [[ .Name ]], nil
	[[ end ]] }
	return [[ .constUnknown | raw ]], Invalid[[ .typeName ]]
}

func ([[ .typeName ]]) TypeName() string {
	return "[[ .typePkgPath ]].[[ .typeName ]]"
}

func (v [[ .typeName ]]) String() string {
	switch v {
	[[ range .options ]] case [[ .Name ]]:
		return [[ .QuotedValue ]] 
	[[ end ]] }
	return "UNKNOWN"
}


func (v [[ .typeName ]]) Label() string {
	switch v {
	[[ range .options ]] case [[ .Name ]]:
		return [[ .QuotedLabel ]] 
	[[ end ]] }
	return "UNKNOWN"
}

func (v [[ .typeName ]]) Int() int {
	return int(v)
}

func ([[ .typeName ]]) ConstValues() [][[ .IntStringerEnum | raw ]] {
	return [][[ .IntStringerEnum | raw ]]{
		[[ range .options ]][[ .Name ]], 
		[[ end ]] }
}

func (v [[ .typeName ]]) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, Invalid[[ .typeName ]]
	}
	return []byte(str), nil
}

func (v *[[ .typeName ]]) UnmarshalText(data []byte) (err error) {
	*v, err = Parse[[ .typeName ]]FromString(string([[ .ToUpper | raw ]](data)))
	return
}


func (v [[ .typeName ]]) Value() ([[ .SqlDriverValue | raw ]], error) {
	offset := 0
	if o, ok := (interface{})(v).([[ .DriverValueOffset | raw ]]); ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}

func (v *[[ .typeName ]]) Scan(src interface{}) error {
	offset := 0
	if o, ok := (interface{})(v).([[ .DriverValueOffset | raw ]]); ok {
		offset = o.Offset()
	}

	i, err := [[ .ScanIntEnumStringer | raw ]](src, offset)
	if err != nil {
		return err
	}
	*v = [[ .typeName ]](i)
	return nil
}
`, a)
}
