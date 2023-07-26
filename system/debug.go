package system

import "context"

func WithDebug(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, debugKey, true)

	print(ctx, logLevelWarning, 0, "Включен режим отладки")

	return ctx
}

func IsDebug(ctx context.Context) bool {
	v := ctx.Value(debugKey)
	if v == nil {
		return false
	}

	// Значение интересует только если истина;
	// его отсутствие, неправильный формат, лож эквивалентны
	debug, _ := v.(bool)

	return debug
}

func Debug(ctx context.Context, args ...interface{}) {
	if !IsDebug(ctx) {
		return
	}

	print(ctx, logLevelDebug, 0, args...)
}
