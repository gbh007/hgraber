package main

import (
	"app/db"
	"app/handler"
	"bufio"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	lf, err := os.Create("log.txt")
	if err != nil {
		log.Println(err)
		return
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(io.MultiWriter(os.Stderr, lf))

	db.Connect()

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
