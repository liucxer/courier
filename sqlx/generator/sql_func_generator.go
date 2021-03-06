package generator

import (
	"go/types"
	"path"
	"path/filepath"

	"github.com/liucxer/courier/codegen"
	"github.com/liucxer/courier/packagesx"
)

func NewSqlFuncGenerator(pkg *packagesx.Package) *SqlFuncGenerator {
	return &SqlFuncGenerator{
		pkg: pkg,
	}
}

type SqlFuncGenerator struct {
	Config
	pkg   *packagesx.Package
	model *Model
}

type Config struct {
	StructName string
	TableName  string
	Database   string

	WithComments        bool
	WithTableName       bool
	WithTableInterfaces bool
	WithMethods         bool

	FieldPrimaryKey   string
	FieldKeyDeletedAt string
	FieldKeyCreatedAt string
	FieldKeyUpdatedAt string
}

func (g *Config) SetDefaults() {
	if g.FieldKeyDeletedAt == "" {
		g.FieldKeyDeletedAt = "DeletedAt"
	}

	if g.FieldKeyCreatedAt == "" {
		g.FieldKeyCreatedAt = "CreatedAt"
	}

	if g.FieldKeyUpdatedAt == "" {
		g.FieldKeyUpdatedAt = "UpdatedAt"
	}

	if g.TableName == "" {
		g.TableName = toDefaultTableName(g.StructName)
	}
}

func (g *SqlFuncGenerator) Scan() {
	for ident, obj := range g.pkg.TypesInfo.Defs {
		if typeName, ok := obj.(*types.TypeName); ok {
			if typeName.Name() == g.StructName {
				if _, ok := typeName.Type().Underlying().(*types.Struct); ok {
					g.model = NewModel(g.pkg, typeName, g.pkg.CommentsOf(ident), &g.Config)
				}
			}
		}
	}
}

func (g *SqlFuncGenerator) Output(cwd string) {
	if g.model == nil {
		return
	}

	dir, _ := filepath.Rel(cwd, filepath.Dir(g.pkg.GoFiles[0]))
	filename := codegen.GeneratedFileSuffix(path.Join(dir, codegen.LowerSnakeCase(g.StructName)+".go"))

	file := codegen.NewFile(g.pkg.Name, filename)
	g.model.WriteTo(file)
	_, _ = file.WriteFile()
}
