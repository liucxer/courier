package workeridutil

import (
	"net"
	"os"
)

func ResolveLocalIP() net.IP {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = os.Getenv("HOSTNAME")
	}

	addrList, _ := net.LookupIP(hostname)
	for _, addr := range addrList {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4
		}
	}
	return nil
}

func WorkerIDFromIP(ip net.IP) uint32 {
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return uint32(ip[2])<<8 + uint32(ip[3])
}
