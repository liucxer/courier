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

type MultiPoint []Point

func (mp MultiPoint) ToGeom() Geom {
	return mp
}

func (mp MultiPoint) Clip(b Bound) Geom {
	if !b.Intersects(mp.Bound()) {
		return nil
	}

	var result MultiPoint
	for _, p := range mp {
		if b.Contains(p) {
			result = append(result, p)
		}
	}
	return result
}

func (mp MultiPoint) Project(transform Transform) Geom {
	nextMp := make(MultiPoint, len(mp))
	for i := range mp {
		nextMp[i] = transform(mp[i])
	}
	return nextMp
}

func (MultiPoint) Type() string {
	return "MultiPoint"
}

func (mp MultiPoint) Equal(g Geom) bool {
	switch multiPoint := g.(type) {
	case MultiPoint:
		if len(mp) != len(multiPoint) {
			return false
		}
		for i := range mp {
			if !mp[i].Equal(multiPoint[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (mp MultiPoint) Bound() Bound {
	if len(mp) == 0 {
		return emptyBound
	}

	b := Bound{mp[0], mp[0]}
	for _, p := range mp {
		b = b.Extend(p)
	}
	return b
}

func (mp MultiPoint) MarshalWKT(w *wkt.WKTWriter) {
	for i, p := range mp {
		if i > 0 {
			w.WriteByte(',')
		}
		p.MarshalWKT(w)
	}
}

func (mp *MultiPoint) UnmarshalWKB(r *wkb.WKBReader, order binary.ByteOrder, tpe wkb.GeometryType) error {
	if tpe != wkb.MultiPointType {
		return errors.New("not multi point wkb")
	}

	var numOfPoints uint32
	if err := r.ReadBinary(order, &numOfPoints); err != nil {
		return err
	}

	result := make(MultiPoint, 0, numOfPoints)

	for i := 0; i < int(numOfPoints); i++ {
		p := Point{}
		if err := r.ReadWKB(&p); err != nil {
			return err
		}
		result = append(result, p)
	}

	*mp = result
	return nil
}

func (mp MultiPoint) Cap() int {
	return 1 + 2*len(mp)
}

func (mp MultiPoint) DrawFeature(w *mvt.FeatureWriter) {
	w.MoveTo(len(mp), func(i int) mvt.Coord {
		return mp[i]
	})
}

func (mp MultiPoint) Geometry() []uint32 {
	w := mvt.NewFeatureWriter(mp.Cap())
	mp.DrawFeature(w)
	return w.Data()
}

func (MultiPoint) DataType(driverName string) string {
	if driverName == "mysql" {
		return "MULTIPOINT"
	}
	return "geometry(MULTIPOINT)"
}

func (MultiPoint) ValueEx() string {
	return "ST_GeomFromText(?," + strconv.FormatInt(SRS3857, 10) + ")"
}

func (mp MultiPoint) Value() (driver.Value, error) {
	return wkt.MarshalWKT(mp, SRS3857), nil
}

func (mp *MultiPoint) Scan(src interface{}) error {
	return scan(src, mp)
}
