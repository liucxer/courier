package snapshotmacther

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func TestSnapshot(t *testing.T) {
	scope := t.Name()

	t.Run("match snapshot", func(t *testing.T) {
		NewWithT(t).Expect([]byte(`121

todo
`)).To(MatchSnapshot(scope, "expect.txt"))
	})

	t.Run("not match snapshot", func(t *testing.T) {
		NewWithT(t).Expect("222").NotTo(MatchSnapshot(scope, "test.txt"))
	})

	t.Run("update snapshot", func(t *testing.T) {
		_ = os.Setenv(EnvKeyUpdateSnapshot, "all")

		NewWithT(t).Expect("222").To(MatchSnapshot(scope, "test.txt"))
		NewWithT(t).Expect("111").To(MatchSnapshot(scope, "test.txt"))
	})

	t.Run("update snapshot", func(t *testing.T) {
		_ = os.Setenv(EnvKeyUpdateSnapshot, filepath.Join(scope, "test.txt"))

		NewWithT(t).Expect("222").To(MatchSnapshot(scope, "test.txt"))
		NewWithT(t).Expect("111").To(MatchSnapshot(scope, "test.txt"))
	})
}
