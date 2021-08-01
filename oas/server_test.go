package oas

import (
	"testing"
)

func TestServer(t *testing.T) {
	g := NewCaseGroup("Server")

	g.It("empty", `{"url":""}`, Server{})

	g.It("with variables", `{"url":"$HOST","variables":{"HOST":{"default":"google.com"}}}`, func() *Server {
		server := NewServer("$HOST")
		server.AddVariable("SCHEME", nil)
		server.AddVariable("HOST", NewServerVariable("google.com"))
		return server
	}())

	g.Run(t)
}
