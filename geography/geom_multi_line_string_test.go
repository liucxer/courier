package geography

import (
	"testing"

	"github.com/liucxer/courier/geography/encoding/wkt"
	"github.com/stretchr/testify/require"
)

func TestMultiLineString(t *testing.T) {
	p := MultiLineString{
		{{0, 0}, {1, 1}},
		{{1, 1}, {2, 2}},
	}

	require.True(t, p.Equal(MultiLineString{
		{{0, 0}, {1, 1}},
		{{1, 1}, {2, 2}},
	}))
	require.Equal(t, "MultiLineString((0 0,1 1),(1 1,2 2))", wkt.MarshalWKT(p, 0))
}
