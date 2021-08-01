package a

import (
	"github.com/pkg/errors"
	"strings"

	"github.com/liucxer/courier/gengo/testdata/a/b"
)

func Example() {

}

func (String) Method() string {
	return strings.Join(strings.Split("1", ","), ",")
}

func FuncSingleReturn() interface{} {
	// should skip
	_ = func() bool {
		a := true
		return !a
	}()

	var a interface{}
	a = "" + "1"
	a = 2

	return a
}

func FuncSelectExprReturn() string {
	v := struct{ s string }{}
	v.s = "2"
	return v.s
}

func FuncWillCall() (a interface{}, s String) {
	return FuncSingleReturn(), String(FuncSelectExprReturn())
}

func FuncReturnWithCallDirectly() (a interface{}, b String) {
	return FuncWillCall()
}

func FuncWithNamedReturn() (a interface{}, b String) {
	a, b = FuncWillCall()
	return
}

func newErr() error {
	return errors.New("some err")
}

func FuncSingleNamedReturnByAssign() (a interface{}, s String, err error) {
	a = "" + "1"
	s = "2"
	return a, s, newErr()
}

func FunWithSwitch() (a interface{}, b String) {
	switch a {
	case "1":
		a = "a1"
		b = "b1"
		return
	case "2":
		a = "a2"
		b = "b2"
		return
	default:
		a = "a3"
		b = "b3"
	}
	return
}

func str(a string, b string) string {
	return a + b
}

func FuncWithIf() (a interface{}) {
	if true {
		a = "a0"
		return
	} else if true {
		a = "a1"
		return
	} else {
		a = str("a", "b")
		return
	}
}

func FuncCallReturnAssign() (a interface{}, b String) {
	return FuncSingleReturn(), String(FuncSelectExprReturn())
}

func FuncCallWithFuncLit() (a interface{}, b String) {
	call := func() interface{} {
		return 1
	}
	return call(), "s"
}

func FuncWithImportedCall() interface{} {
	return b.V()
}

type Func func() func() int

func curryCall(r bool) Func {
	if r {
		return func() func() int {
			return func() int {
				return 1
			}
		}
	}

	return func() func() int {
		return b.V
	}
}

func FuncCurryCall() interface{} {
	return curryCall(true)()()
}
