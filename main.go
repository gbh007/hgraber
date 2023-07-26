package main

import (
	"app/service/fileStorage"
	"app/service/parser"
	"app/service/titleHandler"
	"app/service/webServer"
	"app/storage/jdb"
	"app/storage/schema"
	"app/super"
	"app/system"
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// базовые опции
	webPort := flag.Int("p", 8080, "порт веб сервера")
	onlyView := flag.Bool("v", false, "режим только просмотра")
	token := flag.String("access-token", "", "токен для доступа к контенту")

	// потоки логирования
	disableStdErr := flag.Bool("no-stderr", false, "отключить стандартный поток ошибок")
	disableFileErr := flag.Bool("no-stdfile", false, "отключить поток ошибок в файл")
	enableAppendFileErr := flag.Bool("stdfile-append", false, "режим дозаписи файла потока ошибок")

	// размещение данных
	fileStoragePath := flag.String("fs", "loads", "директория для данных")
	fileExport := flag.String("fe", "exported", "директория для экспорта файлов")
	dbFileName := flag.String("db", "db.json", "файл базы")
	staticDirName := flag.String("static", "", "папка со статическими файлами")

	// отладка
	debugMode := flag.Bool("debug", false, "активировать режим отладки (дебага)")
	// debugCopyMode := flag.Bool("debug-copy", false, "включает при активном дебаге, информацию о копировании данных в памяти")
	debugFullpathMode := flag.Bool("debug-fullpath", false, "включает длинные пути файлов в логах")

	flag.Parse()

	system.Init(system.LogConfig{
		EnableFile:   !*disableFileErr,
		AppendMode:   *enableAppendFileErr,
		EnableStdErr: !*disableStdErr,
	})

	notifyCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	mainContext := system.NewSystemContext(notifyCtx, "Main")

	if *debugMode {
		system.EnableDebug(mainContext)
	}

	if *debugFullpathMode {
		system.EnableFullpath(mainContext)
	}

	system.Debug(mainContext, "Версия", system.Version)
	system.Debug(mainContext, "Коммит", system.Commit)
	system.Debug(mainContext, "Собрано", system.BuildAt)

	system.Info(mainContext, "Инициализация базы")

	storage := jdb.Init(mainContext, *dbFileName)

	err := storage.Load(mainContext, *dbFileName)
	if err != nil {
		os.Exit(1)
	}

	system.Info(mainContext, "База загружена")

	titleService := titleHandler.Init(storage)
	pageService := fileStorage.Init(storage)

	so := super.NewObject(storage, titleService)
	so.RegisterRunner(mainContext, storage)

	err = system.SetFileStoragePath(mainContext, *fileStoragePath)
	if err != nil {
		os.Exit(2)
	}

	err = system.SetFileExportPath(mainContext, *fileExport)
	if err != nil {
		os.Exit(3)
	}

	if !*onlyView {
		go parseTaskFile(mainContext, titleService)

		so.RegisterRunner(mainContext, titleService)
		so.RegisterRunner(mainContext, pageService)
	}

	webServer := &webServer.WebServer{
		Storage:   storage,
		Title:     titleService,
		Page:      pageService,
		Addr:      fmt.Sprintf(":%d", *webPort),
		StaticDir: *staticDirName,
		Token:     *token,
	}
	so.RegisterRunner(mainContext, webServer)

	system.Info(mainContext, "Завершение работы, ожидание завершения процессов")

	err = so.Run(mainContext)
	if err != nil {
		os.Exit(4)
	}

	// FIXME: удалить
	<-system.WaitingChan(mainContext)

	system.Info(mainContext, "Процессы завершены")

	if storage.Save(mainContext, *dbFileName, false) == nil {
		system.Info(mainContext, "База сохранена")
	} else {
		system.Warning(mainContext, "База не сохранена")
	}

	system.Info(mainContext, "Выход")
}

func parseTaskFile(ctx context.Context, titleService super.TitleHandler) {
	f, err := os.Open("task.txt")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			system.Error(ctx, err)
		}
		return
	}
	defer system.IfErrFunc(ctx, f.Close)

	var (
		totalCount     = 0
		loadedCount    = 0
		duplicateCount = 0
		errorCount     = 0
	)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		u := sc.Text()
		if u == "" {
			continue
		}

		totalCount++

		err = titleService.FirstHandle(ctx, u)

		switch {
		case errors.Is(err, schema.TitleDuplicateError):
			duplicateCount++

		case errors.Is(err, parser.ErrInvalidLink):
			errorCount++

			system.Warning(ctx, "не поддерживаемая ссылка", u)
		case err != nil:
			errorCount++

			system.Error(ctx, err)
		default:
			loadedCount++
		}
	}

	system.Info(ctx,
		fmt.Sprintf(
			"всего: %d загружено: %d дубликаты: %d ошибки: %d",
			totalCount, loadedCount, duplicateCount, errorCount,
		))
}
