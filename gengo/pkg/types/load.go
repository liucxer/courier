package types

import (
	"fmt"
	"golang.org/x/tools/go/packages"
)

const (
	LoadFiles     = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles
	LoadImports   = LoadFiles | packages.NeedImports
	LoadTypes     = LoadImports | packages.NeedTypes | packages.NeedTypesSizes
	LoadSyntax    = LoadTypes | packages.NeedSyntax | packages.NeedTypesInfo
	LoadAllSyntax = LoadSyntax | packages.NeedDeps | packages.NeedModule
)

func Load(patterns []string) (Universe, error) {
	c := &packages.Config{
		Mode: LoadAllSyntax,
	}

	pkgs, err := packages.Load(c, patterns...)
	if err != nil {
		return nil, err
	}

	u := Universe{}

	var register func(p *packages.Package)

	register = func(p *packages.Package) {
		for k := range p.Imports {
			importedPkg := p.Imports[k]

			if _, ok := u[importedPkg.PkgPath]; !ok {
				register(importedPkg)
			}
		}

		if len(p.Errors) > 0 {
			for i := range p.Errors {
				e := p.Errors[i]
				fmt.Println("[warning]", e.Pos, e.Msg)
			}
		}

		u[p.PkgPath] = newPkg(p, u)
	}

	for i := range pkgs {
		register(pkgs[i])
	}

	return u, nil
}

type Universe map[string]Package

func (u Universe) Package(pkgPath string) Package {
	return u[pkgPath]
}
