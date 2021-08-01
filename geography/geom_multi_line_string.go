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

type MultiLineString []LineString

func (mls MultiLineString) ToGeom() Geom {
	return mls
}

func (mls MultiLineString) Clip(b Bound) Geom {
	var result MultiLineString
	for i := range mls {
		result = append(result, mls[i].Clip(b).(MultiLineString)...)
	}
	return result
}

func (mls MultiLineString) Project(transform Transform) Geom {
	nextMls := make(MultiLineString, len(mls))
	for i := range mls {
		nextMls[i] = mls[i].Project(transform).(LineString)
	}
	return nextMls
}

func (MultiLineString) Type() string {
	return "MultiLineString"
}

func (mls MultiLineString) Bound() Bound {
	if len(mls) == 0 {
		return emptyBound
	}

	bound := mls[0].Bound()
	for i := 1; i < len(mls); i++ {
		bound = bound.Union(mls[i].Bound())
	}

	return bound
}

func (mls MultiLineString) Equal(g Geom) bool {
	switch multiLineString := g.(type) {
	case MultiLineString:
		if len(mls) != len(multiLineString) {
			return false
		}

		for i, ls := range mls {
			if !ls.Equal(multiLineString[i]) {
				return false
			}
		}

		return true
	}
	return false
}

func (mls MultiLineString) MarshalWKT(w *wkt.WKTWriter) {
	for i, line := range mls {
		w.WriteGroup(line.MarshalWKT, i)
	}
}

func (mls *MultiLineString) UnmarshalWKB(r *wkb.WKBReader, order binary.ByteOrder, tpe wkb.GeometryType) error {
	if tpe != wkb.MultiLineStringType {
		return errors.New("not multi line string wkb")
	}

	var numOfLineStrings uint32
	if err := r.ReadBinary(order, &numOfLineStrings); err != nil {
		return err
	}

	result := make(MultiLineString, 0, numOfLineStrings)

	for i := 0; i < int(numOfLineStrings); i++ {
		p := LineString{}
		if err := r.ReadWKB(&p); err != nil {
			return err
		}
		result = append(result, p)
	}

	*mls = result
	return nil
}

func (mls MultiLineString) Cap() int {
	c := 0
	for _, ls := range mls {
		c += ls.Cap()
	}
	return c
}

func (mls MultiLineString) DrawFeature(w *mvt.FeatureWriter) {
	for _, ls := range mls {
		ls.DrawFeature(w)
	}
}

func (mls MultiLineString) Geometry() []uint32 {
	w := mvt.NewFeatureWriter(mls.Cap())
	mls.DrawFeature(w)
	return w.Data()
}

func (MultiLineString) DataType(driverName string) string {
	if driverName == "mysql" {
		return "MULTILINESTRING"
	}
	return "geometry(MULTILINESTRING)"
}

func (mls MultiLineString) Value() (driver.Value, error) {
	return wkt.MarshalWKT(mls, SRS3857), nil
}

func (MultiLineString) ValueEx() string {
	return "ST_GeomFromText(?," + strconv.FormatInt(SRS3857, 10) + ")"
}

func (mls *MultiLineString) Scan(src interface{}) error {
	return scan(src, mls)
}
