package openapi

import (
	"testing"

	"github.com/liucxer/courier/schema/pkg/testutil"
	"github.com/onsi/gomega"
)

func TestTag(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(Tag{})).To(gomega.Equal(`{"name":""}`))
	})
}
