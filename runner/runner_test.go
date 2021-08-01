package runner

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunner(t *testing.T) {
	tt := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wp := NewRunner(ctx, "normal", 3)

	doneCount := int32(0)

	for i := 0; i < int(9); i++ {
		count := i
		wp.Add(func(ctx context.Context) error {
			time.Sleep(500 * time.Millisecond)
			atomic.AddInt32(&doneCount, 1)
			fmt.Printf("job done %d\n", count)
			return nil
		}, fmt.Sprintf("%d", count))
	}

	if err := wp.Commit(); err != nil {
		t.Logf("err %s", err.Error())
	}

	tt.Equal(doneCount, int32(9))
}

func TestRunnerWhenNoTask(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	wp := NewRunner(ctx, "no task", 5)

	err := wp.Commit()
	t.Logf("err %v", err)
}

func TestRunnerWhenTimeout(t *testing.T) {
	tt := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	wp := NewRunner(ctx, "timeout", 3)

	for i := 0; i < int(9); i++ {
		count := i
		wp.Add(func(ctx context.Context) error {
			fmt.Printf("job timeout %d\n", count)
			time.Sleep(50 * time.Millisecond)
			return nil
		}, fmt.Sprintf("%d", i))
	}

	err := wp.Commit()
	tt.Error(err)
	t.Log(err.Error())
	tt.True(IsErrTimeout(err))
}

func TestRunnerWhenPanic(t *testing.T) {
	tt := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	wp := NewRunner(ctx, "panic", 3)

	for i := 0; i < int(9); i++ {
		wp.Add(func(ctx context.Context) error {
			panic("panic")
		}, fmt.Sprintf("%d", i))
	}

	err := wp.Commit()
	tt.Error(err)
	t.Log(err.Error())
}

func TestRunnerWhenTaskSmallThankWorkers(t *testing.T) {
	tt := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	wp := NewRunner(ctx, "panic", 3)

	for i := 0; i < int(1); i++ {
		wp.Add(func(ctx context.Context) error {
			return nil
		}, fmt.Sprintf("%d", i))
	}

	err := wp.Commit()
	tt.NoError(err)
}
