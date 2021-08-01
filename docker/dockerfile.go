package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func EnvVarInDocker(key string) string {
	return fmt.Sprintf("$${%s}", key)
}

func EnvVar(key string) string {
	return fmt.Sprintf("${%s}", key)
}

type Dockerfile struct {
	From       string            `docker:"FROM" yaml:"from,omitempty"`
	Image      string            `yaml:"image"`
	Label      map[string]string `docker:"LABEL,multi" yaml:"label,omitempty"`
	WorkDir    string            `docker:"WORKDIR" yaml:"workdir,omitempty"`
	Env        map[string]string `docker:"ENV,multi" yaml:"env,omitempty"`
	Add        map[string]string `docker:"ADD,join" yaml:"add,omitempty"`
	Run        string            `docker:"RUN,inline" yaml:"run,omitempty"`
	Expose     []string          `docker:"EXPOSE" yaml:"expose,omitempty"`
	Volume     []string          `docker:"VOLUME,array" yaml:"volume,omitempty"`
	Cmd        []string          `docker:"CMD,array" yaml:"cmd,omitempty"`
	EntryPoint []string          `docker:"ENTRYPOINT,array" yaml:"entrypoint,omitempty"`
}

func (d *Dockerfile) String() string {
	return ParseEnvVars(GetDockerfileTemplate(*d), EnvVarsFromEnviron(os.Environ()))
}

func (d Dockerfile) AddContent(from string, to string) *Dockerfile {
	if d.Add == nil {
		d.Add = map[string]string{}
	}
	d.Add[from] = to
	return &d
}

func (d Dockerfile) AddLabel(label string, content string) *Dockerfile {
	if d.Label == nil {
		d.Label = map[string]string{}
	}
	d.Label[label] = content
	return &d
}

func (d Dockerfile) AddEnv(key string, value string) *Dockerfile {
	if d.Env == nil {
		d.Env = map[string]string{}
	}
	d.Env[key] = value
	return &d
}

func (d Dockerfile) WithExpose(exposes ...string) *Dockerfile {
	d.Expose = exposes
	return &d
}

func (d Dockerfile) WithVolume(volumes ...string) *Dockerfile {
	d.Volume = volumes
	return &d
}

func (d Dockerfile) WithWorkDir(dir string) *Dockerfile {
	d.WorkDir = dir
	return &d
}

func (d Dockerfile) WithCmd(cmd ...string) *Dockerfile {
	d.Cmd = cmd
	return &d
}

func GetDockerfileTemplate(d Dockerfile) string {
	dockerfileConfig := make([]string, 0)

	appendDockerConfig := func(dockerKey string, value string) {
		dockerfileConfig = append(
			dockerfileConfig,
			strings.Join([]string{dockerKey, value}, " "),
		)
	}

	tpe := reflect.TypeOf(d)
	rv := reflect.Indirect(reflect.ValueOf(d))

	for i := 0; i < tpe.NumField(); i++ {
		field := tpe.Field(i)
		dockerTag := field.Tag.Get("docker")
		dockerKeys := strings.Split(dockerTag, ",")
		dockerKey := dockerKeys[0]
		dockerFlags := dockerKeys[1:]

		if len(dockerKey) > 0 {
			value := rv.FieldByName(field.Name)

			switch field.Type.Kind() {
			case reflect.String:
				if len(value.String()) > 0 {
					inline := stringIncludes(dockerFlags, "inline")
					if inline {
						appendDockerConfig(dockerKey, value.String())
					} else {
						appendDockerConfig(dockerKey, mayQuote(value.String()))
					}
				}
			case reflect.Slice:
				jsonArray := stringIncludes(dockerFlags, "array")
				slice := make([]string, 0)
				for i := 0; i < value.Len(); i++ {
					slice = append(slice, value.Index(i).String())
				}
				if len(slice) > 0 {
					if jsonArray {
						jsonString, err := json.Marshal(slice)
						if err != nil {
							panic(err)
						}
						appendDockerConfig(
							dockerKey,
							string(jsonString),
						)
					} else {
						appendDockerConfig(
							dockerKey,
							strings.Join(slice, ""),
						)
					}
				}

			case reflect.Map:
				multi := stringIncludes(dockerFlags, "multi")
				join := stringIncludes(dockerFlags, "join")

				if join {
					destMap := map[string][]string{}
					for _, key := range value.MapKeys() {
						dest := value.MapIndex(key).String()
						if destMap[dest] == nil {
							destMap[dest] = []string{}
						}
						destMap[dest] = append(destMap[dest], key.String())
					}
					for dest, src := range destMap {
						appendDockerConfig(
							dockerKey,
							strings.Join(append(src, dest), " "),
						)
					}
				} else if multi {
					keyValues := make([]string, 0)
					for _, key := range value.MapKeys() {
						keyValues = append(
							keyValues,
							key.String()+"="+mayQuote(value.MapIndex(key).String()),
						)
					}
					if len(keyValues) > 0 {
						appendDockerConfig(
							dockerKey,
							strings.Join(keyValues, " "),
						)
					}
				} else {
					for _, key := range value.MapKeys() {
						appendDockerConfig(
							dockerKey,
							strings.Join([]string{key.String(), mayQuote(value.MapIndex(key).String())}, " "),
						)
					}
				}
			}
		}
	}

	return strings.Join(dockerfileConfig, "\n")
}

func mayQuote(s string) string {
	if s == "" || strings.Index(s, " ") > -1 {
		return strconv.Quote(s)
	}
	return s
}
