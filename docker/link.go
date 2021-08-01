package docker

import (
	"fmt"
	"strings"
)

type Link struct {
	Service string
	Host    string
}

func ParseLinkString(s string) (*Link, error) {
	l := Link{}
	serviceHost := strings.Split(s, ":")

	if len(serviceHost) != 2 {
		return nil, fmt.Errorf("service:host format error")
	}

	l.Service = serviceHost[0]
	l.Host = serviceHost[1]

	return &l, nil
}

func (link Link) String() string {
	return fmt.Sprintf("%s:%s", link.Service, link.Host)
}

func (link Link) MarshalYAML() (interface{}, error) {
	return link.String(), nil
}

func (link *Link) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	err := unmarshal(&v)
	if err == nil {
		l, err := ParseLinkString(v)
		if err != nil {
			return err
		}
		*link = *l
	}
	return nil
}
