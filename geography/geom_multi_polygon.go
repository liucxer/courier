package geography

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"strconv"

	"github.com/liucxer/courier/geography/encoding/mvt"
	"github.com/liucxer/courier/geography/encoding/wkb"
	"github.com/liucxer/courier/geography/encoding/wkt"
)

type MultiPolygon []Polygon

func (mp MultiPolygon) ToGeom() Geom {
	return mp
}

func (mp MultiPolygon) Clip(b Bound) Geom {
	var result MultiPolygon
	for i := range mp {
		g := mp[i].Clip(b)
		if g != nil {
			result = append(result, g.(Polygon))
		}
	}
	return result
}

func (mp MultiPolygon) Project(transform Transform) Geom {
	nextMp := make(MultiPolygon, len(mp))
	for i := range mp {
		nextMp[i] = mp[i].Project(transform).(Polygon)
	}
	return nextMp
}

func (MultiPolygon) Type() string {
	return "MultiPolygon"
}

func (mp MultiPolygon) Bound() Bound {
	if len(mp) == 0 {
		return emptyBound
	}
	bound := mp[0].Bound()
	for i := 1; i < len(mp); i++ {
		bound = bound.Union(mp[i].Bound())
	}

	return bound
}

func (mp MultiPolygon) Equal(g Geom) bool {
	switch multiPolygon := g.(type) {
	case MultiPolygon:
		if len(mp) != len(multiPolygon) {
			return false
		}

		for i, p := range mp {
			if !p.Equal(multiPolygon[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (mp MultiPolygon) MarshalWKT(w *wkt.WKTWriter) {
	for i, polygon := range mp {
		w.WriteGroup(polygon.MarshalWKT, i)
	}
}

func (mp *MultiPolygon) UnmarshalWKB(r *wkb.WKBReader, order binary.ByteOrder, tpe wkb.GeometryType) error {
	if tpe != wkb.MultiPolygonType {
		return errors.New("not multi line string wkb")
	}

	var numOfPolygons uint32
	if err := r.ReadBinary(order, &numOfPolygons); err != nil {
		return err
	}

	result := make(MultiPolygon, 0, numOfPolygons)

	for i := 0; i < int(numOfPolygons); i++ {
		p := Polygon{}
		if err := r.ReadWKB(&p); err != nil {
			return err
		}
		result = append(result, p)
	}

	*mp = result
	return nil
}

func (mp MultiPolygon) Cap() int {
	c := 0
	for _, p := range mp {
		c += p.Cap()
	}
	return c
}

func (mp MultiPolygon) DrawFeature(w *mvt.FeatureWriter) {
	for _, p := range mp {
		p.DrawFeature(w)
	}
}

func (mp MultiPolygon) Geometry() []uint32 {
	w := mvt.NewFeatureWriter(mp.Cap())
	mp.DrawFeature(w)
	return w.Data()
}

func (MultiPolygon) DataType(driverName string) string {
	if driverName == "mysql" {
		return "MULTIPOLYGON"
	}
	return "geometry(MULTIPOLYGON)"
}

func (MultiPolygon) ValueEx() string {
	return "ST_GeomFromText(?," + strconv.FormatInt(SRS3857, 10) + ")"
}

func (mp MultiPolygon) Value() (driver.Value, error) {
	return wkt.MarshalWKT(mp, SRS3857), nil
}

func (mp *MultiPolygon) Scan(src interface{}) error {
	return scan(src, mp)
}
