package docker

import (
	"fmt"
	"strconv"
	"strings"
)

type Protocol string

const (
	ProtocolTCP Protocol = "TCP"
	ProtocolUDP Protocol = "UDP"
)

type Port struct {
	Port          int16
	ContainerPort int16
	Protocol      Protocol
}

func ParsePortString(s string) (*Port, error) {
	p := Port{}
	portsAndProtocol := strings.Split(s, "/")
	ports := strings.Split(portsAndProtocol[0], ":")

	if len(ports) != 2 {
		return nil, fmt.Errorf("port:containerPort format error")
	}

	port, err := strconv.ParseInt(ports[0], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid port %s", ports[0])
	}
	p.Port = int16(port)

	containerPort, err := strconv.ParseInt(ports[1], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid container port %s", ports[1])
	}

	p.ContainerPort = int16(containerPort)

	if len(portsAndProtocol) == 2 && strings.ToUpper(portsAndProtocol[1]) == string(ProtocolUDP) {
		p.Protocol = ProtocolUDP
	} else {
		p.Protocol = ProtocolTCP
	}

	return &p, nil
}

func (port Port) String() string {
	return fmt.Sprintf("%d:%d/%s", port.Port, port.ContainerPort, port.Protocol)
}

func (port Port) MarshalYAML() (interface{}, error) {
	return port.String(), nil
}

func (port *Port) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	err := unmarshal(&v)
	if err == nil {
		p, err := ParsePortString(v)
		if err != nil {
			return err
		}
		*port = *p
	}
	return nil
}
