package main

import (
	"app/db"
	"app/handler"
	"app/system/clog"
	"app/system/coreContext"
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

	mainContext := coreContext.NewSystemContext()
	mainContext.SetRequestID("MAIN")

	err := db.Connect(mainContext)
	if err != nil {
		clog.Error(mainContext, err)
		return
	}

	if !*onlyView {
		go loadPages(mainContext)
		go completeTitle(mainContext)
		go parseTaskFile(mainContext)
	}

	done := webgin.Run(mainContext, fmt.Sprintf(":%d", *webPort))
	<-done
}

func loadPages(ctx coreContext.CoreContext) {
	timer := time.NewTimer(time.Minute)
	for range timer.C {
		handler.AddUnloadedPagesToQueue(ctx)
		time.Sleep(time.Second)
		handler.FileWait()
		timer.Reset(time.Minute)
	}
}

func completeTitle(ctx coreContext.CoreContext) {
	timer := time.NewTicker(time.Minute)
	for range timer.C {
		for _, t := range db.SelectUnloadTitles(ctx) {
			_ = handler.Update(ctx, t)
		}
	}
}

func parseTaskFile(ctx coreContext.CoreContext) {
	f, err := os.Open("task.txt")
	if err != nil {
		clog.Error(ctx, err)
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
