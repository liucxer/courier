package raw

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega"
)

var subCases = [][]interface{}{
	{1, 2, int64(1)},
	{1, uint(2), int64(1)},
	{uint(1), uint(2), uint64(1)},
	{1, float64(2), float64(1)},
}

func TestSub(t *testing.T) {
	for _, c := range subCases {
		t.Run(fmt.Sprintf("%T(%v) - %T(%v)  = %T(%v)", c[1], c[1], c[0], c[0], c[2], c[2]), func(t *testing.T) {
			v, err := Sub(ValueOf(c[0]), ValueOf(c[1]))
			gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
			gomega.NewWithT(t).Expect(v).To(gomega.Equal(c[2]))
		})
	}
}

func BenchmarkSub(b *testing.B) {
	for _, c := range subCases {
		b.Run(fmt.Sprintf("%T(%v) - %T(%v)", c[1], c[1], c[0], c[0]), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = Sub(ValueOf(c[0]), ValueOf(c[1]))
			}
		})
	}
}
