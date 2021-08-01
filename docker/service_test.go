package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestDockerComposeService(t *testing.T) {
	tt := require.New(t)

	serviceString := []byte(`
image: nginx:alpine
labels:
  base_path: /kubemgr
environment:
  GOENV: ${GOENV}
  S_LOG_LEVEL: ${S_LOG_LEVEL}
  S_MASTERDB_PASSWORD: ${S_MASTERDB_PASSWORD}
  S_MASTERDB_USER: ${S_MASTERDB_USER}
entrypoint:
  - sh
command:
  - "-c"
  - "sleep 50000"
ports:
  - 99:80/tcp
volumes:
  - /tmp:/tmp
volumes_from:
  - content:ro
`)

	service := Service{}

	err := yaml.Unmarshal(serviceString, &service)
	tt.Nil(err)

	tt.Equal(Image{
		Name:    "nginx",
		Version: "alpine",
	}, service.Image)

	tt.Equal([]string{"sh"}, service.EntryPoint.Value())
	tt.Equal([]string{"-c", "sleep 50000"}, service.Command.Value())

	tt.Equal([]Port{
		{
			Port:          99,
			ContainerPort: 80,
			Protocol:      ProtocolTCP,
		},
	}, service.Ports)

	tt.Equal([]Volume{
		{
			LocalPath:  "/tmp",
			MountPath:  "/tmp",
			AccessMode: VolumeAccessModeReadWrite,
		},
	}, service.Volumes)

	tt.Equal([]Volume{
		{
			Name:       "content",
			AccessMode: VolumeAccessModeReadOnly,
		},
	}, service.VolumesFrom)

	bytes, err := yaml.Marshal(service)
	tt.Nil(err)

	t.Log(string(bytes))
}
