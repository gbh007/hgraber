package main

import (
	"app/db"
	"app/parser"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func escapeFileName(n string) string {
	const r = " "
	if len([]rune(n)) > 200 {
		n = string([]rune(n)[:200])
	}
	for _, e := range []string{`\`, `/`, `|`, `:`, `"`, `*`, `?`, `<`, `>`} {
		n = strings.ReplaceAll(n, e, r)
	}
	return n
}

func main() {
	lf, err := os.Create("log.txt")
	if err != nil {
		log.Println(err)
		return
	}
	log.SetOutput(io.MultiWriter(os.Stderr, lf))

	db.Connect()

	p := parser.Parser_IMHENTAI_XXX{}
	log.Println(p.Load("https://imhentai.xxx/gallery/692183/"))
	log.Println(p.ParseAuthors())
	log.Println(p.ParseCharacters())
	log.Println(p.ParseName())
	log.Println(p.ParsePages())
	log.Println(p.ParseTags())

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
	notComplete := make([]string, 0)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		wg.Add(1)
		go func(u string) {
			log.Println(p.Load(u))
			log.Println(p.ParseAuthors())
			log.Println(p.ParseCharacters())
			log.Println(p.ParseName())
			log.Println(p.ParsePages())
			log.Println(p.ParseTags())
			notComplete = append(notComplete, u)
			wg.Done()
		}(sc.Text())
	}
	wg.Wait()
	f.Close()
	f, err = os.Create("task.txt")
	if err != nil {
		log.Println(err)
		return
	}
	for _, u := range notComplete {
		fmt.Fprintln(f, u)
	}
	f.Close()
	time.Sleep(time.Second)
	// load("https://imhentai.xxx/gallery/686547/")
}
