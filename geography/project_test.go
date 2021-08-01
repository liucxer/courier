package geography

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPseudoMercatorToWGS84(t *testing.T) {
	lon, lat := PseudoMercatorToWGS84(6669821.53117059, -10271108.1804419)
	require.Equal(t, fmt.Sprintf("%.10f", 59.9160262380001), fmt.Sprintf("%.10f", lon))
	require.Equal(t, fmt.Sprintf("%.10f", -67.4004859349999), fmt.Sprintf("%.10f", lat))
}

func TestWGS84ToPseudoMercator(t *testing.T) {
	lon, lat := WGS84ToPseudoMercator(59.9160262380001, -67.4004859349999)
	require.Equal(t, fmt.Sprintf("%.7f", 6669821.53117059), fmt.Sprintf("%.7f", lon))
	require.Equal(t, fmt.Sprintf("%.7f", -10271108.1804419), fmt.Sprintf("%.7f", lat))
}
