package generators

import (
	"context"
	"testing"

	"github.com/liucxer/courier/gengo/pkg/gengo"
)

func TestGenerator(t *testing.T) {
	c, _ := gengo.NewContext(&gengo.GeneratorArgs{
		Inputs: []string{
			"github.com/liucxer/courier/schema/testdata/a",
			"github.com/liucxer/courier/schema/testdata/b",
		},
		OutputFileBaseName: "zz_generated",
	})

	if err := c.Execute(context.Background(), gengo.GetRegisteredGenerators()...); err != nil {
		t.Fatal(err)
	}
}
