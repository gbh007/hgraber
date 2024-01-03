package logger

import (
	"runtime"
	"strconv"
)

func simpleTrace(skip, count int) string {
	result := ""

	pc := make([]uintptr, count)
	n := runtime.Callers(skip, pc)

	pc = pc[:n]

	frames := runtime.CallersFrames(pc)

	for {
		frame, more := frames.Next()

		if result != "" {
			result += "\n"
		}

		result += "=> " + frame.File + ":" + strconv.Itoa(frame.Line)

		if !more {
			break
		}
	}

	return result
}
