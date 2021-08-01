package otp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTOTP(t *testing.T) {
	totp := NewDefaultTOTP(RandomSecret(10))

	t.Run("#Now", func(t *testing.T) {
		require.Equal(t, totp.At(time.Now().Unix()), totp.Now())
	})

	t.Run("#NowWithExpiration", func(t *testing.T) {
		otp, exp := totp.NowWithExpiration()
		cts := time.Now().Unix()

		require.Equal(t, totp.Now(), otp)
		require.Equal(t, totp.At(exp), totp.At(cts+30))
	})

	t.Run("#Verify", func(t *testing.T) {
		cts := time.Now().Unix()
		otp := totp.At(cts)
		require.True(t, totp.Verify(otp, cts))
	})
}
