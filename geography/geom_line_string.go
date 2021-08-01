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

type LineString []Point

func (ls LineString) ToGeom() Geom {
	return ls
}

func (ls LineString) Clip(b Bound) Geom {
	if !b.Intersects(ls.Bound()) {
		return nil
	}
	return line(b, ls, false)
}

func (ls LineString) Project(transform Transform) Geom {
	return LineString(MultiPoint(ls).Project(transform).(MultiPoint))
}

func (LineString) Type() string {
	return "LineString"
}

func (ls LineString) Closed() bool {
	return ls[0].Equal(ls[len(ls)-1])
}

func (ls LineString) Equal(g Geom) bool {
	switch linestring := g.(type) {
	case LineString:
		return MultiPoint(ls).Equal(MultiPoint(linestring))
	}
	return false
}

func (ls LineString) Bound() Bound {
	return MultiPoint(ls).Bound()
}

func (ls LineString) MarshalWKT(w *wkt.WKTWriter) {
	MultiPoint(ls).MarshalWKT(w)
}

func (ls *LineString) UnmarshalWKB(r *wkb.WKBReader, order binary.ByteOrder, tpe wkb.GeometryType) error {
	if tpe != wkb.LineStringType {
		return errors.New("not line string wkb")
	}

	var numOfPoints uint32
	if err := r.ReadBinary(order, &numOfPoints); err != nil {
		return err
	}

	result := make(LineString, 0, numOfPoints)

	for i := 0; i < int(numOfPoints); i++ {
		p := Point{}
		if err := p.UnmarshalWKB(r, order, wkb.PointType); err != nil {
			return err
		}
		result = append(result, p)
	}

	*ls = result
	return nil
}

func (ls LineString) IsValid() bool {
	n := len(ls)
	if ls.Closed() {
		return n >= 3
	}
	return n >= 2
}

func (ls LineString) DrawFeature(w *mvt.FeatureWriter) {
	if !ls.IsValid() {
		return
	}

	w.MoveTo(1, func(i int) mvt.Coord {
		return ls[0]
	})

	if ls.Closed() {
		points := ls[1 : len(ls)-1]
		w.LineTo(len(points), func(i int) mvt.Coord {
			return points[i]
		})
		w.ClosePath()
	} else {
		points := ls[1:]
		w.LineTo(len(points), func(i int) mvt.Coord {
			return points[i]
		})
	}
}

func (ls LineString) Cap() int {
	return 2 + 2*len(ls)
}

func (ls LineString) Geometry() []uint32 {
	w := mvt.NewFeatureWriter(ls.Cap())
	ls.DrawFeature(w)
	return w.Data()
}

func (LineString) DataType(driverName string) string {
	if driverName == "mysql" {
		return "LINESTRING"
	}
	return "geometry(LINESTRING)"
}

func (LineString) ValueEx() string {
	return "ST_GeomFromText(?," + strconv.FormatInt(SRS3857, 10) + ")"
}

func (ls LineString) Value() (driver.Value, error) {
	return wkt.MarshalWKT(ls, SRS3857), nil
}

func (ls *LineString) Scan(src interface{}) error {
	return scan(src, ls)
}
