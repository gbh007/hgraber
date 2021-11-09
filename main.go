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

	go func() {
		timer := time.NewTimer(time.Minute)
		for range timer.C {
			handler.AddUnloadedPagesToQueue()
			time.Sleep(time.Second)
			handler.FileWait()
			timer.Reset(time.Minute)
		}
	}()

	go func() {
		timer := time.NewTicker(time.Minute)
		for range timer.C {
			for _, t := range db.SelectUnloadTitles() {
				_ = handler.Update(t)
			}
		}
	}()

	go func() {
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
	}()

	done := web.Run(fmt.Sprintf(":%d", *webPort))
	<-done
}
