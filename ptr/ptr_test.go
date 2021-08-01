package ptr

import (
	testing "testing"

	github_com_onsi_gomega "github.com/onsi/gomega"
)

func TestBool(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(bool(true)).To(github_com_onsi_gomega.Equal(*Bool(true)))
	github_com_onsi_gomega.NewWithT(t).Expect(bool(false)).To(github_com_onsi_gomega.Equal(*Bool(false)))
}

func TestInt(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(int(1)).To(github_com_onsi_gomega.Equal(*Int(1)))
}

func TestInt8(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(int8(1)).To(github_com_onsi_gomega.Equal(*Int8(1)))
}

func TestInt16(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(int16(1)).To(github_com_onsi_gomega.Equal(*Int16(1)))
}

func TestInt32(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(int32(1)).To(github_com_onsi_gomega.Equal(*Int32(1)))
}

func TestInt64(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(int64(1)).To(github_com_onsi_gomega.Equal(*Int64(1)))
}

func TestUint(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(uint(1)).To(github_com_onsi_gomega.Equal(*Uint(1)))
}

func TestUint8(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(uint8(1)).To(github_com_onsi_gomega.Equal(*Uint8(1)))
}

func TestUint16(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(uint16(1)).To(github_com_onsi_gomega.Equal(*Uint16(1)))
}

func TestUint32(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(uint32(1)).To(github_com_onsi_gomega.Equal(*Uint32(1)))
}

func TestUint64(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(uint64(1)).To(github_com_onsi_gomega.Equal(*Uint64(1)))
}

func TestUintptr(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(uintptr(1)).To(github_com_onsi_gomega.Equal(*Uintptr(1)))
}

func TestFloat32(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(float32(1)).To(github_com_onsi_gomega.Equal(*Float32(1)))
}

func TestFloat64(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(float64(1)).To(github_com_onsi_gomega.Equal(*Float64(1)))
}

func TestComplex64(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(complex64(1)).To(github_com_onsi_gomega.Equal(*Complex64(1)))
}

func TestComplex128(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(complex128(1)).To(github_com_onsi_gomega.Equal(*Complex128(1)))
}

func TestString(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(string("string")).To(github_com_onsi_gomega.Equal(*String("string")))
}

func TestByte(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(byte([]uint8{
		98,
		121,
		116,
		101,
		115,
	}[0])).To(github_com_onsi_gomega.Equal(*Byte([]uint8{
		98,
		121,
		116,
		101,
		115,
	}[0])))
}

func TestRune(t *testing.T) {
	github_com_onsi_gomega.NewWithT(t).Expect(rune('r')).To(github_com_onsi_gomega.Equal(*Rune('r')))
}
