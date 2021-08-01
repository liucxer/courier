package metax

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type SomeString struct {
	Ctx
}

func (s *SomeString) WithContext(ctx context.Context) *SomeString {
	return &SomeString{
		Ctx: s.Ctx.WithContext(ctx),
	}
}

func TestCtx(t *testing.T) {
	s := &SomeString{}
	s2 := s.WithContext(ContextWith(context.Background(), "k", "1"))

	require.Equal(t, "", MetaFromContext(s.Context()).Get("k"))
	require.Equal(t, "1", MetaFromContext(s2.Context()).Get("k"))
}
