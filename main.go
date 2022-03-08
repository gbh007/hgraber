package main

import (
	"app/service/async"
	"app/service/jdb"
	"app/service/titleHandler"
	"app/service/webServer"
	"app/system"
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// базовые опции
	webPort := flag.Int("p", 8080, "порт веб сервера")
	onlyView := flag.Bool("v", false, "режим только просмотра")
	export := flag.Bool("e", false, "экспортировать данные и выйти")

	// потоки логирования
	disableStdErr := flag.Bool("no-stderr", false, "отключить стандартный поток ошибок")
	disableFileErr := flag.Bool("no-stdfile", false, "отключить поток ошибок в файл")
	enableAppendFileErr := flag.Bool("stdfile-append", false, "режим дозаписи файла потока ошибок")

	// размещение данных
	fileStorage := flag.String("fs", "loads", "директория для данных")
	fileExport := flag.String("fe", "exported", "директория для экспорта файлов")
	dbFileName := flag.String("db", "db.json", "файл базы")

	// отладка
	debugMode := flag.Bool("debug", false, "активировать режим отладки (дебага)")
	debugCopyMode := flag.Bool("debug-copy", false, "включает при активном дебаге, информацию о копировании данных в памяти")
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
	if *debugCopyMode {
		jdb.EnableCopyStopwatch(mainContext)
	}
	if *debugFullpathMode {
		system.EnableFullpath(mainContext)
	}

	system.Debug(mainContext, "Версия", system.Version)
	system.Debug(mainContext, "Коммит", system.Commit)
	system.Debug(mainContext, "Собрано", system.BuildAt)

	system.Info(mainContext, "Инициализация базы")
	jdb.Init(mainContext)
	err := jdb.Get().Load(mainContext, *dbFileName)
	if err != nil {
		os.Exit(1)
	}
	system.Info(mainContext, "База загружена")

	err = system.SetFileStoragePath(mainContext, *fileStorage)
	if err != nil {
		os.Exit(2)
	}

	err = system.SetFileExportPath(mainContext, *fileExport)
	if err != nil {
		os.Exit(3)
	}

	if *export {
		exportData(mainContext)
		os.Exit(0)
	}

	if !*onlyView {
		go parseTaskFile(mainContext)
		async.Init(mainContext, *dbFileName)
	}

	webServer.Start(mainContext, fmt.Sprintf(":%d", *webPort))

	<-mainContext.Done()
	system.Info(mainContext, "Завершение работы, ожидание завершения процессов")
	<-system.WaitingChan(mainContext)
	system.Info(mainContext, "Процессы завершены")
	if jdb.Get().Save(mainContext, *dbFileName, false) == nil {
		system.Info(mainContext, "База сохранена")
	} else {
		system.Warning(mainContext, "База не сохранена")
	}
	system.Info(mainContext, "Выход")
}

func exportData(ctx context.Context) {
	system.Info(ctx, "Экспорт начат")
	exporter := jdb.Get()
	_ = exporter.Save(
		ctx,
		fmt.Sprintf(
			"%s/db-%s.json",
			system.GetFileExportPath(ctx),
			time.Now().Format("2006-01-02-150405"),
		),
		true,
	)
	system.Info(ctx, "Экспорт завершен")
}

func parseTaskFile(ctx context.Context) {
	f, err := os.Open("task.txt")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			system.Error(ctx, err)
		}
		return
	}
	defer system.IfErrFunc(ctx, f.Close)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		_ = titleHandler.FirstHandle(ctx, sc.Text())
	}
}
