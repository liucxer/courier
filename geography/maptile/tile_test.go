package maptile

import (
	"fmt"
	"github.com/liucxer/courier/geography/coordstransform"
	"testing"

	"github.com/liucxer/courier/geography/encoding/mvt"
	"github.com/liucxer/courier/geography/encoding/mvt/vector_tile"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"

	"github.com/liucxer/courier/geography"
)

func TestTile(t *testing.T) {
	// https://overpass-api.de/api/map?bbox=101.250,21.943,112.500,31.952
	mt := NewMapTile(5, 25, 13)

	min := geography.Point{101.250, 21.94304553343818}
	max := geography.Point{112.5, 31.952162238024968}

	require.Equal(t, geography.Point{0, 4096}, mt.NewTransform(4096)(min))
	require.Equal(t, geography.Point{4096, 0}, mt.NewTransform(4096)(max))

	mt.AddTileLayers(
		LayerPoi{N: "1"},
		LayerPoi{N: "2"},
		LayerPoi{N: "3"},
		LayerPoi{N: "4"},
	)

	v, _ := mvt.ToMVT(mt)
	tile := vector_tile.Tile{}
	if err := proto.Unmarshal(v.Bytes(), &tile); err != nil {
		panic(err)
	}

	for i := range tile.Layers {
		layer := tile.Layers[i]

		fmt.Printf("%s\n", *layer.Name)
		fmt.Printf("\t%d\n", *layer.Version)

		fmt.Printf("\t%v\n", layer.Keys)
		fmt.Printf("\t%v\n", layer.Values)

		for j := range layer.Features {
			f := layer.Features[j]
			fmt.Printf("\t\t%v\n", f.Type)
			fmt.Printf("\t\t%v\n", f.Geometry)
		}
	}

	{
		mt.SetCoordsTransform(coordstransform.CoordsTransform{})
		v, _ := mvt.ToMVT(mt)
		tile := vector_tile.Tile{}
		if err := proto.Unmarshal(v.Bytes(), &tile); err != nil {
			panic(err)
		}

		for i := range tile.Layers {
			layer := tile.Layers[i]

			fmt.Printf("%s\n", *layer.Name)
			fmt.Printf("\t%d\n", *layer.Version)

			fmt.Printf("\t%v\n", layer.Keys)
			fmt.Printf("\t%v\n", layer.Values)

			for j := range layer.Features {
				f := layer.Features[j]
				fmt.Printf("\t\t%v\n", f.Type)
				fmt.Printf("\t\t%v\n", f.Geometry)
			}
		}
	}
}

type LayerPoi struct {
	N string
}

func (p LayerPoi) Name() string {
	return "poi" + p.N
}

func (LayerPoi) Fields() map[string]FieldType {
	return map[string]FieldType{
		"name": FieldTypeString,
	}
}

func (LayerPoi) Features(tile *MapTile) ([]Feature, error) {
	return []Feature{
		FeaturePoi{Geom: geography.Point{110, 20}},
		FeaturePoi{Geom: geography.LineString{{110, 22}, {111, 23}}},
		FeaturePoi{Geom: geography.Polygon{{{110, 24}, {110, 24}, {110, 24}}}},
	}, nil
}

type FeaturePoi struct {
	geography.Geom
}

func (FeaturePoi) ID() uint64  {
	return 1
}

func (w FeaturePoi) ToGeom() geography.Geom {
	return w.Geom
}

func (FeaturePoi) Properties() map[string]interface{} {
	return map[string]interface{}{
		"name": "string",
	}
}
