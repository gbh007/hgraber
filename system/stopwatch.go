package system

import (
	"context"
	"time"
)

// Stopwatch секундомер для defer
func Stopwatch(ctx context.Context, text string) func() {
	if !debugMode {
		return func() {}
	}

	tNow := time.Now()

	return func() {
		print(ctx, logLevelDebug, 1, text, time.Since(tNow))
	}
}
