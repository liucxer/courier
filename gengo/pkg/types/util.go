package types

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

func StringifyNode(fset *token.FileSet, node ast.Node) string {
	buf := bytes.Buffer{}
	if err := format.Node(&buf, fset, node); err != nil {
		panic(err)
	}
	return buf.String()
}
