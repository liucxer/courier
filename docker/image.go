package docker

import (
	"fmt"
	"strings"
)

type Image struct {
	Name    string
	Version string
}

func (image Image) IsZero() bool {
	return image.Name == ""
}

func ParseImageString(s string) (*Image, error) {
	i := Image{}
	nameAndVersion := strings.Split(s, ":")

	if len(nameAndVersion) == 2 {
		i.Version = nameAndVersion[1]
	}

	i.Name = nameAndVersion[0]
	return &i, nil
}

func (i Image) String() string {
	if i.Version != "" {
		return fmt.Sprintf(`%s:%s`, i.Name, i.Version)
	}
	return i.Name
}

func (i Image) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

func (i *Image) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	err := unmarshal(&v)
	if err == nil {
		image, err := ParseImageString(v)
		if err != nil {
			return err
		}
		*i = *image
	}
	return nil
}
