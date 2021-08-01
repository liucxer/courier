package geography

import (
	"testing"

	"github.com/liucxer/courier/geography/encoding/wkt"
	"github.com/stretchr/testify/require"
)

func TestMultiPolygon(t *testing.T) {
	p := MultiPolygon{
		{
			{{0, 0}, {1, 1}, {2, 2}, {0, 0}},
		},
		{
			{{0, 0}, {1, 1}, {2, 2}, {0, 0}},
		},
	}

	require.Equal(t, "MultiPolygon(((0 0,1 1,2 2,0 0)),((0 0,1 1,2 2,0 0)))", wkt.MarshalWKT(p, 0))
	require.True(t, p.Equal(MultiPolygon{
		{
			{{0, 0}, {1, 1}, {2, 2}, {0, 0}},
		},
		{
			{{0, 0}, {1, 1}, {2, 2}, {0, 0}},
		},
	}))
}
