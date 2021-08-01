package logr

import (
	"context"
)

func ExampleLogger() {
	ctx := WithLogger(context.Background(), StdLogger())

	logr := FromContext(ctx).WithValues("k", "k")

	logr.Debug("test %d", 1)
	logr.Trace("test %d", 1)
	logr.Info("test %d", 1)
	// Output:
}

func ExampleLogger_Start() {
	ctx := WithLogger(context.Background(), StdLogger())

	_, log := Start(ctx, "span", "k", "k")
	defer log.End()

	log.Debug("test %d", 1)
	log.Trace("test %d", 1)
	log.Info("test %d", 1)
	// Output:
}
