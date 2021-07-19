package handler

import (
	"app/db"
	"app/file"
	"app/parser"
	"log"
	"sync"
	"time"
)

// maxQueueSize максимальный размер очереди для загрузки файлов
const maxQueueSize = 100000

// fileQueue очередь для загрузки файлов
var fileQueue chan db.Page

// maxFileHandlersCount количество одновременно запущенных файловых загрузчиков
const maxFileHandlersCount = 10

// fileWG группа ожидания для фаловых обработчиков
var fileWG *sync.WaitGroup = &sync.WaitGroup{}

// FileWait ожидает завершения файловых обработчиков
func FileWait() {
	for len(fileQueue) > 0 {
		time.Sleep(time.Second)
	}
	fileWG.Wait()
}

func init() {
	fileQueue = make(chan db.Page, maxQueueSize)
	for i := 0; i < maxFileHandlersCount; i++ {
		go handleFileQueue()
	}
}

// handleFileQueue обработчик файловой очереди
func handleFileQueue() {
	for page := range fileQueue {
		fileWG.Add(1)
		err := file.Load(page.TitleID, page.PageNumber, page.URL, page.Ext)
		if err == nil {
			db.UpdatePageSuccess(page.TitleID, page.PageNumber, true)
		}
		fileWG.Done()
	}
}

// AddUnloadedPagesToQueue добавляет незагруженные страницы в очередь
func AddUnloadedPagesToQueue() {
	for _, p := range db.SelectUnsuccessPages() {
		fileQueue <- p
	}
}

// HandleFull обрабатывает данные тайтла (новое добавление)
func HandleFull(u string) error {
	log.Println("начата обработка", u)
	p, ok, err := parser.Load(u)
	if err != nil {
		return err
	}
	id, err := db.InsertTitle(p.ParseName(), u, ok)
	if err != nil {
		return err
	}
	if ok {
		err = db.UpdateTitleAuthors(id, p.ParseAuthors())
		if err != nil {
			return err
		}
		err = db.UpdateTitleTags(id, p.ParseTags())
		if err != nil {
			return err
		}
		err = db.UpdateTitleCharacters(id, p.ParseCharacters())
		if err != nil {
			return err
		}
		pp := true
		pages := p.ParsePages()
		for _, page := range pages {
			if db.InsertPage(id, page.Ext, page.URL, page.Number) != nil {
				pp = false
			}
		}
		db.UpdateTitleParsedPage(id, len(pages), pp)
	}
	log.Println("завершена обработка", u)
	return nil
}

// UpdateFull обрабатывает данные тайтла (переобработка)
func UpdateFull(id int, u string) error {
	log.Println("начата обработка", u)
	p, ok, err := parser.Load(u)
	if err != nil {
		return err
	}
	err = db.UpdateTitle(id, p.ParseName(), ok)
	if err != nil {
		return err
	}
	if ok {
		err = db.UpdateTitleAuthors(id, p.ParseAuthors())
		if err != nil {
			return err
		}
		err = db.UpdateTitleTags(id, p.ParseTags())
		if err != nil {
			return err
		}
		err = db.UpdateTitleCharacters(id, p.ParseCharacters())
		if err != nil {
			return err
		}
		pp := true
		pages := p.ParsePages()
		for _, page := range pages {
			if db.InsertPage(id, page.Ext, page.URL, page.Number) != nil {
				pp = false
			}
		}
		db.UpdateTitleParsedPage(id, len(pages), pp)
	}
	log.Println("завершена обработка", u)
	return nil
}
