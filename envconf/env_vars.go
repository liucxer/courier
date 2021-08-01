package envconf

import (
	"bytes"
	"encoding"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func EnvVarsFromEnviron(prefix string, envs []string) *EnvVars {
	e := NewEnvVars(prefix)
	for _, kv := range envs {
		keyValuePair := strings.SplitN(kv, "=", 2)
		if len(keyValuePair) == 2 {
			if strings.HasPrefix(keyValuePair[0], prefix) {
				e.Set(&EnvVar{
					KeyPath: strings.Replace(keyValuePair[0], e.Prefix+"__", "", 1),
					Value:   keyValuePair[1],
				})
			}
		}
	}
	return e
}

func NewEnvVars(prefix string) *EnvVars {
	e := &EnvVars{
		Prefix: prefix,
	}
	return e
}

type EnvVars struct {
	Prefix string
	Values map[string]*EnvVar
}

func (e *EnvVars) SetKeyValue(k string, v string) {
	e.Set(&EnvVar{
		KeyPath: strings.Replace(k, e.Prefix+"__", "", 1),
		Value:   v,
	})
}

func (e *EnvVars) Set(envVar *EnvVar) {
	if e.Values == nil {
		e.Values = map[string]*EnvVar{}
	}
	e.Values[envVar.KeyPath] = envVar
}

func (e *EnvVars) MaskBytes() []byte {
	values := map[string]string{}
	for _, envVar := range e.Values {
		if envVar.Mask != "" {
			values[envVar.Key(e.Prefix)] = envVar.Mask
			continue
		}
		values[envVar.Key(e.Prefix)] = envVar.Value
	}
	return DotEnv(values)
}

func (e *EnvVars) Bytes() []byte {
	values := map[string]string{}
	for _, envVar := range e.Values {
		values[envVar.Key(e.Prefix)] = envVar.Value
	}
	return DotEnv(values)
}

var interfaceTextMarshaller = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
var interfaceTextUnmarshaller = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func (e *EnvVars) Len(key string) int {
	maxIdx := -1

	for _, envVar := range e.Values {
		if strings.HasPrefix(envVar.KeyPath, key) {
			v := strings.TrimLeft(envVar.KeyPath, key+"_")
			parts := strings.Split(v, "_")
			i, err := strconv.ParseInt(parts[0], 10, 64)
			if err == nil {
				if int(i) > maxIdx {
					maxIdx = int(i)
				}
			}
		}
	}

	return maxIdx + 1
}

func (e *EnvVars) Get(key string) *EnvVar {
	if e.Values == nil {
		return nil
	}
	return e.Values[key]
}

func DotEnv(keyValues map[string]string) []byte {
	buf := bytes.NewBuffer(nil)

	keys := make([]string, 0)
	for k := range keyValues {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(keyValues[k])
		buf.WriteRune('\n')
	}

	return buf.Bytes()
}
