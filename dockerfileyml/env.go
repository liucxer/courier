package dockerfileyml

import "fmt"

func ContainerEnvVar(key string) string {
	return fmt.Sprintf("$${%s}", key)
}

func EnvVar(key string) string {
	return fmt.Sprintf("${%s}", key)
}
