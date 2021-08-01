package docker

import (
	"fmt"
	"path"
	"strings"
)

type VolumeAccessMode string

const (
	VolumeAccessModeReadWrite = "rw"
	VolumeAccessModeReadOnly  = "ro"
)

type Volume struct {
	Name       string
	LocalPath  string
	MountPath  string
	AccessMode VolumeAccessMode
}

func (v Volume) ReadOnly() bool {
	return v.AccessMode == VolumeAccessModeReadOnly
}

func (v Volume) IsMountFromStorage() bool {
	return v.Name != ""
}

func (v Volume) IsMountFromLocal() bool {
	return v.LocalPath != ""
}

func (v Volume) IsProvider() bool {
	return v.LocalPath == "" && v.Name == ""
}

func isRel(p string) bool {
	return len(p) > 0 && p[0] == '.'
}

func ParseVolumeString(s string) (*Volume, error) {
	v := Volume{}
	values := strings.Split(s, ":")

	if values[0] == "container" {
		values = values[1:]
	}

	accessMode := strings.ToLower(values[len(values)-1])
	if accessMode == VolumeAccessModeReadOnly || accessMode == VolumeAccessModeReadWrite {
		v.AccessMode = VolumeAccessMode(accessMode)
		values = values[0 : len(values)-1]
	}

	mouthPath := values[len(values)-1]

	if path.IsAbs(mouthPath) {
		v.MountPath = mouthPath
		if len(values) == 2 {
			if path.IsAbs(values[0]) || isRel(values[0]) {
				v.LocalPath = values[0]
			} else {
				v.Name = values[0]
			}
		}
	} else {
		v.Name = values[0]
	}

	if v.AccessMode == "" {
		if !v.IsProvider() {
			v.AccessMode = VolumeAccessModeReadWrite
		}
	}

	return &v, nil
}

func (v Volume) String() string {
	if v.AccessMode == "" {
		v.AccessMode = VolumeAccessModeReadWrite
	}

	if v.IsMountFromLocal() {
		return fmt.Sprintf("%s:%s:%s", v.LocalPath, v.MountPath, v.AccessMode)
	}
	if v.IsMountFromStorage() {
		if v.MountPath != "" {
			return fmt.Sprintf("%s:%s:%s", v.Name, v.MountPath, v.AccessMode)
		}
		return fmt.Sprintf("%s:%s", v.Name, v.AccessMode)
	}
	return v.MountPath
}

func (v Volume) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

func (v *Volume) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err == nil {
		volume, err := ParseVolumeString(s)
		if err != nil {
			return err
		}
		*v = *volume
	}
	return nil
}
