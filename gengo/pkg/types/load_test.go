package types

import (
	"go/types"
	"testing"

	"github.com/onsi/gomega"
)

func TestLoad(t *testing.T) {
	u, err := Load([]string{
		"github.com/liucxer/courier/gengo/testdata/a",
		"github.com/liucxer/courier/gengo/testdata/a/b",
	})

	gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

	p := u.Package("github.com/liucxer/courier/gengo/testdata/a")
	tpe := p.Type("FakeBool")

	gomega.NewWithT(t).Expect(p.Doc(tpe.Pos())).To(gomega.HaveLen(3))
	gomega.NewWithT(t).Expect(p.MethodsOf(tpe.Type().(*types.Named), false)).To(gomega.HaveLen(1))
	gomega.NewWithT(t).Expect(p.MethodsOf(tpe.Type().(*types.Named), true)).To(gomega.HaveLen(1))

	t.Run("ResultsOf", func(t *testing.T) {
		funcResults := map[string]string{
			"FuncSingleReturn":              `(2)`,
			"FuncSelectExprReturn":          `(string)`,
			"FuncWillCall":                  `(2, github.com/liucxer/courier/gengo/testdata/a.String)`,
			"FuncReturnWithCallDirectly":    `(2, github.com/liucxer/courier/gengo/testdata/a.String)`,
			"FuncWithNamedReturn":           `(2, github.com/liucxer/courier/gengo/testdata/a.String)`,
			"FuncSingleNamedReturnByAssign": `("1", "2", *github.com/pkg/errors.fundamental)`,
			"FunWithSwitch":                 `("a1" | "a2" | "a3", "b1" | "b2" | "b3")`,
			"FuncWithIf":                    `("a0" | "a1" | string)`,
			"FuncCallReturnAssign":          `(2, github.com/liucxer/courier/gengo/testdata/a.String)`,
			"FuncCallWithFuncLit":           `(1, "s")`,
			"FuncWithImportedCall":          `(2)`,
			"FuncCurryCall":                 `(int)`,
		}

		for k, r := range funcResults {
			t.Run(k, func(t *testing.T) {
				fn := p.Function(k)
				ar, _ := p.ResultsOf(fn)
				gomega.NewWithT(t).Expect(ar.String()).To(gomega.Equal(r))
			})
		}
	})
}
