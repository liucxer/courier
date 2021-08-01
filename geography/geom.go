package geography

import (
	"github.com/liucxer/courier/geography/encoding/wkb"
)

type Transform func(point Point) Point

type Geom interface {
	Clip(b Bound) Geom
	Project(transform Transform) Geom

	Type() string

	ToGeom() Geom
	Geometry() []uint32

	Bound() Bound
	Equal(g Geom) bool
}

func scan(src interface{}, g Geom) error {
	if src == nil {
		return nil
	}
	if data, ok := src.([]byte); ok {
		if err := wkb.UnmarshalWKB(data, g); err != nil {
			return err
		}
	}
	return nil
}
