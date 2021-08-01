package geography

import (
	"testing"

	"github.com/liucxer/courier/geography/encoding/wkt"
	"github.com/stretchr/testify/require"
)

func TestMultiPoint(t *testing.T) {
	p := MultiPoint{{0, 0}, {1, 1}}

	require.Equal(t, "MultiPoint(0 0,1 1)", wkt.MarshalWKT(p, 0))
	require.True(t, MultiPoint{{0, 0}, {1, 1}}.Equal(MultiPoint{{0, 0}, {1, 1}}))
}
