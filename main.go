package main

import (
	"app/db"
	"app/service/jdb"
	"app/service/titleHandler"
	"app/service/webServer"
	"app/system"
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {

	webPort := flag.Int("p", 8080, "порт веб сервера")
	onlyView := flag.Bool("v", false, "режим только просмотра")
	export := flag.Bool("e", false, "экспортировать данные и выйти")
	disableStdErr := flag.Bool("no-stderr", false, "отключить стандартный поток ошибок")
	disableFileErr := flag.Bool("no-stdfile", false, "отключить поток ошибок в файл")
	enableAppendFileErr := flag.Bool("stdfile-append", false, "режим дозаписи файла потока ошибок")
	fileStorage := flag.String("fs", "loads", "директория для данных")
	debugMode := flag.Bool("debug", false, "активировать режим отладки")
	flag.Parse()

	system.Init(system.LogConfig{
		EnableFile:   !*disableFileErr,
		AppendMode:   *enableAppendFileErr,
		EnableStdErr: !*disableStdErr,
	})

	notifyCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	mainContext := system.NewSystemContext(notifyCtx, "MAIN")

	if *debugMode {
		system.EnableDebug(mainContext)
	}

	err := system.SetFileStoragePath(mainContext, *fileStorage)
	if err != nil {
		os.Exit(1)
	}

	err = db.Connect(mainContext)
	if err != nil {
		os.Exit(2)
	}

	if *export {
		exportData(mainContext)
	}

	if !*onlyView {
		go loadPages(mainContext)
		go completeTitle(mainContext)
		go parseTaskFile(mainContext)
		system.Info(mainContext, "Запущены асинхронные обработчики")
	}

	webServer.Run(mainContext, fmt.Sprintf(":%d", *webPort))

	<-mainContext.Done()
	system.Info(mainContext, "Завершение работы, ожидание завершения процессов")
	<-system.WaitingChan(mainContext)
	system.Info(mainContext, "Процессы завершены, выход")
}

func exportData(ctx context.Context) {
	system.Info(ctx, "Экспорт начат")
	exporter := jdb.New()
	system.Info(ctx, "Конвертирование данных")
	exporter.FetchFromSQL(ctx)
	system.Info(ctx, "Сохранение данных")
	_ = exporter.Save(ctx, fmt.Sprintf("exported-%s.json", time.Now().Format("2006-01-02-150405")))
	system.Info(ctx, "Экспорт завершен")
	os.Exit(0)

}

func loadPages(ctx context.Context) {
	titleHandler.Init(ctx)
	timer := time.NewTimer(time.Minute)
	for range timer.C {
		titleHandler.AddUnloadedPagesToQueue(ctx)
		time.Sleep(time.Second)
		titleHandler.FileWait()
		timer.Reset(time.Minute)
	}
}

func completeTitle(ctx context.Context) {
	timer := time.NewTicker(time.Minute)
	for range timer.C {
		for _, t := range db.SelectUnloadTitles(ctx) {
			_ = titleHandler.Update(ctx, t)
		}
	}
}

func parseTaskFile(ctx context.Context) {
	f, err := os.Open("task.txt")
	defer system.IfErrFunc(ctx, f.Close)
	if err != nil {
		system.Error(ctx, err)
		return
	}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		_ = titleHandler.FirstHandle(ctx, sc.Text())
	}
}
