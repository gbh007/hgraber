package logger

import (
	"app/pkg/ctxtool"
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

func attrsToString(attrs []slog.Attr) string {
	if len(attrs) == 0 {
		return ""
	}

	return " " + fmt.Sprint(attrs)
}

func print(ctx context.Context, t time.Time, msg string, lv slog.Level, attrs []slog.Attr) error {
	attrsToPrint := attrsToString(attrs)

	_, err := fmt.Fprintf(
		os.Stderr,
		"[%s] [%s] %s - %s%s\n",
		lv.String(),
		ctxtool.GetRequestID(ctx),
		t.Format(timeFormat),
		msg,
		attrsToPrint,
	)

	return err
}

func printWithTrace(ctx context.Context, t time.Time, msg string, lv slog.Level, attrs []slog.Attr) error {
	attrsToPrint := attrsToString(attrs)

	_, err := fmt.Fprintf(
		os.Stderr,
		"[%s] [%s] %s - %s%s\n%s\n",
		lv.String(),
		ctxtool.GetRequestID(ctx),
		t.Format(timeFormat),
		msg,
		attrsToPrint,
		simpleTrace(6, 4),
	)

	return err
}
