package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestDockerCompose(t *testing.T) {
	tt := require.New(t)

	d := DockerCompose{}

	yaml.Unmarshal([]byte(`version: "2"

services:
  web: 
    image: "test"
    environment:
      KEY_1: "2"
      KEY_2: "2"
`), &d)

	yaml.Unmarshal([]byte(`version: "2"

services:
  add:
    image: "add"
  web: 
    image: "web"
    environment:
      KEY_1: "1"
      KEY_3: "1"
`), &d)

	tt.Equal(d.Services["web"].Environment, map[string]string{
		"KEY_1": "1",
		"KEY_2": "2",
		"KEY_3": "1",
	})
}
