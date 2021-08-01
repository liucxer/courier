package workeridutil

import (
	"net"
	"testing"
)

func TestWorkerIDFromIP(t *testing.T) {
	ip := net.ParseIP("255.255.255.255")
	t.Log(WorkerIDFromIP(ip))
	t.Log(WorkerIDFromIP(ResolveLocalIP()))
}
