package handler

import (
	"app/db"
	"app/file"
	"app/parser"
	"log"
	"sync"
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
func HandleFull(u string) {
	log.Println("начата обработка", u)
	p, ok := parser.Load(u)
	id, err := db.InsertTitle(p.ParseName(), u, ok)
	if err != nil {
		return
	}
	if ok {
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
}
