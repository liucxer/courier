package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/liucxer/courier/packagesx"
)

func TestGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	p, _ := packagesx.Load(filepath.Join(cwd, "../__examples__"))

	g := NewGenerator(p)

	g.Scan("Protocol", "PullPolicy")
	g.Output(cwd)
}
