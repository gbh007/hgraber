package handler

import (
	"app/db"
	"app/file"
	"app/parser"
	"fmt"
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
			_ = db.UpdatePageSuccess(page.TitleID, page.PageNumber, true)
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

// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
func FirstHandle(u string) error {
	log.Println("начата обработка", u)
	p, ok, err := parser.Load(u)
	if err != nil {
		return err
	}
	_, err = db.InsertTitle(p.ParseName(), u, ok)
	if err != nil {
		return err
	}
	log.Println("завершена обработка", u)
	return nil
}

// Update обрабатывает данные тайтла (только недостающие)
func Update(title db.TitleShortInfo) error {
	log.Println("начата обработка", title.URL)
	p, ok, err := parser.Load(title.URL)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("not load")
	}
	if !title.Loaded {
		err = db.UpdateTitle(title.ID, p.ParseName(), ok)
		if err != nil {
			return err
		}
		log.Println("обновлено название", title.URL)
	}
	if !title.ParsedAuthors {
		err = db.UpdateTitleMeta(title.ID, db.AuthorsMetaType, p.ParseAuthors())
		if err != nil {
			return err
		}
		log.Println("обновлены авторы", title.URL)
	}
	if !title.ParsedTags {
		err = db.UpdateTitleMeta(title.ID, db.TagsMetaType, p.ParseTags())
		if err != nil {
			return err
		}
		log.Println("обновлены теги", title.URL)
	}
	if !title.ParsedCharacters {
		err = db.UpdateTitleMeta(title.ID, db.CharactersMetaType, p.ParseCharacters())
		if err != nil {
			return err
		}
		log.Println("обновлены персонажи", title.URL)
	}
	if !title.ParsedCategories {
		err = db.UpdateTitleMeta(title.ID, db.CategoriesMetaType, p.ParseCategories())
		if err != nil {
			return err
		}
		log.Println("обновлены категории", title.URL)
	}
	if !title.ParsedGroups {
		err = db.UpdateTitleMeta(title.ID, db.GroupsMetaType, p.ParseGroups())
		if err != nil {
			return err
		}
		log.Println("обновлены группы", title.URL)
	}
	if !title.ParsedLanguages {
		err = db.UpdateTitleMeta(title.ID, db.LanguagesMetaType, p.ParseLanguages())
		if err != nil {
			return err
		}
		log.Println("обновлены языки", title.URL)
	}
	if !title.ParsedParodies {
		err = db.UpdateTitleMeta(title.ID, db.ParodiesMetaType, p.ParseParodies())
		if err != nil {
			return err
		}
		log.Println("обновлены пародии", title.URL)
	}
	if !title.ParsedPage {
		pp := true
		pages := p.ParsePages()
		for _, page := range pages {
			if db.InsertPage(title.ID, page.Ext, page.URL, page.Number) != nil {
				pp = false
			}
		}
		err = db.UpdateTitleParsedPage(title.ID, len(pages), pp && (len(pages) > 0))
		if err != nil {
			return err
		}
		log.Println("обновлены страницы", title.URL)
	}
	log.Println("завершена обработка", title.URL)
	return nil
}
