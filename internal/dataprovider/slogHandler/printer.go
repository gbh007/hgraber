package slogHandler

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type Printer func(ctx context.Context, t time.Time, msg string, lv slog.Level, attrs []slog.Attr) error

const timeFormat = "2006-01-02 15:04:05"

func stdoutPrinter(
	ctx context.Context,
	t time.Time,
	msg string,
	lv slog.Level,
	attrs []slog.Attr,
) error {
	attrsToPrint := fmt.Sprint(attrs)

	fmt.Fprintf(
		os.Stderr,
		"[%s] %s - %s %s\n",
		lv.String(),
		t.Format(timeFormat),
		msg,
		attrsToPrint,
	)

	return nil
}
