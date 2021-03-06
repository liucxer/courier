package codegen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/liucxer/courier/codegen/formatx"
	"golang.org/x/tools/go/packages"
)

func NewFile(pkgName string, filename string) *File {
	return &File{
		PkgName:  LowerSnakeCase(pkgName),
		filename: filename,
	}
}

type File struct {
	PkgName  string
	filename string
	imports  map[string]string
	bytes.Buffer
}

func (file *File) WriteBlock(ss ...Snippet) {
	for _, s := range ss {
		file.Write(s.Bytes())
		file.WriteString("\n\n")
	}
}

func (file *File) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(`package ` + LowerSnakeCase(file.PkgName) + `
`)

	if file.imports != nil {
		buf.WriteString(`import (
`)
		for importPath, alias := range file.imports {
			buf.WriteString(alias)
			buf.WriteString(" ")
			buf.WriteString(strconv.Quote(importPath))
			buf.WriteString("\n")
		}

		buf.WriteString(`)
`)
	}

	io.Copy(buf, &file.Buffer)

	return formatx.MustFormat(file.filename, buf.Bytes(), formatx.SortImportsProcess)
}

func (file *File) Expr(f string, args ...interface{}) SnippetExpr {
	return createExpr(file.importAliaser)(f, args...)
}

func (file *File) TypeOf(tpe reflect.Type) SnippetType {
	return createTypeOf(file.importAliaser)(tpe)
}

func (file *File) Val(v interface{}) Snippet {
	return createVal(file.importAliaser)(v)
}

func (file *File) importAliaser(importPath string) string {
	if file.imports == nil {
		file.imports = map[string]string{}
	}
	if file.imports[importPath] == "" {
		pkgs, err := packages.Load(nil, importPath)
		if err != nil {
			panic(err)
		}
		if len(pkgs) == 0 {
			panic(fmt.Errorf("`%s` not found", importPath))
		}
		importPath = pkgs[0].PkgPath
		file.imports[importPath] = LowerSnakeCase(importPath)
	}
	return file.imports[importPath]
}

func (file *File) Use(importPath string, exposedName string) string {
	return file.importAliaser(importPath) + "." + exposedName
}

func deVendor(importPath string) string {
	parts := strings.Split(importPath, "/vendor/")
	return parts[len(parts)-1]
}

func (file *File) WriteFile() (int, error) {
	dir := filepath.Dir(file.filename)

	if dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return -1, err
		}
	}

	f, err := os.Create(file.filename)
	defer f.Close()
	if err != nil {
		return -1, err
	}

	n3, err := f.Write(file.Bytes())
	if err != nil {
		return -1, err
	}

	if err := f.Sync(); err != nil {
		return -1, err
	}

	return n3, nil
}
