package semver

import (
	"testing"
)

func BenchmarkNewConstraint(b *testing.B) {
	benchNewConstraint := func(c string, b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, _ = NewConstraint(c)
		}
	}

	b.Run("Unary", func(b *testing.B) {
		benchNewConstraint("=2.0", b)
	})

	b.Run("Tilde", func(b *testing.B) {
		benchNewConstraint("~2.0.0", b)
	})

	b.Run("Caret", func(b *testing.B) {
		benchNewConstraint("^2.0.0", b)
	})

	b.Run("Wildcard", func(b *testing.B) {
		benchNewConstraint("1.x", b)
	})

	b.Run("Range", func(b *testing.B) {
		benchNewConstraint(">=2.1.x, <3.1.0", b)
	})

	b.Run("Union", func(b *testing.B) {
		benchNewConstraint("~2.0.0 || =3.1.0", b)
	})
}

func BenchmarkCheckVersion(b *testing.B) {
	benchCheckVersion := func(c, v string, b *testing.B) {
		version, _ := ParseVersion(v)
		constraint, _ := NewConstraint(c)

		for i := 0; i < b.N; i++ {
			constraint.Check(version)
		}
	}

	b.Run("Unary", func(b *testing.B) {
		benchCheckVersion("=2.0", "2.0.0", b)
	})

	b.Run("Tilde", func(b *testing.B) {
		benchCheckVersion("~2.0.0", "2.0.5", b)
	})

	b.Run("Caret", func(b *testing.B) {
		benchCheckVersion("^2.0.0", "2.1.0", b)
	})

	b.Run("Wildcard", func(b *testing.B) {
		benchCheckVersion("1.x", "1.4.0", b)
	})

	b.Run("Range", func(b *testing.B) {
		benchCheckVersion(">=2.1.x, <3.1.0", "2.4.5", b)
	})

	b.Run("Union", func(b *testing.B) {
		benchCheckVersion("~2.0.0 || =3.1.0", "3.1.0", b)
	})
}

func BenchmarkNewVersion(b *testing.B) {
	benchNewVersion := func(v string, b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = ParseVersion(v)
		}
	}

	b.Run("simple", func(b *testing.B) {
		benchNewVersion("1.0.0", b)
	})

	b.Run("pre", func(b *testing.B) {
		benchNewVersion("1.0.0-alpha", b)
	})

	b.Run("meta", func(b *testing.B) {
		benchNewVersion("1.0.0+buildMetadata", b)
	})

	b.Run("all", func(b *testing.B) {
		benchNewVersion("1.0.0-alpha.1+meta.data", b)
	})
}
