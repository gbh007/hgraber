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
	"sync"
	"time"
)

func main() {

	webPort := flag.Int("p", 8080, "порт веб сервера")
	flag.Parse()

	lf, err := os.Create("log.txt")
	if err != nil {
		log.Println(err)
		return
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(io.MultiWriter(os.Stderr, lf))

	db.Connect()

	go func() {
		timer := time.NewTimer(time.Minute)
		for range timer.C {
			handler.AddUnloadedPagesToQueue()
			time.Sleep(time.Second)
			handler.FileWait()
			timer.Reset(time.Minute)
		}
	}()

	done := web.Run(fmt.Sprintf(":%d", *webPort))
	<-done
	if true {
		return
	}

	_, err = os.Stat("loads")
	if os.IsNotExist(err) {
		os.MkdirAll("loads", 0777)
	}
	f, err := os.Open("task.txt")
	if err != nil {
		log.Println(err)
		return
	}
	sc := bufio.NewScanner(f)
	wg := &sync.WaitGroup{}
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		wg.Add(1)
		go func(u string) {
			handler.HandleFull(u)
			wg.Done()
		}(sc.Text())
	}
	wg.Wait()
	f.Close()

	handler.AddUnloadedPagesToQueue()
	time.Sleep(time.Second)
	handler.FileWait()

	time.Sleep(time.Second)
}
