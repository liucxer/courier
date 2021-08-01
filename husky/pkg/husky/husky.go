package husky

import (
	"context"
	"os"

	"github.com/liucxer/courier/husky/pkg/lintcommit"
	"github.com/liucxer/courier/husky/pkg/lintstaged"
	"gopkg.in/yaml.v2"
)

func HuskyFrom(ctx context.Context, huskyFile string) *Husky {
	data, err := os.ReadFile(huskyFile)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}

	s := NewSpec()

	if err := yaml.Unmarshal(data, s); err != nil {
		panic(err)
	}

	if s.VersionFile == "" {
		s.VersionFile = ".version"
	}

	h := &Husky{
		Spec: *s,
	}

	h.Init(ctx)

	return h
}

type Husky struct {
	Spec
	RunLintCommit func() error
	RunLintStated func() error
}

func (h *Husky) Init(ctx context.Context) {
	h.RunLintCommit = h.LintCommit.NewLint(ctx)
	h.RunLintStated = h.LintStaged.NewLint(ctx)
}

func NewSpec() *Spec {
	s := &Spec{
		Hooks:      map[string][]string{},
		LintStaged: lintstaged.LintStaged{},
		LintCommit: lintcommit.LintCommit{},
	}
	return s
}

type Spec struct {
	VersionFile string                `yaml:"version-file,omitempty"`
	Hooks       map[string][]string   `yaml:"hooks,omitempty"`
	LintStaged  lintstaged.LintStaged `yaml:"lint-staged,omitempty"`
	LintCommit  lintcommit.LintCommit `yaml:"lint-commit,omitempty"`
}
