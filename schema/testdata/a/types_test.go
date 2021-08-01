package a

import (
	"context"
	"reflect"
	"testing"

	"github.com/liucxer/courier/reflectx/typesutil"
	"github.com/liucxer/courier/schema/pkg/jsonschema/extractors"
	"github.com/liucxer/courier/schema/pkg/ptr"
	"github.com/liucxer/courier/schema/pkg/testutil"
	"github.com/liucxer/courier/schema/pkg/validator"
	"github.com/liucxer/courier/schema/testdata/b"
)

type Struct2 Struct

func TestSchemaFrom(t *testing.T) {
	testutil.PrintJSON(extractors.SchemaFrom(&Struct{}))
	testutil.PrintJSON(extractors.SchemaFrom(&Struct2{}))
}

func TestValidate(t *testing.T) {
	vf := validator.ValidatorMgrDefault.MustCompile(
		validator.ContextWithNamedTagKey(context.Background(), "json"),
		nil,
		typesutil.FromRType(reflect.TypeOf(&Struct{})),
		nil,
	)

	s := Struct{
		Name:       ptr.String("xx"),
		PullPolicy: b.PullAlways,
		Protocol:   PROTOCOL__HTTP,
		Slice:      []int{3},
		Map: map[string]map[string]struct {
			ID int `json:"id" validate:"@int[0,10]"`
		}{
			"a": {"b": {ID: 3}},
		},
	}

	{
		err := vf.Validate(s)
		t.Log(err)
	}

	{
		err := s.Validate()
		t.Log(err)
	}
}

func BenchmarkValidate(b *testing.B) {
	v := validator.ValidatorMgrDefault.MustCompile(
		validator.ContextWithNamedTagKey(context.Background(), "json"),
		nil,
		typesutil.FromRType(reflect.TypeOf(&Struct{})),
		nil,
	)

	b.Run("full errors", func(b *testing.B) {
		s := Struct{
			Slice:    []int{11},
			Int:      1024,
			Name:     ptr.String("x"),
			Protocol: 10,
			Map: map[string]map[string]struct {
				ID int `json:"id" validate:"@int[0,10]"`
			}{
				"a": {"b": {ID: 11}},
			},
		}

		b.Run("by reflect", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = v.Validate(s)
			}
		})

		b.Run("by generated", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = s.Validate()
			}
		})
	})

	b.Run("no errors", func(b *testing.B) {
		s := Struct{
			Name:     ptr.String("xx"),
			Protocol: PROTOCOL__HTTP,
			Slice:    []int{3},
			Map: map[string]map[string]struct {
				ID int `json:"id" validate:"@int[0,10]"`
			}{
				"a": {"b": {ID: 3}},
			},
		}

		b.Run("by reflect", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = v.Validate(s)
			}
		})

		b.Run("by generated", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = s.Validate()
			}
		})
	})
}
