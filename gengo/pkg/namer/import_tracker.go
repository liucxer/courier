package namer

import (
	"go/token"
	"strings"

	gengotypes "github.com/liucxer/courier/gengo/pkg/types"
)

type ImportTracker interface {
	AddType(o gengotypes.TypeName)

	LocalNameOf(packagePath string) string
	PathOf(localName string) (string, bool)

	Imports() map[string]string
}

type defaultImportTracker struct {
	pathToName map[string]string
	nameToPath map[string]string
}

func NewDefaultImportTracker() ImportTracker {
	return &defaultImportTracker{
		pathToName: map[string]string{},
		nameToPath: map[string]string{},
	}
}

func (tracker *defaultImportTracker) AddType(o gengotypes.TypeName) {
	path := o.Pkg().Path()

	if _, ok := tracker.pathToName[path]; ok {
		return
	}

	localName := golangTrackerLocalName(path)

	tracker.nameToPath[localName] = path
	tracker.pathToName[path] = localName
}

func golangTrackerLocalName(name string) string {
	name = strings.Replace(name, "/", "_", -1)
	name = strings.Replace(name, ".", "_", -1)
	name = strings.Replace(name, "-", "_", -1)

	if token.Lookup(name).IsKeyword() {
		name = "_" + name
	}

	return name
}

func (tracker *defaultImportTracker) Imports() map[string]string {
	return tracker.pathToName
}

func (tracker *defaultImportTracker) LocalNameOf(path string) string {
	return tracker.pathToName[path]
}

func (tracker *defaultImportTracker) PathOf(localName string) (string, bool) {
	name, ok := tracker.nameToPath[localName]
	return name, ok
}
