package envconf

import (
	"strings"
	"testing"
	"time"

	"github.com/liucxer/courier/ptr"
	. "github.com/liucxer/courier/snapshotmacther"
	. "github.com/onsi/gomega"
)

type SubConfig struct {
	Duration     Duration
	Password     Password `env:""`
	Key          string   `env:""`
	Bool         bool
	Map          map[string]string
	Func         func() error
	defaultValue bool
}

func (c *SubConfig) SetDefaults() {
	c.defaultValue = true
}

type Config struct {
	Map       map[string]string
	Slice     []string `env:""`
	PtrString *string  `env:""`
	Host      string   `env:",upstream"`
	SubConfig
	Config SubConfig
}

func TestEnvVars(t *testing.T) {
	c := Config{}

	c.Duration = Duration(10 * time.Second)
	c.Password = Password("123123")
	c.Key = "123456"
	c.PtrString = ptr.String("123456=")
	c.Slice = []string{"1", "2"}
	c.Config.Key = "key"
	c.Config.defaultValue = true
	c.defaultValue = true

	envVars := NewEnvVars("S")

	t.Run("Encoding", func(t *testing.T) {
		data, _ := NewDotEnvEncoder(envVars).Encode(&c)

		NewWithT(t).Expect(string(data)).To(MatchSnapshot(".dotenv"))
	})

	t.Run("Decoding", func(t *testing.T) {
		data, _ := NewDotEnvEncoder(envVars).Encode(&c)

		envVars := EnvVarsFromEnviron("S", strings.Split(string(data), "\n"))

		c2 := Config{}
		err := NewDotEnvDecoder(envVars).Decode(&c2)

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(c2).To(Equal(c))
	})
}
