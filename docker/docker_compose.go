package docker

func NewDockerCompose() *DockerCompose {
	return &DockerCompose{
		Version: "2",
	}
}

type DockerCompose struct {
	Version  string   `yaml:"version"`
	Services Services `yaml:"services,omitempty"`
}

type Services map[string]*Service

func (v *Services) UnmarshalYAML(unmarshal func(interface{}) error) error {
	set := map[string]*Service{}
	err := unmarshal(&set)
	if err != nil {
		return err
	}

	ss := *v
	if ss == nil {
		*v = set
	} else {
		for k, s := range set {
			if service, ok := ss[k]; ok {
				ss[k] = service.Merge(s)
			} else {
				ss[k] = service
			}
		}
		*v = ss
	}

	return nil
}

func (dc DockerCompose) AddService(name string, s *Service) *DockerCompose {
	if dc.Services == nil {
		dc.Services = map[string]*Service{}
	}
	dc.Services[name] = s
	return &dc
}
