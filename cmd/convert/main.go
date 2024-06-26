package main

import (
	"app/internal/dataprovider/logger"
	"app/internal/dataprovider/storage/jdb"
	"app/internal/dataprovider/storage/postgresql"
	"app/internal/usecase/converter"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type connector func(ctx context.Context, logger *slog.Logger, builder *converter.Builder, path string, from bool) error

func main() {
	dbFromFilePath := flag.String("from", "db.json", "файл базы")
	dbFromType := flag.String("from-type", "jdb", "Тип БД: jdb, pg")

	dbToFilePath := flag.String("to", "main.db", "файл базы")
	dbToType := flag.String("to-type", "sqlite", "Тип БД: jdb, pg")

	offset := flag.Int("offset", 0, "Пропустить количество")

	// Отладка
	debug := flag.Bool("debug", false, "Режим отладки")
	debugTrace := flag.Bool("debug-trace", false, "Режим стектрейсов")

	flag.Parse()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	logger := logger.New(*debug, *debugTrace)

	builder := converter.New(logger)

	var fromConnector, toConnector connector

	switch *dbFromType {
	case "jdb":
		fromConnector = jdbConnect
	case "pg":
		fromConnector = pgConnect
	}

	switch *dbToType {
	case "jdb":
		toConnector = jdbConnect
	case "pg":
		toConnector = pgConnect
	}

	if fromConnector == nil || toConnector == nil {
		logger.ErrorContext(ctx, "nil connector")
		os.Exit(1)
	}

	err := fromConnector(ctx, logger, builder, *dbFromFilePath, true)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		os.Exit(1)
	}

	err = toConnector(ctx, logger, builder, *dbToFilePath, false)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		os.Exit(1)
	}

	builder.Convert(ctx, *offset, true)
}

func jdbConnect(ctx context.Context, logger *slog.Logger, builder *converter.Builder, path string, from bool) error {
	storageJDB := jdb.Init(ctx, logger, nil)
	err := storageJDB.Load(ctx, path)
	if err != nil {
		return err
	}

	if from {
		builder.WithFrom(storageJDB)
	} else {
		builder.WithTo(storageJDB)
	}

	return nil
}

func pgConnect(ctx context.Context, logger *slog.Logger, builder *converter.Builder, path string, from bool) error {
	postgresql, err := postgresql.Connect(ctx, path, logger)
	if err != nil {
		return err
	}

	err = postgresql.MigrateAll(ctx)
	if err != nil {
		return err
	}

	if from {
		builder.WithFrom(postgresql)
	} else {
		builder.WithTo(postgresql)
	}

	return nil
}
