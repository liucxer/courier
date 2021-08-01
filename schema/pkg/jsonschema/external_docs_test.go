package jsonschema

import (
	"net/url"
	"testing"

	"github.com/liucxer/courier/schema/pkg/testutil"
	"github.com/onsi/gomega"
)

func TestExternalDoc(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(ExternalDoc{})).To(gomega.Equal(`{}`))
	})

	t.Run("with url", func(t *testing.T) {
		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(ExternalDoc{
			URL: (&url.URL{
				Scheme: "https",
				Host:   "google.com",
			}).String(),
		})).To(gomega.Equal(`{"url":"https://google.com"}`))
	})

	t.Run("with url and description", func(t *testing.T) {
		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(ExternalDoc{
			URL: (&url.URL{
				Scheme: "https",
				Host:   "google.com",
			}).String(),
			Description: "google",
		})).To(gomega.Equal(`{"description":"google","url":"https://google.com"}`))
	})
}
