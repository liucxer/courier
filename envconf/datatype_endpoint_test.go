package envconf

import (
	"net/url"
	"testing"

	"github.com/davecgh/go-spew/spew"
	. "github.com/onsi/gomega"
)

func TestEndpoint(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		var id = Endpoint{
			Scheme:   "stmps",
			Hostname: "smtp.exmail.qq.com",
			Port:     465,
		}
		NewWithT(t).Expect(id.String()).To(Equal("stmps://smtp.exmail.qq.com:465"))
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		var id Endpoint

		err := id.UnmarshalText([]byte("stmps://smtp.exmail.qq.com:465"))
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(id).To(Equal(Endpoint{
			Scheme:   "stmps",
			Hostname: "smtp.exmail.qq.com",
			Port:     465,
		}))
	})

	t.Run("UnmarshalText full", func(t *testing.T) {
		var id Endpoint

		err := id.UnmarshalText([]byte("postgres://postgres:postgres@127.0.0.1:5432/postgres/xxx?sslmode=disable"))
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(id).To(Equal(Endpoint{
			Scheme:   "postgres",
			Hostname: "127.0.0.1",
			Password: "postgres",
			Username: "postgres",
			Port:     5432,
			Base:     "postgres",
			Extra:    url.Values{"sslmode": {"disable"}},
		}))

		NewWithT(t).Expect(id.String()).To(Equal("postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"))
		NewWithT(t).Expect(id.SecurityString()).To(Equal("postgres://postgres:--------@127.0.0.1:5432/postgres?sslmode=disable"))
	})

	t.Run("UnmarshalExtra", func(t *testing.T) {
		opt := struct {
			ConnectTimeout Duration `name:"connectTimeout" default:"10s"`
			ReadTimeout    Duration `name:"readTimeout" default:"10s"`
			WriteTimeout   Duration `name:"writeTimeout" default:"10s"`
			IdleTimeout    Duration `name:"idleTimeout" default:"240s"`
			MaxActive      int      `name:"maxActive" default:"5"`
			MaxIdle        int      `name:"maxIdle" default:"3"`
			DB             int      `name:"db" default:"10"`
		}{}

		err := UnmarshalExtra(url.Values{}, &opt)
		NewWithT(t).Expect(err).To(BeNil())

		spew.Dump(opt)
	})
}
