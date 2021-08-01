package geography

import (
	"testing"

	"github.com/liucxer/courier/geography/encoding/wkt"
	"github.com/stretchr/testify/require"
)

func TestLineString(t *testing.T) {
	p := LineString{{0, 0}, {1, 1}}

	require.Equal(t, "LineString(0 0,1 1)", wkt.MarshalWKT(p, 0))
	require.True(t, LineString{{0, 0}, {1, 1}}.Equal(LineString{{0, 0}, {1, 1}}))
}
