package types

import (
	"go/ast"
	"go/types"
)

type TypeName interface {
	Pkg() *types.Package
	Name() string
	String() string
	Exported() bool
}

var _ TypeName = &types.TypeName{}

func Ref(pkgPath string, name string) TypeName {
	return &ref{pkgPath: pkgPath, name: name}
}

type ref struct {
	pkgPath string
	name    string
}

func (ref) Underlying() types.Type {
	return nil
}

func (r *ref) String() string {
	return r.pkgPath + "." + r.name
}

func (r *ref) Pkg() *types.Package {
	return types.NewPackage(r.pkgPath, "")
}

func (r *ref) Name() string {
	return r.name
}

func (r *ref) Exported() bool {
	return ast.IsExported(r.name)
}
