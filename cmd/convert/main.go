package main

import (
	"app/internal/converter"
	"app/internal/storage/jdb"
	"app/internal/storage/postgresql"
	"app/system"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

type connector func(ctx context.Context, builder *converter.Builder, path string, from bool) error

func main() {
	dbFromFilePath := flag.String("from", "db.json", "файл базы")
	dbFromType := flag.String("from-type", "jdb", "Тип БД: jdb, pg")

	dbToFilePath := flag.String("to", "main.db", "файл базы")
	dbToType := flag.String("to-type", "sqlite", "Тип БД: jdb, pg")

	offset := flag.Int("offset", 0, "Пропустить количество")

	flag.Parse()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	ctx = system.WithDebug(ctx)

	builder := new(converter.Builder)

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
		system.ErrorText(ctx, "nil connector")
		os.Exit(1)
	}

	err := fromConnector(ctx, builder, *dbFromFilePath, true)
	if err != nil {
		system.Error(ctx, err)

		os.Exit(1)
	}

	err = toConnector(ctx, builder, *dbToFilePath, false)
	if err != nil {
		system.Error(ctx, err)

		os.Exit(1)
	}

	builder.Convert(ctx, *offset, true)
}

func jdbConnect(ctx context.Context, builder *converter.Builder, path string, from bool) error {
	storageJDB := jdb.Init(ctx, nil)
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

func pgConnect(ctx context.Context, builder *converter.Builder, path string, from bool) error {
	postgresql, err := postgresql.Connect(ctx, path)
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
