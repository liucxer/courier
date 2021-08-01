package generator

import (
	"fmt"
	"go/types"
	"log"
	"path"
	"path/filepath"

	"github.com/liucxer/courier/codegen"
	"github.com/liucxer/courier/packagesx"
	"golang.org/x/tools/go/packages"

	"github.com/liucxer/courier/statuserror"
)

func NewStatusErrorGenerator(pkg *packagesx.Package) *StatusErrorGenerator {
	return &StatusErrorGenerator{
		pkg:          pkg,
		scanner:      NewStatusErrorScanner(pkg),
		statusErrors: map[string]*StatusError{},
	}
}

type StatusErrorGenerator struct {
	pkg          *packagesx.Package
	scanner      *StatusErrorScanner
	statusErrors map[string]*StatusError
}

func (g *StatusErrorGenerator) Scan(names ...string) {
	for _, name := range names {
		typeName := g.pkg.TypeName(name)
		g.statusErrors[name] = &StatusError{
			TypeName: typeName,
			Errors:   g.scanner.StatusError(typeName),
		}
	}
}

func getPkgDir(importPath string) string {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles,
	}, importPath)
	if err != nil {
		panic(err)
	}
	if len(pkgs) == 0 {
		panic(fmt.Errorf("package `%s` not found", importPath))
	}
	return filepath.Dir(pkgs[0].GoFiles[0])
}

func (g *StatusErrorGenerator) Output(cwd string) {
	for _, statusErr := range g.statusErrors {
		dir, _ := filepath.Rel(cwd, getPkgDir(statusErr.TypeName.Pkg().Path()))
		filename := codegen.GeneratedFileSuffix(path.Join(dir, codegen.LowerSnakeCase(statusErr.Name())+".go"))

		file := codegen.NewFile(statusErr.TypeName.Pkg().Name(), filename)
		statusErr.WriteToFile(file)

		if _, err := file.WriteFile(); err != nil {
			log.Printf("%s generated", file)
		}
	}
}

type StatusError struct {
	TypeName *types.TypeName
	Errors   []*statuserror.StatusErr
}

func (s *StatusError) Name() string {
	return s.TypeName.Name()
}

func (s *StatusError) WriteToFile(file *codegen.File) {
	s.WriteMethodImplements(file)

	s.WriteMethodStatusErrAndError(file)
	s.WriteMethodStatus(file)
	s.WriteMethodCode(file)

	s.WriteMethodKey(file)
	s.WriteMethodMsg(file)
	s.WriteMethodCanBeTalkError(file)
}

func (s *StatusError) WriteMethodImplements(file *codegen.File) {
	tpe := codegen.Type(file.Use("github.com/liucxer/courier/statuserror", "StatusError"))

	file.WriteBlock(
		file.Expr("var _ ? = (*?)(nil)", codegen.Interface(tpe), codegen.Type(s.Name())),
	)
}

func (s *StatusError) WriteMethodStatusErrAndError(file *codegen.File) {
	tpe := codegen.Type(file.Use("github.com/liucxer/courier/statuserror", "StatusErr"))

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("StatusErr").
			Return(codegen.Var(codegen.Star(tpe))).Do(
			file.Expr(`return &?{
Key: v.Key(),
Code: v.Code(),
Msg: v.Msg(),
CanBeTalkError: v.CanBeTalkError(),
}`, tpe)),
	)

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Unwrap").
			Return(codegen.Var(codegen.Error)).Do(
			file.Expr(`return v.StatusErr()`)),
	)

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Error").
			Return(codegen.Var(codegen.String)).Do(
			file.Expr(`return v.StatusErr().Error()`)),
	)
}

func (s *StatusError) WriteMethodStatus(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("StatusCode").
			Return(codegen.Var(codegen.Int)).Do(
			file.Expr(`return ?(int(v))`, codegen.Id(file.Use("github.com/liucxer/courier/statuserror", "StatusCodeFromCode"))),
		),
	)
}

func (s *StatusError) WriteMethodCode(file *codegen.File) {
	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Code").
			Return(codegen.Var(codegen.Int)).Do(
			file.Expr(`if withServiceCode, ok := (interface{})(v).(?); ok {
	return withServiceCode.ServiceCode() + int(v)
}
return int(v)
`, codegen.Id(file.Use("github.com/liucxer/courier/statuserror", "StatusErrorWithServiceCode"))),
		),
	)
}

func (s *StatusError) WriteMethodKey(file *codegen.File) {
	clauses := make([]*codegen.SnippetClause, 0)

	for _, statusErr := range s.Errors {
		clauses = append(clauses, codegen.Clause(codegen.Id(statusErr.Key)).Do(
			codegen.Return(
				file.Val(statusErr.Key),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Key").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val("UNKNOWN")),
		),
	)
}

func (s *StatusError) WriteMethodMsg(file *codegen.File) {
	clauses := make([]*codegen.SnippetClause, 0)

	for _, statusErr := range s.Errors {
		clauses = append(clauses, codegen.Clause(codegen.Id(statusErr.Key)).Do(
			codegen.Return(
				file.Val(statusErr.Msg),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("Msg").
			Return(codegen.Var(codegen.String)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val("-")),
		),
	)
}

func (s *StatusError) WriteMethodCanBeTalkError(file *codegen.File) {
	clauses := make([]*codegen.SnippetClause, 0)

	for _, statusErr := range s.Errors {
		clauses = append(clauses, codegen.Clause(codegen.Id(statusErr.Key)).Do(
			codegen.Return(
				file.Val(statusErr.CanBeTalkError),
			),
		))
	}

	file.WriteBlock(
		codegen.Func().
			MethodOf(codegen.Var(codegen.Type(s.Name()), "v")).
			Named("CanBeTalkError").
			Return(codegen.Var(codegen.Bool)).Do(
			codegen.Switch(codegen.Id("v")).When(
				clauses...,
			),
			codegen.Return(file.Val(false)),
		),
	)
}
