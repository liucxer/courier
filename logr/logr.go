package logr

import (
	"context"
)

type Logger interface {
	// Start start span for tracing
	//
	// 	ctx log = logr.Start(ctx, "SpanName")
	// 	defer log.End()
	//
	Start(ctx context.Context, name string, keyAndValues ...interface{}) (context.Context, Logger)
	// End end span
	End()

	// WithValues key value pairs
	WithValues(keyAndValues ...interface{}) Logger

	// info methods
	//
	// Trace trace info
	Trace(msg string, args ...interface{})
	// Debug debug info
	Debug(msg string, args ...interface{})
	// Info info
	Info(msg string, args ...interface{})

	// error methods
	//
	// consider use `github.com/pkg/errors` instead of `errors`
	//
	// Warn
	Warn(err error)
	// Error
	Error(err error)
	// Fatal
	Fatal(err error)
	// Panic
	Panic(err error)
}

type contextKey struct{}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v
	}
	return StdLogger()
}

func Start(ctx context.Context, name string, keyAndValues ...interface{}) (context.Context, Logger) {
	return FromContext(ctx).Start(ctx, name, keyAndValues...)
}
