package system

import (
	"context"
	"time"
)

// Stopwatch секундомер для defer.
//
// Пример использования:
//
//	defer system.Stopwatch(ctx, "some name")()
func Stopwatch(ctx context.Context, text string) func() {
	if !IsDebug(ctx) {
		return func() {}
	}

	tNow := time.Now()

	return func() {
		print(ctx, logLevelDebug, 1, text, time.Since(tNow))
	}
}

// StopwatchWithDepth секундомер для defer с настройкой грубины лога.
//
// Пример использования:
//
//	defer system.StopwatchWithDepth(ctx, "some name", 1)()
func StopwatchWithDepth(ctx context.Context, text string, depth int) func() {
	if !IsDebug(ctx) {
		return func() {}
	}

	tNow := time.Now()

	return func() {
		print(ctx, logLevelDebug, 1+depth, text, time.Since(tNow))
	}
}
