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
