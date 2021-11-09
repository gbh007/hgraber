package main

import (
	"app/db"
	"app/handler"
	"app/web"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {

	webPort := flag.Int("p", 8080, "порт веб сервера")
	onlyView := flag.Bool("v", false, "режим только просмотра")
	flag.IntVar(&web.PageLimit, "pl", 12, "количество тайтлов на странице")
	flag.Parse()

	lf, err := os.Create("log.txt")
	if err != nil {
		log.Println(err)
		return
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(io.MultiWriter(os.Stderr, lf))

	err = db.Connect()
	if err != nil {
		log.Println(err)
		return
	}

	if !*onlyView {
		go loadPages()
		go completeTitle()
		go parseTaskFile()
	}

	done := web.Run(fmt.Sprintf(":%d", *webPort))
	<-done
}

func loadPages() {
	timer := time.NewTimer(time.Minute)
	for range timer.C {
		handler.AddUnloadedPagesToQueue()
		time.Sleep(time.Second)
		handler.FileWait()
		timer.Reset(time.Minute)
	}
}

func completeTitle() {
	timer := time.NewTicker(time.Minute)
	for range timer.C {
		for _, t := range db.SelectUnloadTitles() {
			_ = handler.Update(t)
		}
	}
}

func parseTaskFile() {
	f, err := os.Open("task.txt")
	if err != nil {
		log.Println(err)
		return
	}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		_ = handler.FirstHandle(sc.Text())
	}
	f.Close()
}
