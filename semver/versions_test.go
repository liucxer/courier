package semver

import (
	"sort"
	"testing"

	. "github.com/onsi/gomega"
)

func TestVersions(t *testing.T) {
	raw := []string{
		"1.2.3",
		"1.0",
		"1.3",
		"2",
		"0.4.2",
	}

	vs := make([]Version, len(raw))
	for i, r := range raw {
		v, err := ParseVersion(r)
		NewWithT(t).Expect(err).To(BeNil())
		vs[i] = *v
	}

	sort.Sort(Versions(vs))

	e := []string{
		"0.4.2",
		"1.0.0",
		"1.2.3",
		"1.3.0",
		"2.0.0",
	}

	a := make([]string, len(vs))
	for i, v := range vs {
		a[i] = v.String()
	}

	NewWithT(t).Expect(a).To(Equal(e))
}
