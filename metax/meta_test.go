package metax

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMeta(t *testing.T) {
	t.Run("parse id", func(t *testing.T) {
		meta := ParseMeta("xxxxxx")
		require.Equal(t, "_id=xxxxxx", meta.String())
	})

	t.Run("parse meta", func(t *testing.T) {
		meta := ParseMeta("operator=1&operator=2&_id=xxx")
		require.Equal(t, "1", meta.Get("operator"))
		require.Equal(t, "_id=xxx&operator=1&operator=2", meta.String())
	})
}

func TestMeta(t *testing.T) {

	t.Run("ContextConcat", func(t *testing.T) {
		ctx := ContextWith(context.Background(), "key", "1")
		ctx = ContextWithMeta(ctx, (Meta{}).With("key", "2", "3"))

		require.Equal(t, []string{"1", "2", "3"}, MetaFromContext(ctx)["key"])
	})

	t.Run("ContextOverwrite", func(t *testing.T) {
		ctx := ContextWith(context.Background(), "_key", "1")
		ctx = ContextWithMeta(ctx, (Meta{}).With("_key", "2", "3"))

		require.Equal(t, []string{"2", "3"}, MetaFromContext(ctx)["_key"])
	})

	t.Run("EmptyKeyIgnore", func(t *testing.T) {
		ctx := ContextWith(context.Background(), "", "1")
		require.Len(t, MetaFromContext(ctx), 0)
	})
}
