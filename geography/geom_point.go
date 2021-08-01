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

type Point [2]float64

func (p Point) IsZero() bool {
	return p[0] == 0 && p[1] == 0
}

func (p Point) ToGeom() Geom {
	return p
}

func (p Point) Clip(bound Bound) Geom {
	if !bound.Intersects(p.Bound()) {
		return nil
	}
	return p
}

func (p Point) Project(transform Transform) Geom {
	return transform(p)
}

func (Point) Type() string {
	return "Point"
}

func (p Point) Bound() Bound {
	return Bound{p, p}
}

func (p Point) Equal(geom Geom) bool {
	if point, ok := geom.(Point); ok {
		return p[0] == point[0] && p[1] == point[1]
	}
	return false
}

func (p Point) X() float64 {
	return p[0]
}

func (p Point) Y() float64 {
	return p[1]
}

func (p Point) Lon() float64 {
	return p[0]
}

func (p Point) Lat() float64 {
	return p[1]
}

func (p *Point) UnmarshalWKB(r *wkb.WKBReader, order binary.ByteOrder, tpe wkb.GeometryType) error {
	if tpe != wkb.PointType {
		return errors.New("not point wkb")
	}
	if err := r.ReadBinary(order, &p[0]); err != nil {
		return err
	}
	if err := r.ReadBinary(order, &p[1]); err != nil {
		return err
	}

	switch r.SRS {
	case SRS3857:
		p[0], p[1] = PseudoMercatorToWGS84(p[0], p[1])
	}
	return nil
}

func (Point) DataType(driverName string) string {
	if driverName == "mysql" {
		return "POINT"
	}
	return "geometry(POINT)"
}

func (Point) Cap() int {
	return 3
}

func (p Point) DrawFeature(w *mvt.FeatureWriter) {
	w.MoveTo(1, func(i int) mvt.Coord {
		return p
	})
}

func (p Point) Geometry() []uint32 {
	w := mvt.NewFeatureWriter(p.Cap())
	p.DrawFeature(w)
	return w.Data()
}

func (p Point) MarshalWKT(w *wkt.WKTWriter) {
	x, y := p[0], p[1]
	switch w.SRS {
	case SRS3857:
		x, y = WGS84ToPseudoMercator(x, y)
	}
	fmt.Fprintf(w, "%g %g", x, y)
}

func (Point) ValueEx() string {
	return "ST_GeomFromText(?," + strconv.FormatInt(SRS3857, 10) + ")"
}

func (p Point) Value() (driver.Value, error) {
	return wkt.MarshalWKT(p, SRS3857), nil
}

func (p *Point) Scan(src interface{}) error {
	return scan(src, p)
}
