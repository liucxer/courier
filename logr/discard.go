package logr

import "context"

func Discard() Logger {
	return &discardLogger{}
}

type discardLogger struct {
}

func (d *discardLogger) WithValues(keyAndValues ...interface{}) Logger {
	return d
}

func (d *discardLogger) Start(ctx context.Context, name string, keyAndValues ...interface{}) (context.Context, Logger) {
	return ctx, d
}

func (discardLogger) End() {
}

func (discardLogger) Trace(format string, args ...interface{}) {
}

func (discardLogger) Debug(format string, args ...interface{}) {
}

func (discardLogger) Info(format string, args ...interface{}) {
}

func (discardLogger) Warn(err error) {
}

func (discardLogger) Error(err error) {
}

func (discardLogger) Fatal(err error) {
}

func (discardLogger) Panic(err error) {
}
