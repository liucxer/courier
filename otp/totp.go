package otp

import "time"

func NewDefaultTOTP(secret string) *TOTP {
	return NewTOTP(secret, 6, 30, nil)
}

func NewTOTP(secret string, digits int, interval int, hasher *Hasher) *TOTP {
	otp := NewOTP(secret, digits, hasher)
	return &TOTP{OTP: otp, interval: interval}
}

type TOTP struct {
	*OTP
	interval int
}

func (t *TOTP) At(timestamp int64) string {
	return t.GenerateOTP(t.timecode(timestamp))
}

func (t *TOTP) Now() string {
	return t.At(time.Now().Unix())
}

func (t *TOTP) NowWithExpiration() (string, int64) {
	interval64 := int64(t.interval)
	timeCodeInt64 := time.Now().Unix() / interval64
	expirationTime := (timeCodeInt64 + 1) * interval64

	return t.GenerateOTP(int(timeCodeInt64)), expirationTime
}

func (t *TOTP) Verify(otp string, timestamp int64) bool {
	return otp == t.At(timestamp)
}

func (t *TOTP) timecode(timestamp int64) int {
	return int(timestamp / int64(t.interval))
}
