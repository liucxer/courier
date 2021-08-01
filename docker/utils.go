package docker

import (
	"regexp"
	"strings"
)

var reEnvVar = regexp.MustCompile("(\\$?\\$)\\{?([A-Za-z0-9_]+)\\}?")

func EnvVarsFromEnviron(environ []string) map[string]string {
	m := map[string]string{}
	for _, kv := range environ {
		kvParts := strings.Split(kv, "=")
		if len(kvParts) == 2 {
			m[kvParts[0]] = kvParts[1]
		}
	}
	return m
}

func ParseEnvVars(s string, envVars map[string]string) string {
	result := reEnvVar.ReplaceAllStringFunc(s, func(str string) string {
		matched := reEnvVar.FindAllStringSubmatch(str, -1)[0]

		// skip $${ }
		if matched[1] == "$$" {
			return "${" + matched[2] + "}"
		}

		if value, ok := envVars[matched[2]]; ok {
			return value
		}

		return "${" + matched[2] + "}"
	})

	return result
}

func stringSome(list []string, checker func(item string, i int) bool) bool {
	for i, item := range list {
		if checker(item, i) {
			return true
		}
	}
	return false
}

func stringIncludes(list []string, target string) bool {
	return stringSome(list, func(item string, i int) bool {
		return item == target
	})
}
