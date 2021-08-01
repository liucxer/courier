package wkb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
)

var (
	ErrNotWKB = errors.New("wkb: invalid data")
)

type GeometryType uint32

const (
	PointType GeometryType = iota + 1
	LineStringType
	PolygonType
	MultiPointType
	MultiLineStringType
	MultiPolygonType
)

type WKBUnmarshaler interface {
	UnmarshalWKB(r *WKBReader, order binary.ByteOrder, tpe GeometryType) error
}

func NewWKBReader(data []byte) *WKBReader {
	return &WKBReader{
		Reader: bytes.NewReader(data),
	}
}

func UnmarshalWKB(data []byte, v interface{}) error {
	srid := uint32(0)
	src := data
	if n, err := hex.Decode(src, src); err == nil {
		// postgres ewkb
		// 01 01_00_00_20 11_0f_00_00 31_0c_4a_da_77_2d_fb_c0 48_f4_04_ab_e1_2e_0b_41
		// o  ewkbtype    srid        x                       y
		data = src[:n]
	} else {
		// mysql srid + wkb
		// 11_0f_00_00 01 01_00_00_00 31_0c_4a_da_77_2d_fb_c0 48_f4_04_ab_e1_2e_0b_41
		// srid        o  type        x                       y
		if err := binary.Read(bytes.NewBuffer(data[0:4]), binary.LittleEndian, &srid); err != nil {
			return err
		}
		data = data[4:]
	}
	r := NewWKBReader(data)
	r.SRS = srid
	if err := r.ReadWKB(v); err != nil {
		return err
	}
	return nil
}

type WKBReader struct {
	SRS uint32
	io.Reader
}

func (r *WKBReader) ReadBinary(order binary.ByteOrder, data interface{}) error {
	return binary.Read(r, order, data)
}

const (
	ewkbZ    = 0x80000000
	ewkbM    = 0x40000000
	ewkbSRID = 0x20000000
)

func (r *WKBReader) ReadWKB(v interface{}) error {
	var bom = make([]byte, 1)

	if _, err := r.Read(bom); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder

	if bom[0] == 0 {
		byteOrder = binary.BigEndian
	} else if bom[0] == 1 {
		byteOrder = binary.LittleEndian
	} else {
		return ErrNotWKB
	}

	var typ uint32
	if err := r.ReadBinary(byteOrder, &typ); err != nil {
		return err
	}

	// ewkb
	if typ > 10 {
		// srid
		var srid uint32
		if err := r.ReadBinary(byteOrder, &srid); err != nil {
			return err
		}
		r.SRS = srid
		typ = typ &^ (ewkbZ | ewkbM | ewkbSRID)
	}

	if wkbUnmarshaler, ok := v.(WKBUnmarshaler); ok {
		if err := wkbUnmarshaler.UnmarshalWKB(r, byteOrder, GeometryType(typ)); err != nil {
			return err
		}
	}
	return nil
}
