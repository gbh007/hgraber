package system

var debugMode = false

func EnableDebug(ctx Context) {
	debugMode = true
	print(ctx, logLevelWarning, 0, "Включен режим отладки")
}

func DisableDebug(ctx Context) {
	debugMode = false
	print(ctx, logLevelWarning, 0, "Отключен режим отладки")
}

func DebugStatus() bool {
	return debugMode
}

func Debug(ctx Context, args ...interface{}) {
	if !debugMode {
		return
	}

	print(ctx, logLevelDebug, 0, args...)
}
