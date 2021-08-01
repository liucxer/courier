package maptile

import (
	"github.com/liucxer/courier/geography"
)

type Feature interface {
	ToGeom() geography.Geom
	Properties() map[string]interface{}
}

type FeatureID interface {
	ID() uint64
}

type TileLayer interface {
	Name() string
	Fields() map[string]FieldType
	Features(tile *MapTile) ([]Feature, error)
}

type TileLayerExtentConf interface {
	Extent() uint32
}

func NewLayer(name string, extent uint32, features ...Feature) *Layer {
	if extent == 0 {
		extent = 4096
	}
	return &Layer{
		Name:     name,
		Extent:   extent,
		Features: features,
	}
}

type Layer struct {
	Name     string
	Extent   uint32
	Features []Feature
}
