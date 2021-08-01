package a

import (
	"time"

	"github.com/liucxer/courier/schema/testdata/b"
)

type Interface interface{}

// String
type String string

// Bool
type Bool bool

// Float
type Float float32

// Double
type Double float64

// Int
type Int int

// Uint
type Uint uint

// FakeBool
// +gengo:jsonschema=false
type FakeBool int

func (FakeBool) OpenAPISchemaType() []string { return []string{"boolean"} }

// Map
type Map map[string]String

// ArrayString
type ArrayString [2]string

// SliceString
type SliceString []string

// SliceNamed
type SliceNamed []String

type TimeAlias = time.Time

// +gengo:validator
// Struct
type Struct struct {
	private string

	Int int `json:"int,omitempty" default:"7" validate:"@int(,1024)"`
	// name
	Name *string `json:"name" validate:"@string[2,]"`
	// id
	ID **string `json:"id,omitempty" default:"11" validate:"@string/\\d+/"`
	// pull policy
	PullPolicy b.PullPolicy `json:"pullPolicy" validate:"@string{Always}"`
	Protocol   Protocol     `json:"protocol" validate:"@string{HTTP}"`
	Slice      []int        `json:"slice" validate:"@slice<@int[0,10]>[1,30]"`
	Map        map[string]map[string]struct {
		ID int `json:"id" validate:"@int[0,10]"`
	} `json:"map,omitempty" validate:"@map<,@map<,@struct>>[0,3]"`
}

type Part struct {
	Name string `json:",omitempty" validate:"@string[2,]"`
	Skip string `json:"-"`
	skip string
}

type PartConflict struct {
	Name string `json:"name" validate:"@string[0,]"`
}

type Composed struct {
	Part
}

type NamedComposed struct {
	Part `json:"part"`
}

type InvalidComposed struct {
	Part
	PartConflict
}

type Node struct {
	Type     string  `json:"type"`
	Children []*Node `json:"children"`
}
