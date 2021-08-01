package wkt

import (
	"bytes"
)

func MarshalWKT(v WKTMarshaller, srs uint32) string {
	w := NewWKTWriter(srs)
	w.WriteString(v.Type())
	w.WriteGroup(v.MarshalWKT, 0)
	return w.String()
}

type WKTMarshaller interface {
	Type() string
	MarshalWKT(w *WKTWriter)
}

func NewWKTWriter(srs uint32) *WKTWriter {
	return &WKTWriter{
		SRS: srs,
	}
}

type WKTWriter struct {
	SRS uint32
	bytes.Buffer
}

func (w *WKTWriter) WriteGroup(fn func(w *WKTWriter), idx int) {
	if idx > 0 {
		w.WriteByte(',')
	}
	w.WriteByte('(')
	fn(w)
	w.WriteByte(')')
}
