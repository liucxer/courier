package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/liucxer/courier/packagesx"
)

func TestGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	pkg, _ := packagesx.Load(filepath.Join(cwd, "../__examples__"))

	g := NewStatusErrorGenerator(pkg)

	g.Scan("StatusError")
	g.Output(cwd)
}
