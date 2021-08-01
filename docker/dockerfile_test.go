package docker

import (
	"testing"
)

func TestDockerfile(t *testing.T) {
	SOME_ENV := "SOME_ENV"

	d := &Dockerfile{
		From: "busybox:latest",
		EntryPoint: []string{
			"sh",
		},
		Env: map[string]string{
			SOME_ENV: "hello",
		},
		Cmd: []string{
			"-c",
			"echo",
			EnvVarInDocker(SOME_ENV),
		},
	}

	t.Log(d.String())
}
