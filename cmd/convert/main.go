package main

import (
	"app/internal/converter"
	"app/internal/storage/jdb"
	"app/internal/storage/sqlite"
	"app/system"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbFromFilePath := flag.String("from", "db.json", "файл базы")
	dbFromType := flag.String("from-type", "jdb", "Тип БД: jdb, sqlite")

	dbToFilePath := flag.String("to", "main.db", "файл базы")
	dbToType := flag.String("to-type", "sqlite", "Тип БД: jdb, sqlite")

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

	switch *dbFromType {
	case "jdb":
		storageJDB := jdb.Init(ctx, *dbFromFilePath)
		err := storageJDB.Load(ctx, *dbFromFilePath)
		if err != nil {
			system.Error(ctx, err)

			os.Exit(1)
		}

		builder.WithFrom(storageJDB)

	case "sqlite":
		sqliteDB, err := sqlite.Connect(ctx, *dbFromFilePath)
		if err != nil {
			system.Error(ctx, err)

			os.Exit(1)
		}

		err = sqliteDB.MigrateAll(ctx)
		if err != nil {
			system.Error(ctx, err)

			os.Exit(1)
		}

		builder.WithFrom(sqliteDB)
	}

	switch *dbToType {
	case "jdb":
		storageJDB := jdb.Init(ctx, *dbToFilePath)
		err := storageJDB.Load(ctx, *dbToFilePath)
		if err != nil {
			system.Error(ctx, err)

			os.Exit(1)
		}

		builder.WithTo(storageJDB)

	case "sqlite":
		sqliteDB, err := sqlite.Connect(ctx, *dbToFilePath)
		if err != nil {
			system.Error(ctx, err)

			os.Exit(1)
		}

		err = sqliteDB.MigrateAll(ctx)
		if err != nil {
			system.Error(ctx, err)

			os.Exit(1)
		}

		builder.WithTo(sqliteDB)
	}

	builder.Convert(ctx, *offset, true)
}
