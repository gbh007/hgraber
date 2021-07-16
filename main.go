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

func handle(u string) {
	log.Println("начата обработка", u)
	p := parser.Parser_IMHENTAI_XXX{}
	ok := p.Load(u)
	id, err := db.InsertTitle(p.ParseName(), u, ok)
	if err != nil {
		return
	}
	err = db.UpdateTitleAuthors(id, p.ParseAuthors())
	if err != nil {
		return
	}
	err = db.UpdateTitleTags(id, p.ParseTags())
	if err != nil {
		return
	}
	err = db.UpdateTitleCharacters(id, p.ParseCharacters())
	if err != nil {
		return
	}
	for _, page := range p.ParsePages() {
		db.InsertPage(id, page.Name, page.URL, page.Number)
	}
	log.Println("завершена обработка", u)
}

func main() {
	lf, err := os.Create("log.txt")
	if err != nil {
		log.Println(err)
		return
	}
	log.SetOutput(io.MultiWriter(os.Stderr, lf))

	db.Connect()

	handle("https://imhentai.xxx/gallery/692183/")
	/*
		p := parser.Parser_IMHENTAI_XXX{}
		log.Println(p.Load("https://imhentai.xxx/gallery/692183/"))
		log.Println(p.ParseAuthors())
		log.Println(p.ParseCharacters())
		log.Println(p.ParseName())
		log.Println(p.ParsePages())
		log.Println(p.ParseTags())

		log.Println(db.GetTagID("test"))

		id, err := db.InsertTitle("123", "https://imhentai.xxx/gallery/692183/", true)
		log.Println(id, err)
		log.Println(db.UpdateTitleTags(id, []string{"t1", "t2", "t3"}))
		log.Println(db.UpdateTitleAuthors(id, []string{"a1", "a2", "a3"}))
		log.Println(db.UpdateTitleCharacters(id, []string{"c1", "c2", "c3"}))
		log.Println(db.InsertPage(id, "1.tft", "http://", 3))
		log.Println(db.InsertPage(id, "3.tft", "http://", 3))
	*/
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
			handle(u)
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
