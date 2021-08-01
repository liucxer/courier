package gengo

import (
	"go/types"
	"io"
)

type Generator interface {
	Name() string
	New() Generator
	Init(*Context, io.Writer) error
	Imports(*Context) map[string]string
	GenerateType(*Context, *types.Named, io.Writer) error
}

type GeneratorArgs struct {
	Inputs             []string
	OutputFileBaseName string
}
