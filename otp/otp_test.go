package otp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOTP(t *testing.T) {
	otp := NewOTP(RandomSecret(16), 9, nil)

	t.Run("#GenerateOTP", func(t *testing.T) {
		cts := int(time.Now().Unix())
		codes := map[string]bool{}

		for i := 0; i < 100; i++ {
			codes[otp.GenerateOTP(cts)] = true
		}

		require.Len(t, codes, 1)
	})
}

func TestOTPType(t *testing.T) {
	t.Run("#OTPAuthURI", func(t *testing.T) {
		s := OTPTypeTOTP.OTPAuthURI("qazwsxedc", 6, OTPAuthOption{
			Issuer:  "X",
			Account: "xxx",
		})

		require.Equal(t, "otpauth://totp/X:xxx?issuer=X&secret=qazwsxedc", s)
	})
}

func TestRandomSecret(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Logf(RandomSecret(32))
	}

}
