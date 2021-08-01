package openapi

import (
	"testing"

	"github.com/liucxer/courier/schema/pkg/testutil"
	"github.com/onsi/gomega"
)

func TestInfo(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(Info{})).To(gomega.Equal(`{"title":"","version":""}`))
	})

	t.Run("with contact", func(t *testing.T) {
		info := Info{
			InfoObject: InfoObject{
				Contact: &Contact{
					ContactObject: ContactObject{
						Name:  "name",
						URL:   "url",
						Email: "email",
					},
				},
			},
		}

		expect := `{"title":"","contact":{"name":"name","url":"url","email":"email"},"version":""}`

		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(info)).To(gomega.Equal(expect))
	})

	t.Run("with licence", func(t *testing.T) {
		actual := Info{
			InfoObject: InfoObject{
				License: &License{
					LicenseObject: LicenseObject{
						Name: "MIT",
					},
				},
			},
		}

		expect := `{"title":"","license":{"name":"MIT"},"version":""}`

		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(actual)).To(gomega.Equal(expect))
	})

	t.Run("with specification_extensions", func(t *testing.T) {
		actual := Info{
			InfoObject: InfoObject{
				Contact: &Contact{
					SpecExtensions: SpecExtensions{
						Extensions: map[string]interface{}{
							"x-x": "x",
						},
					},
				},
				License: &License{
					SpecExtensions: SpecExtensions{
						Extensions: map[string]interface{}{
							"x-x": "x",
						},
					},
				},
			},
			SpecExtensions: SpecExtensions{
				Extensions: map[string]interface{}{
					"x-x": "x",
				},
			},
		}

		expect := `{"title":"","contact":{"x-x":"x"},"license":{"name":"","x-x":"x"},"version":"","x-x":"x"}`

		gomega.NewWithT(t).Expect(testutil.MustJSONRaw(actual)).To(gomega.Equal(expect))
	})
}
