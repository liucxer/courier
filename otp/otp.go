package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"hash"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Hasher struct {
	HashName string
	Digest   func() hash.Hash
}

type OTP struct {
	byteSecret []byte  // byteSecret in base32 format
	digits     int     // number of integers in the OTP. Some apps expect this to be 6 digits, others support more.
	hasher     *Hasher // digest function to use in the HMAC (expected to be sha1)
}

func NewOTP(secret string, digits int, hasher *Hasher) *OTP {
	if hasher == nil {
		hasher = &Hasher{
			HashName: "sha1",
			Digest:   sha1.New,
		}
	}

	missingPadding := len(secret) % 8
	if missingPadding != 0 {
		secret = secret + strings.Repeat("=", 8-missingPadding)
	}

	byteSecret, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		panic("decode secret failed")
	}

	return &OTP{
		byteSecret: byteSecret,
		digits:     digits,
		hasher:     hasher,
	}
}

func (o *OTP) GenerateOTP(input int) string {
	if input < 0 {
		panic("input must be positive integer")
	}

	hasher := hmac.New(o.hasher.Digest, o.byteSecret)
	hasher.Write(itob(input))

	hmacHash := hasher.Sum(nil)

	offset := int(hmacHash[len(hmacHash)-1] & 0xf)
	code := ((int(hmacHash[offset]) & 0x7f) << 24) |
		((int(hmacHash[offset+1] & 0xff)) << 16) |
		((int(hmacHash[offset+2] & 0xff)) << 8) |
		(int(hmacHash[offset+3]) & 0xff)

	code = code % int(math.Pow10(o.digits))

	return fmt.Sprintf(fmt.Sprintf("%%0%dd", o.digits), code)
}

func itob(integer int) []byte {
	byteArr := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteArr[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteArr
}

type OTPType string

const (
	OTPTypeTOTP OTPType = "totp"
	OTPTypeHOTP OTPType = "hotp"
)

type OTPAuthOption struct {
	Issuer       string
	Account      string
	InitialCount int
	Period       int
	Algorithm    string
}

func (otpType OTPType) OTPAuthURI(secret string, digits int, opt OTPAuthOption) string {
	urlParams := &url.Values{
		"secret": {secret},
	}

	if otpType == OTPTypeHOTP {
		urlParams.Add("counter", strconv.Itoa(opt.InitialCount))
	}

	label := url.QueryEscape(opt.Account)

	if opt.Issuer != "" {
		issuerNameEscape := url.QueryEscape(opt.Issuer)

		label = issuerNameEscape + ":" + label

		urlParams.Add("issuer", opt.Issuer)
	}

	if opt.Algorithm != "" && opt.Algorithm != "sha1" {
		urlParams.Add("algorithm", opt.Algorithm)
	}

	if digits != 0 && digits != 6 {
		urlParams.Add("digits", strconv.Itoa(digits))
	}

	if opt.Period != 0 && opt.Period != 30 {
		urlParams.Add("period", strconv.Itoa(opt.Period))
	}

	return fmt.Sprintf("otpauth://%s/%s?%s", otpType, label, urlParams.Encode())
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567")

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomSecret(length int) string {
	rand.Seed(time.Now().UnixNano())

	bytes := make([]rune, length)

	for i := range bytes {
		bytes[i] = letterRunes[rnd.Intn(len(letterRunes))]
	}

	return string(bytes)
}
