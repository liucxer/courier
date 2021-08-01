package namer

import (
	gengotypes "github.com/liucxer/courier/gengo/pkg/types"
)

type Namer interface {
	Name(gengotypes.TypeName) string
}

type NameSystems map[string]Namer

type Names map[gengotypes.TypeName]string

func NewRawNamer(pkgPath string, tracker ImportTracker) Namer {
	return &rawNamer{pkgPath: pkgPath, tracker: tracker}
}

type rawNamer struct {
	pkgPath string
	tracker ImportTracker
	Names
}

func (n *rawNamer) Name(typeName gengotypes.TypeName) string {
	if n.Names == nil {
		n.Names = Names{}
	}

	if name, ok := n.Names[typeName]; ok {
		return name
	}

	pkgPath := typeName.Pkg().Path()

	if pkgPath == n.pkgPath {
		name := typeName.Name()
		if name != "" {
			return name
		}
		return typeName.String()
	} else {
		n.tracker.AddType(typeName)
		return n.tracker.LocalNameOf(pkgPath) + "." + typeName.Name()
	}
}
