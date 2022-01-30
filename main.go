package main

import (
	"app/db"
	"app/handler"
	"app/system"
	"app/webgin"
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	webPort := flag.Int("p", 8080, "порт веб сервера")
	onlyView := flag.Bool("v", false, "режим только просмотра")
	flag.Parse()

	mainContext := system.NewSystemContext("MAIN")

	err := db.Connect(mainContext)
	if err != nil {
		system.Error(mainContext, err)
		return
	}

	if !*onlyView {
		go loadPages(mainContext)
		go completeTitle(mainContext)
		go parseTaskFile(mainContext)
		system.Info(mainContext, "Запущены асинхронные обработчики")
	}

	done := webgin.Run(mainContext, fmt.Sprintf(":%d", *webPort))
	<-done
}

func loadPages(ctx system.Context) {
	timer := time.NewTimer(time.Minute)
	for range timer.C {
		handler.AddUnloadedPagesToQueue(ctx)
		time.Sleep(time.Second)
		handler.FileWait()
		timer.Reset(time.Minute)
	}
}

func completeTitle(ctx system.Context) {
	timer := time.NewTicker(time.Minute)
	for range timer.C {
		for _, t := range db.SelectUnloadTitles(ctx) {
			_ = handler.Update(ctx, t)
		}
	}
}

func parseTaskFile(ctx system.Context) {
	f, err := os.Open("task.txt")
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
	f.Close()
}
