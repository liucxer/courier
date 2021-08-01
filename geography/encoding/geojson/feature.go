package geojson

import (
	"encoding/json"
	"errors"
)

type Feature struct {
	ID         interface{}            `json:"id,omitempty"`
	Type       string                 `json:"type"`
	Geometry   *Geometry              `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
	CRS        map[string]interface{} `json:"crs,omitempty"`
}

func (f *Feature) MarshalJSON() ([]byte, error) {
	type feature Feature
	fea := &feature{
		ID:       f.ID,
		Type:     "Feature",
		Geometry: f.Geometry,
	}

	if f.Properties != nil && len(f.Properties) != 0 {
		fea.Properties = f.Properties
	}

	if f.CRS != nil && len(f.CRS) != 0 {
		fea.CRS = f.CRS
	}

	return json.Marshal(fea)
}

func (f Feature) MarshalText() ([]byte, error) {
	return f.MarshalJSON()
}

func (f *Feature) UnmarshalText(data []byte) error {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return err
	}
	return decodeFeature(f, object)
}

func decodeFeature(f *Feature, object map[string]interface{}) error {
	t, ok := object["type"]
	if !ok {
		return errors.New("type property not defined")
	}

	if str, ok := t.(string); ok {
		f.Type = str
	} else {
		return errors.New("type property not string")
	}

	properties, ok := object["properties"]
	if !ok {
		return errors.New("properties property not defined")
	}
	if p, ok := properties.(map[string]interface{}); ok {
		f.Properties = p
	}

	crs, ok := object["crs"]
	if ok {
		if c, ok := crs.(map[string]interface{}); ok {
			f.CRS = c
		}
	}

	geometry, ok := object["geometry"]
	if !ok {
		return errors.New("geometry property not defined")
	}
	if geo, ok := geometry.(map[string]interface{}); ok {
		var g = &Geometry{}
		err := decodeGeometry(g, geo)
		if err != nil {
			return err
		}
		f.Geometry = g
	} else {
		return errors.New("geometry property not geometry")
	}
	return nil
}
