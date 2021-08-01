package gengo

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/liucxer/courier/gengo/pkg/namer"
	"github.com/onsi/gomega"
)

func TestDumper_TypeLit(t *testing.T) {
	tt := gomega.NewWithT(t)

	d := NewDumper("", namer.NewDefaultImportTracker())

	t.Run("TypeLit", func(t *testing.T) {
		tt.Expect("*bytes.Buffer").To(gomega.Equal(d.ReflectTypeLit(reflect.TypeOf(&bytes.Buffer{}))))
		tt.Expect("[]string").To(gomega.Equal(d.ReflectTypeLit(reflect.TypeOf([]string{}))))
		tt.Expect("map[string]string").To(gomega.Equal(d.ReflectTypeLit(reflect.TypeOf(map[string]string{}))))
		tt.Expect("*struct {V string `json:\"v\" validate:\"@int[0,10]\"`\n}").To(gomega.Equal(d.ReflectTypeLit(reflect.TypeOf(&struct {
			V string `json:"v" validate:"@int[0,10]"`
		}{}))))
	})

	t.Run("ValueLit", func(t *testing.T) {
		tt.Expect("&(bytes.Buffer{})").To(gomega.Equal(d.ValueLit(reflect.ValueOf(&(bytes.Buffer{})))))
		tt.Expect(`[]string{
"1",
"2",
}`).To(gomega.Equal(d.ValueLit(reflect.ValueOf([]string{"1", "2"}))))
	})
}
