package main

import (
	"app/db"
	"app/handler"
	"app/jdb"
	"app/system"
	"app/webgin"
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
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
	flag.Parse()

	system.Init(system.LogConfig{
		EnableFile:   !*disableFileErr,
		AppendMode:   *enableAppendFileErr,
		EnableStdErr: !*disableStdErr,
	})

	mainContext := system.NewSystemContext(context.Background(), "MAIN")

	err := system.SetFileStoragePath(mainContext, *fileStorage)
	if err != nil {
		os.Exit(1)
	}

	err = db.Connect(mainContext)
	if err != nil {
		system.Error(mainContext, err)
		os.Exit(2)
	}

	if *export {
		system.Info(mainContext, "Экспорт начат")
		exporter := jdb.New()
		system.Info(mainContext, "Конвертирование данных")
		exporter.FetchFromSQL(mainContext)
		system.Info(mainContext, "Сохранение данных")
		_ = exporter.Save(mainContext, fmt.Sprintf("exported-%s.json", time.Now().Format("2006-01-02-150405")))
		system.Info(mainContext, "Экспорт завершен")
		os.Exit(0)
	}

	if !*onlyView {
		go loadPages(mainContext)
		go completeTitle(mainContext)
		go parseTaskFile(mainContext)
		system.Info(mainContext, "Запущены асинхронные обработчики")
	}

	system.Info(mainContext, "Запущен веб сервер")
	done := webgin.Run(mainContext, fmt.Sprintf(":%d", *webPort))
	<-done
}

func loadPages(ctx context.Context) {
	handler.Init(ctx)
	timer := time.NewTimer(time.Minute)
	for range timer.C {
		handler.AddUnloadedPagesToQueue(ctx)
		time.Sleep(time.Second)
		handler.FileWait()
		timer.Reset(time.Minute)
	}
}

func completeTitle(ctx context.Context) {
	timer := time.NewTicker(time.Minute)
	for range timer.C {
		for _, t := range db.SelectUnloadTitles(ctx) {
			_ = handler.Update(ctx, t)
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
		_ = handler.FirstHandle(ctx, sc.Text())
	}
}
