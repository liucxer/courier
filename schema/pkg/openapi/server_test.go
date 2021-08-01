package openapi

import (
	"testing"

	"github.com/liucxer/courier/schema/pkg/testutil"
	"github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(Server{})).To(gomega.Equal(`{"url":""}`))
	})

	t.Run("with variables", func(t *testing.T) {
		server := NewServer("$HOST")
		server.AddVariable("SCHEME", nil)
		server.AddVariable("HOST", NewServerVariable("google.com"))

		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(server)).To(gomega.Equal(`{"url":"$HOST","variables":{"HOST":{"default":"google.com"}}}`))
	})
}
