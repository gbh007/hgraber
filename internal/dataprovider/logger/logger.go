package logger

import (
	"app/internal/dataprovider/slogHandler"
	"app/pkg/ctxtool"
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

func New(debug bool, trace bool) *slog.Logger {
	opts := []slogHandler.HandlerOption{}

	if debug {
		opts = append(opts, slogHandler.WithDebug())
	}

	if trace {
		opts = append(opts, slogHandler.WithPrinter(printWithTrace))
	} else {
		opts = append(opts, slogHandler.WithPrinter(print))
	}

	return slog.New(slogHandler.New(opts...))
}

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
