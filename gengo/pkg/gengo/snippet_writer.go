package gengo

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
	"text/template"

	"github.com/liucxer/courier/gengo/pkg/namer"
)

type SnippetWriter struct {
	w  io.Writer
	ns namer.NameSystems
}

func NewSnippetWriter(w io.Writer, ns namer.NameSystems) *SnippetWriter {
	sw := &SnippetWriter{
		w:  w,
		ns: ns,
	}

	return sw
}

type Args = map[string]interface{}

func (s *SnippetWriter) Render(r func(s *SnippetWriter)) string {
	b := bytes.NewBuffer(nil)
	r(NewSnippetWriter(b, s.ns))
	return b.String()
}

func Snippet(format string, args ...Args) func(s *SnippetWriter) {
	return func(s *SnippetWriter) {
		s.Do(format, args...)
	}
}

func (s *SnippetWriter) Do(format string, args ...Args) {
	_, file, line, _ := runtime.Caller(1)

	funcMap := template.FuncMap{}

	for k := range s.ns {
		funcMap[k] = s.ns[k].Name
	}

	funcMap["render"] = s.Render

	tmpl, err := template.
		New(fmt.Sprintf("%s:%d", file, line)).
		Delims("[[", "]]").
		Funcs(funcMap).
		Parse(strings.TrimLeftFunc(format, func(r rune) bool {
			return r == '\n'
		}))

	if err != nil {
		panic(err)
	}

	finalArgs := Args{}

	for i := range args {
		a := args[i]
		for k := range a {
			finalArgs[k] = a[k]
		}
	}

	if err := tmpl.Execute(s.w, finalArgs); err != nil {
		panic(err)
	}
}

func (s *SnippetWriter) Println(format string) {

}
