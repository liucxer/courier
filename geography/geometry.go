package geography

import (
	"database/sql/driver"
	"encoding/binary"
	"strconv"

	"github.com/liucxer/courier/geography/encoding/wkb"
	"github.com/liucxer/courier/geography/encoding/wkt"
)

func ToGeometry(g Geom) Geometry {
	return Geometry{Geom: g}
}

type Geometry struct {
	Geom
}

func (g Geometry) ToGeom() Geom {
	return g.Geom
}

func (g Geometry) MarshalWKT(w *wkt.WKTWriter) {
	if wktMarshaller, ok := g.Geom.(wkt.WKTMarshaller); ok {
		wktMarshaller.MarshalWKT(w)
	}
}

func (g *Geometry) UnmarshalWKB(r *wkb.WKBReader, order binary.ByteOrder, tpe wkb.GeometryType) error {
	switch tpe {
	case wkb.PointType:
		gg := Point{}
		if err := gg.UnmarshalWKB(r, order, tpe); err != nil {
			return err
		}
		g.Geom = gg
	case wkb.LineStringType:
		gg := LineString{}
		if err := gg.UnmarshalWKB(r, order, tpe); err != nil {
			return err
		}
		g.Geom = gg
	case wkb.PolygonType:
		gg := Polygon{}
		if err := gg.UnmarshalWKB(r, order, tpe); err != nil {
			return err
		}
		g.Geom = gg
	case wkb.MultiPointType:
		gg := MultiPoint{}
		if err := gg.UnmarshalWKB(r, order, tpe); err != nil {
			return err
		}
		g.Geom = gg
	case wkb.MultiLineStringType:
		gg := MultiLineString{}
		if err := gg.UnmarshalWKB(r, order, tpe); err != nil {
			return err
		}
		g.Geom = gg
	case wkb.MultiPolygonType:
		gg := MultiPolygon{}
		if err := gg.UnmarshalWKB(r, order, tpe); err != nil {
			return err
		}
		g.Geom = gg
	}
	return nil
}

func (g Geometry) DataType(driverName string) string {
	return "geometry"
}

func (Geometry) ValueEx() string {
	return "ST_GeomFromText(?," + strconv.FormatInt(SRS3857, 10) + ")"
}

func (g Geometry) Value() (driver.Value, error) {
	return wkt.MarshalWKT(g, SRS3857), nil
}

func (g *Geometry) Scan(src interface{}) error {
	return scan(src, g)
}
