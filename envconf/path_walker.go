package envconf

import (
	"bytes"
	"fmt"
)

func NewPathWalker() *PathWalker {
	return &PathWalker{
		path: []interface{}{},
	}
}

type PathWalker struct {
	path []interface{}
}

func (pw *PathWalker) Enter(i interface{}) {
	pw.path = append(pw.path, i)
}

func (pw *PathWalker) Exit() {
	pw.path = pw.path[:len(pw.path)-1]
}

func (pw *PathWalker) Paths() []interface{} {
	return pw.path
}

func (pw *PathWalker) String() string {
	return StringifyPath(pw.path...)
}

func StringifyPath(paths ...interface{}) string {
	buf := bytes.NewBuffer(nil)

	for i, key := range paths {
		if i > 0 {
			buf.WriteRune('_')
		}
		switch v := key.(type) {
		case string:
			buf.WriteString(v)
		case int:
			buf.WriteString(fmt.Sprintf("%d", v))
		}
	}

	return buf.String()
}
