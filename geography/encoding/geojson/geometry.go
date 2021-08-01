package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liucxer/courier/geography"
)

// 几何体关联对象
type Geometry struct {
	Type            string `json:"type"`
	Point           *geography.Point
	MultiPoint      *geography.MultiPoint
	LineString      *geography.LineString
	MultiLineString *geography.MultiLineString
	Polygon         *geography.Polygon
	MultiPolygon    *geography.MultiPolygon
	Geometries      []*Geometry
	CRS             map[string]interface{} `json:"crs,omitempty"`
}

func (g *Geometry) MarshalJSON() ([]byte, error) {
	type geometry struct {
		Type        string      `json:"type"`
		Geometries  interface{} `json:"geometries,omitempty"`
		Coordinates interface{} `json:"coordinates,omitempty"`
	}

	geo := &geometry{
		Type: g.Type,
	}

	switch g.Type {
	case "Point":
		geo.Coordinates = g.Point
	case "MultiPoint":
		geo.Coordinates = g.MultiPoint
	case "LineString":
		geo.Coordinates = g.LineString
	case "MultiLineString":
		geo.Coordinates = g.MultiLineString
	case "Polygon":
		geo.Coordinates = g.Polygon
	case "MultiPolygon":
		geo.Coordinates = g.MultiPolygon
	case "GeometryCollection":
		geo.Geometries = g.Geometries
	}

	return json.Marshal(geo)
}
func (g *Geometry) UnmarshalJSON(data []byte) error {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return decodeGeometry(g, object)
}

func (g Geometry) MarshalText() ([]byte, error) {
	return g.MarshalJSON()
}

func (g *Geometry) UnmarshalText(data []byte) error {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return err
	}
	return decodeGeometry(g, object)
}

func decodeGeometry(g *Geometry, object map[string]interface{}) error {
	t, ok := object["type"]
	if !ok {
		return errors.New("type property not defined")
	}

	if s, ok := t.(string); ok {
		g.Type = s
	} else {
		return errors.New("type property not string")
	}

	switch g.Type {
	case "Point":
		point, _ := decodePoint(object["coordinates"])
		g.Point = &point
	case "MultiPoint":
		point, _ := decodeMultiPoint(object["coordinates"])
		g.MultiPoint = &point
	case "LineString":
		line, _ := decodeLineString(object["coordinates"])
		g.LineString = &line
	case "MultiLineString":
		line, _ := decodeMultiLineString(object["coordinates"])
		g.MultiLineString = &line
	case "Polygon":
		polygon, _ := decodePolygon(object["coordinates"])
		g.Polygon = &polygon
	case "MultiPolygon":
		polygon, _ := decodeMultiPolygon(object["coordinates"])
		g.MultiPolygon = &polygon
	case "GeometryCollection":
		g.Geometries, _ = decodeGeometries(object["geometries"])
	}

	return nil
}
func decodePoint(data interface{}) (geography.Point, error) {
	points, ok := data.([]interface{})
	result := geography.Point{}
	if !ok {
		return result, fmt.Errorf("not a valid points, got %v", data)
	}
	for k, point := range points {
		// 丢弃高度
		if k > 1 {
			continue
		}
		if f, ok := point.(float64); ok {
			result[k] = f
		} else {
			return result, fmt.Errorf("not a valid points, got %v", points)
		}
	}
	return result, nil
}

func decodeMultiPoint(data interface{}) (geography.MultiPoint, error) {
	points, ok := data.([]interface{})
	result := geography.MultiPoint{}

	if !ok {
		return result, fmt.Errorf("not a valid set of points, got %v", data)
	}

	for _, point := range points {
		if p, err := decodePoint(point); err == nil {
			result = append(result, p)
		} else {
			return result, err
		}
	}

	return result, nil
}

func decodeLineString(data interface{}) (geography.LineString, error) {
	lines, ok := data.([]interface{})
	result := geography.LineString{}

	if !ok {
		return result, fmt.Errorf("not a valid path, got %v", data)
	}
	for _, line := range lines {
		if l, err := decodePoint(line); err == nil {
			result = append(result, l)
		} else {
			return result, err
		}
	}
	return result, nil
}

func decodeMultiLineString(data interface{}) (geography.MultiLineString, error) {
	lines, ok := data.([]interface{})
	result := geography.MultiLineString{}

	if !ok {
		return result, fmt.Errorf("not a valid path, got %v", data)
	}
	for _, line := range lines {
		if l, err := decodeLineString(line); err == nil {
			result = append(result, l)
		} else {
			return result, err
		}
	}
	return result, nil
}

func decodePolygon(data interface{}) (geography.Polygon, error) {
	polygons, ok := data.([]interface{})
	result := geography.Polygon{}

	if !ok {
		return result, fmt.Errorf("not a valid path, got %v", data)
	}
	for _, polygon := range polygons {
		if p, err := decodeLineString(polygon); err == nil {
			result = append(result, p)
		} else {
			return result, err
		}
	}
	return result, nil
}

func decodeMultiPolygon(data interface{}) (geography.MultiPolygon, error) {
	polygons, ok := data.([]interface{})
	result := geography.MultiPolygon{}

	if !ok {
		return result, fmt.Errorf("not a valid path, got %v", data)
	}
	for _, polygon := range polygons {
		if p, err := decodePolygon(polygon); err == nil {
			result = append(result, p)
		} else {
			return result, err
		}
	}
	return result, nil
}

func decodeGeometries(data interface{}) ([]*Geometry, error) {
	if vs, ok := data.([]interface{}); ok {
		geometries := make([]*Geometry, 0, len(vs))
		for _, v := range vs {
			g := &Geometry{}

			vmap, ok := v.(map[string]interface{})
			if !ok {
				break
			}

			err := decodeGeometry(g, vmap)
			if err != nil {
				return nil, err
			}

			geometries = append(geometries, g)
		}

		if len(geometries) == len(vs) {
			return geometries, nil
		}
	}

	return nil, fmt.Errorf("not a valid set of geometries, got %v", data)
}

func (g *Geometry) IsPoint() bool {
	return g.Type == "Point"
}

func (g *Geometry) IsMultiPoint() bool {
	return g.Type == "MultiPoint"
}

func (g *Geometry) IsLineString() bool {
	return g.Type == "LineString"
}

func (g *Geometry) IsMultiLineString() bool {
	return g.Type == "MultiLineString"
}

func (g *Geometry) IsPolygon() bool {
	return g.Type == "Polygon"
}

func (g *Geometry) IsMultiPolygon() bool {
	return g.Type == "MultiPolygon"
}

func (g *Geometry) IsCollection() bool {
	return g.Type == "GeometryCollection"
}
