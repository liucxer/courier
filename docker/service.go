package docker

import (
	"fmt"
	"strings"

	"github.com/liucxer/courier/ptr"
)

func NewService(image string) *Service {
	return Service{}.WithImage(image)
}

type Service struct {
	Image          Image             `yaml:"image"`
	EntryPoint     MaybeListString   `yaml:"entrypoint,omitempty"`
	Command        MaybeListString   `yaml:"command,omitempty"`
	Labels         map[string]string `yaml:"labels,omitempty"`
	Environment    map[string]string `yaml:"environment,omitempty"`
	Ports          []Port            `yaml:"ports,omitempty"`
	Links          []Link            `yaml:"links,omitempty"`
	ExternalLinks  []Link            `yaml:"external_links,omitempty"`
	Volumes        []Volume          `yaml:"volumes,omitempty"`
	VolumesFrom    []Volume          `yaml:"volumes_from,omitempty"`
	WorkingDir     string            `yaml:"working_dir,omitempty"`
	DnsSearch      []string          `yaml:"dns_search,omitempty"`
	Dns            []string          `yaml:"dns,omitempty"`
	TTY            *bool             `yaml:"tty,omitempty"`
	MemLimit       *int64            `yaml:"mem_limit,omitempty"`
	MemSwapLimit   *int64            `yaml:"memswap_limit,omitempty"`
	MemReservation *int64            `yaml:"mem_reservation,omitempty"`
}

func (service Service) Merge(nextService *Service) *Service {
	if !nextService.Image.IsZero() {
		service.Image = nextService.Image
	}

	if !nextService.Command.IsZero() {
		service.Command = nextService.Command
	}

	if service.Labels == nil {
		service.Labels = nextService.Labels
	} else {
		for key, value := range nextService.Labels {
			service.Labels[key] = value
		}
	}

	if service.Environment == nil {
		service.Environment = nextService.Environment
	} else {
		for key, value := range nextService.Environment {
			service.Environment[key] = value
		}
	}

	service.Ports = append(service.Ports, nextService.Ports...)
	service.Volumes = append(service.Volumes, nextService.Volumes...)
	service.VolumesFrom = append(service.VolumesFrom, nextService.VolumesFrom...)

	if nextService.WorkingDir != "" {
		service.WorkingDir = nextService.WorkingDir
	}

	if len(nextService.Dns) > 0 {
		service.Dns = nextService.Dns
	}

	if len(nextService.DnsSearch) > 0 {
		service.DnsSearch = nextService.DnsSearch
	}

	if nextService.TTY != nil {
		service.TTY = nextService.TTY
	}

	if nextService.MemLimit != nil {
		service.MemLimit = nextService.MemLimit
	}

	if nextService.MemLimit != nil {
		service.MemLimit = nextService.MemLimit
	}

	return &service
}

func (service Service) addPort(port int16, containerPort int16, protocol Protocol) *Service {
	service.Ports = append(service.Ports, Port{
		Port:          port,
		ContainerPort: containerPort,
		Protocol:      protocol,
	})
	return &service
}

func (service Service) addVolume(nameOrLocalPath string, mountPath string, accessMode VolumeAccessMode) *Service {
	v, err := ParseVolumeString(strings.Join([]string{
		nameOrLocalPath,
		mountPath,
		string(accessMode),
	}, ":"))
	if err != nil {
		panic(err)
	}

	service.Volumes = append(service.Volumes, *v)
	return &service
}

func (service Service) WithImage(image string) *Service {
	i, err := ParseImageString(image)
	if err != nil {
		panic(fmt.Sprintf("invalid image %s", image))
	}
	service.Image = *i
	return &service
}

func (service Service) EnableTTY() *Service {
	service.TTY = ptr.Bool(true)
	return &service
}

func (service Service) AddLink(s string, host string) *Service {
	service.Links = append(service.Links, Link{
		Service: s,
		Host:    host,
	})
	return &service
}

func (service Service) AddDns(dns string, dnsSearch string) *Service {
	service.DnsSearch = []string{dnsSearch}
	service.Dns = []string{dns}
	return &service
}

func (service Service) AddExternalLink(s string, host string) *Service {
	service.ExternalLinks = append(service.ExternalLinks, Link{
		Service: s,
		Host:    host,
	})
	return &service
}

func (service Service) AddTCPPort(port int16, containerPort int16) *Service {
	service.addPort(port, containerPort, ProtocolTCP)
	return &service
}

func (service Service) AddUDPPort(port int16, containerPort int16) *Service {
	service.addPort(port, containerPort, ProtocolUDP)
	return &service
}

func (service Service) SetCommand(commands ...string) *Service {
	service.Command = MaybeListStringFromStringList(commands...)
	return &service
}

func (service Service) AddVolumeFrom(name string, accessMode VolumeAccessMode) *Service {
	service.VolumesFrom = append(service.VolumesFrom, Volume{
		Name:       name,
		AccessMode: accessMode,
	})
	return &service
}

func (service Service) AddRWVolume(nameOrLocalPath string, mountPath string) *Service {
	service.addVolume(nameOrLocalPath, mountPath, VolumeAccessModeReadWrite)
	return &service
}

func (service Service) AddROVolume(nameOrLocalPath string, mountPath string) *Service {
	service.addVolume(nameOrLocalPath, mountPath, VolumeAccessModeReadOnly)
	return &service
}

func (service Service) SetLabel(key string, value string) *Service {
	if service.Labels == nil {
		service.Labels = map[string]string{}
	}
	service.Labels[key] = value
	return &service
}

func (service Service) SetEnvironment(key string, value string) *Service {
	if service.Environment == nil {
		service.Environment = map[string]string{}
	}
	service.Environment[key] = value
	return &service
}
