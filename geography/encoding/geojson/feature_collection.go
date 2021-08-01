package geojson

import (
	"encoding/json"
	"errors"

	"github.com/liucxer/courier/geography"
	"github.com/liucxer/courier/geography/coordstransform"
	"github.com/liucxer/courier/geography/maptile"
)

type FeatureCollection struct {
	coordsTransform *coordstransform.CoordsTransform
	Type            string                 `json:"type"`
	Features        []*Feature             `json:"features"`
	CRS             map[string]interface{} `json:"crs,omitempty"`
}

// New FeatureCollection
func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]*Feature, 0),
	}
}

func (fc *FeatureCollection) SetCoordsTransform(coordsTransform *coordstransform.CoordsTransform) {
	fc.coordsTransform = coordsTransform
}

func (fc *FeatureCollection) AddMapTileFeature(features ...maptile.Feature) *FeatureCollection {
	for _, v := range features {
		fc.addMapTileFeature(v)
	}
	return fc
}

func (fc *FeatureCollection) addMapTileFeature(feature maptile.Feature) *FeatureCollection {
	feat := feature.ToGeom()
	geo := &Geometry{
		Type: feat.Type(),
	}

	if fc.coordsTransform != nil {
		feat = feat.Project(fc.coordsTransform.ToMars)
	}

	switch feat.Type() {
	case "Point":
		point, _ := feat.(geography.Point)
		geo.Point = &point
		break
	case "MultiPoint":
		point, _ := feat.(geography.MultiPoint)
		geo.MultiPoint = &point
		break
	case "LineString":
		line, _ := feat.(geography.LineString)
		geo.LineString = &line
		break
	case "MultiLineString":
		line, _ := feat.(geography.MultiLineString)
		geo.MultiLineString = &line
		break
	case "Polygon":
		polygon, _ := feat.(geography.Polygon)
		geo.Polygon = &polygon
		break
	case "MultiPolygon":
		polygon, _ := feat.(geography.MultiPolygon)
		geo.MultiPolygon = &polygon
		break
	}

	fe := &Feature{
		Type:       "Feature",
		Geometry:   geo,
		Properties: feature.Properties(),
	}

	if fid, ok := feature.(interface {
		ID() uint64
	}); ok {
		fe.ID = fid.ID()
	}

	fc.Features = append(fc.Features, fe)
	return fc
}

// MarshalJSON
func (fc *FeatureCollection) MarshalJSON() ([]byte, error) {
	type featureCollection FeatureCollection

	fcol := &featureCollection{
		Type: "FeatureCollection",
	}

	fcol.Features = fc.Features
	if fcol.Features == nil {
		fcol.Features = make([]*Feature, 0)
	}
	if fc.CRS != nil && len(fc.CRS) != 0 {
		fcol.CRS = fc.CRS
	}

	return json.Marshal(fcol)
}

func (fc *FeatureCollection) ToJSON() ([]byte, error) {
	return fc.MarshalJSON()
}

func (fc FeatureCollection) MarshalText() ([]byte, error) {
	return fc.ToJSON()
}

func (fc *FeatureCollection) UnmarshalText(data []byte) error {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return decodeFeatureCollection(fc, object)
}

func decodeFeatureCollection(fc *FeatureCollection, object map[string]interface{}) error {
	t, ok := object["type"]
	if !ok {
		return errors.New("type property not defined")
	}

	if str, ok := t.(string); ok {
		fc.Type = str
	} else {
		return errors.New("type property not string")
	}

	crs, ok := object["crs"]
	if ok {
		if c, ok := crs.(map[string]interface{}); ok {
			fc.CRS = c
		}
	}

	features, ok := object["features"]
	if !ok {
		return errors.New("features property not defined")
	}

	feas, ok := features.([]interface{})
	if !ok {
		return errors.New("type property not features")
	}
	for _, fea := range feas {
		if f, ok := fea.(map[string]interface{}); ok {
			fea := &Feature{}
			err := decodeFeature(fea, f)
			if err != nil {
				return err
			}
			fc.Features = append(fc.Features, fea)
		} else {
			return errors.New("type property not features")
		}
	}
	return nil
}
