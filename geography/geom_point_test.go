package geography

import (
	"testing"

	"github.com/liucxer/courier/geography/encoding/wkt"
	"github.com/stretchr/testify/require"
)

func TestPoint(t *testing.T) {
	p := Point{1, 1}

	require.Equal(t, "Point(1 1)", wkt.MarshalWKT(p, 0))
	require.True(t, Point{1, 1}.Equal(Point{1, 1}))
}
