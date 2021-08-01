package worker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	count := int64(0)

	process := func(ctx context.Context) error {
		c := atomic.LoadInt64(&count)
		atomic.StoreInt64(&count, c+1)
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	worker := NewWorker(process, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	worker.Start(ctx)

	t.Log(count)
}
