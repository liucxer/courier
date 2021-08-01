package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	require.Equal(t, "github.com/liucxer/courier/courier", pkgImportPathCourier)
	require.Equal(t, "github.com/liucxer/courier/httptransport", pkgImportPathHttpTransport)
	require.Equal(t, "github.com/liucxer/courier/httptransport/httpx", pkgImportPathHttpx)
}
