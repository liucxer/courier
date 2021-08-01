package geojson_test

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/liucxer/courier/geography"
	"github.com/liucxer/courier/geography/coordstransform"
	"github.com/liucxer/courier/geography/encoding/geojson"
	"github.com/liucxer/courier/geography/maptile"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewFeatureCollection(t *testing.T) {
	fc := geojson.NewFeatureCollection()

	if fc.Type != "FeatureCollection" {
		t.Errorf("should have type of FeatureCollection, got %v", fc.Type)
	}
}

func TestFeatureCollectionToJSON(t *testing.T) {
	fc := geojson.NewFeatureCollection()

	fc.SetCoordsTransform(&coordstransform.CoordsTransform{})

	data, err := fc.AddMapTileFeature([]maptile.Feature{
		&FeaturePoi{Geom: geography.Point{110, 20}},
		&FeaturePoi{Geom: geography.LineString{{110, 22}, {111, 23}}},
		&FeaturePoi{Geom: geography.Polygon{{{110, 24}, {110, 24}, {110, 24}}}},
	}...).ToJSON()

	require.NoError(t, err)

	fmt.Printf("%s\n", string(data))

}

type FeaturePoi struct {
	geography.Geom
}

func (*FeaturePoi) ID() uint64 {
	return 1
}

func (w *FeaturePoi) ToGeom() geography.Geom {
	return w.Geom
}

func (*FeaturePoi) Properties() map[string]interface{} {
	return map[string]interface{}{
		"name": "张三",
		"sex":  "男",
		"age":  11,
	}
}

func TestFeatureCollection_UnmarshalText(t *testing.T) {
	var rawJSON = `
{
    "type": "FeatureCollection",
    "crs": {
        "type": "name",
        "properties": {
            "name": "urn:ogc:def:crs:OGC:1.3:CRS84"
        }
    },
    "features": [{
            "type": "Feature",
            "properties": { "Id": 1234},
            "geometry": {
                "type": "Polygon","coordinates": [[
                        [113.9334716796875,34.87127685546875],
                        [113.9375,34.873901367187543],
                        [113.937683105468764,34.86669921875],
                        [113.939270019531236,34.865905761718722],
                        [113.938476562500043,34.863281250000014]
                    ]
                ]
            },
			"crs":{
				"type":"name"
			}
        }
    ]
}
`
	fc := geojson.NewFeatureCollection()
	err := fc.UnmarshalText([]byte(rawJSON))
	require.NoError(t, err)
	spew.Dump(fc)
	str, err := fc.ToJSON()
	require.NoError(t, err)
	fmt.Println(string(str))
}
