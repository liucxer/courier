package enum

import (
	"fmt"
	"go/constant"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"github.com/liucxer/courier/gengo/pkg/gengo"
)

type EnumTypes map[string]map[types.Type]*EnumType

func (e EnumTypes) ResolveEnumType(t types.Type) (*EnumType, bool) {
	if n, ok := t.(*types.Named); ok {
		if enumTypes, ok := e[n.Obj().Pkg().Path()]; ok && enumTypes != nil {
			if enumType, ok := enumTypes[t]; ok {
				return enumType, ok
			}
		}
	}
	return nil, false
}

func (e EnumTypes) Walk(gc *gengo.Context, inPkgPath string) {
	for i := range gc.Universe {
		p := gc.Universe[i]

		pkgPath := p.Pkg().Path()

		if pkgPath != inPkgPath {
			continue
		}

		constants := p.Constants()

		for k := range p.Constants() {
			if e[pkgPath] == nil {
				e[pkgPath] = map[types.Type]*EnumType{}
			}

			constv := constants[k]

			if e[pkgPath][constv.Type()] == nil {
				e[pkgPath][constv.Type()] = &EnumType{}
			}

			e[pkgPath][constv.Type()].Add(constv, p.Comment(constv.Pos())...)
		}
	}
}

type EnumType struct {
	ConstUnknown *types.Const
	Constants    []*types.Const
	Comments     map[*types.Const][]string
}

func (e *EnumType) IsIntStringer() bool {
	return e.ConstUnknown != nil && len(e.Constants) > 0
}

func (e *EnumType) Label(cv *types.Const) string {
	if comments, ok := e.Comments[cv]; ok {
		label := strings.Join(comments, "")

		if label != "" {
			return label
		}
	}

	return fmt.Sprintf("%v", e.Value(cv))
}

func (e *EnumType) Value(cv *types.Const) interface{} {
	var parts = strings.SplitN(cv.Name(), "__", 2)

	if len(parts) == 2 && strings.HasSuffix(strings.ToUpper(cv.Type().String()), "."+parts[0]) {
		return parts[1]
	}

	val := cv.Val()

	if val.Kind() == constant.Int {
		i, _ := strconv.ParseInt(val.String(), 10, 64)
		return i
	}

	s, _ := strconv.Unquote(val.String())
	return s
}

func (e *EnumType) Add(cv *types.Const, comments ...string) {
	if e.Comments == nil {
		e.Comments = map[*types.Const][]string{}
	}

	n := cv.Name()

	if n[0] == '_' {
		return
	}

	if strings.HasSuffix(n, "_UNKNOWN") {
		e.ConstUnknown = cv
		return
	}

	e.Comments[cv] = comments

	parts := strings.SplitN(n, "__", 2)

	if len(parts) == 2 {
		if strings.HasSuffix(strings.ToUpper(cv.Type().String()), "."+parts[0]) {
			e.Constants = append(e.Constants, cv)
		}
	} else {
		e.Constants = append(e.Constants, cv)
	}

	sort.Slice(e.Constants, func(i, j int) bool {
		return e.Constants[i].Val().String() < e.Constants[j].Val().String()
	})
}
