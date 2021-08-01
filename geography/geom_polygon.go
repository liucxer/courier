package geography

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

	"github.com/liucxer/courier/geography/encoding/mvt"
	"github.com/liucxer/courier/geography/encoding/wkb"
	"github.com/liucxer/courier/geography/encoding/wkt"
)

// Polygon is a closed area.
// The first Polygon is the outer ring.
// The others are the holes.
// Each Polygon is expected to be closed
// ie. the first point matches the last.
type Polygon []LineString

func (p Polygon) ToGeom() Geom {
	return p
}

func (p Polygon) Clip(b Bound) Geom {
	if len(p) == 0 {
		return nil
	}

	circle := ring(b, p[0])
	if circle == nil {
		return nil
	}

	result := Polygon{circle}
	for i := 1; i < len(p); i++ {
		r := ring(b, p[i])
		if r != nil {
			result = append(result, r)
		}
	}

	return result
}

func (p Polygon) Project(transform Transform) Geom {
	nextP := make(Polygon, len(p))
	for i := range p {
		nextP[i] = p[i].Project(transform).(LineString)
	}
	return nextP
}

func (p Polygon) Bound() Bound {
	if len(p) == 0 {
		return emptyBound
	}
	return p[0].Bound()
}

func (p Polygon) Equal(g Geom) bool {
	switch polygon := g.(type) {
	case Polygon:
		if len(p) != len(polygon) {
			return false
		}
		for i := range p {
			if !p[i].Equal(polygon[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (Polygon) Type() string {
	return "Polygon"
}

func (p Polygon) MarshalWKT(w *wkt.WKTWriter) {
	MultiLineString(p).MarshalWKT(w)
}

func (p *Polygon) UnmarshalWKB(r *wkb.WKBReader, order binary.ByteOrder, tpe wkb.GeometryType) error {
	if tpe != wkb.PolygonType {
		return errors.New("not line polygon wkb")
	}

	var numOfLines uint32
	if err := r.ReadBinary(order, &numOfLines); err != nil {
		return err
	}

	result := make(Polygon, 0, numOfLines)

	for i := 0; i < int(numOfLines); i++ {
		p := LineString{}
		if err := p.UnmarshalWKB(r, order, wkb.LineStringType); err != nil {
			return fmt.Errorf("error on %d of %s: %s", i, p.Type(), err)
		}
		result = append(result, p)
	}

	*p = result
	return nil
}

func (p Polygon) Cap() int {
	c := 0
	for _, r := range p {
		c += 3 + 2*len(r)
	}
	return c
}

func (p Polygon) DrawFeature(w *mvt.FeatureWriter) {
	for _, ls := range p {
		ls.DrawFeature(w)
		if !ls.Closed() && ls.IsValid() {
			// force close path
			w.ClosePath()
		}
	}
}

func (p Polygon) Geometry() []uint32 {
	w := mvt.NewFeatureWriter(p.Cap())
	p.DrawFeature(w)
	return w.Data()
}

func (Polygon) DataType(driverName string) string {
	if driverName == "mysql" {
		return "POLYGON"
	}
	return "geometry(POLYGON)"
}

func (Polygon) ValueEx() string {
	return "ST_GeomFromText(?," + strconv.FormatInt(SRS3857, 10) + ")"
}

func (p Polygon) Value() (driver.Value, error) {
	return wkt.MarshalWKT(p, SRS3857), nil
}

func (p *Polygon) Scan(src interface{}) error {
	return scan(src, p)
}
