package handler

import (
	"app/db"
	"app/file"
	"app/parser"
	"app/system/clog"
	"app/system/coreContext"
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
	ctx := coreContext.NewSystemContext()
	ctx.SetRequestID("FILE-HANDLE")
	fileQueue = make(chan db.Page, maxQueueSize)
	for i := 0; i < maxFileHandlersCount; i++ {
		go handleFileQueue(ctx)
	}
}

// handleFileQueue обработчик файловой очереди
func handleFileQueue(ctx coreContext.CoreContext) {
	for page := range fileQueue {
		fileWG.Add(1)
		err := file.Load(page.TitleID, page.PageNumber, page.URL, page.Ext)
		if err == nil {
			_ = db.UpdatePageSuccess(ctx, page.TitleID, page.PageNumber, true)
		}
		fileWG.Done()
	}
}

// AddUnloadedPagesToQueue добавляет незагруженные страницы в очередь
func AddUnloadedPagesToQueue(ctx coreContext.CoreContext) {
	for _, p := range db.SelectUnsuccessPages(ctx) {
		fileQueue <- p
	}
}

// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
func FirstHandle(ctx coreContext.CoreContext, u string) error {
	clog.Info(ctx, "начата обработка", u)
	p, ok, err := parser.Load(u)
	if err != nil {
		return err
	}
	_, err = db.InsertTitle(ctx, p.ParseName(), u, ok)
	if err != nil {
		return err
	}
	clog.Info(ctx, "завершена обработка", u)
	return nil
}

// Update обрабатывает данные тайтла (только недостающие)
func Update(ctx coreContext.CoreContext, title db.TitleShortInfo) error {
	log.Println("начата обработка", title.URL)
	p, ok, err := parser.Load(title.URL)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("not load")
	}
	if !title.Loaded {
		err = db.UpdateTitle(ctx, title.ID, p.ParseName(), ok)
		if err != nil {
			return err
		}
		log.Println("обновлено название", title.URL)
	}
	if !title.ParsedAuthors {
		err = db.UpdateTitleMeta(ctx, title.ID, db.AuthorsMetaType, p.ParseAuthors())
		if err != nil {
			return err
		}
		log.Println("обновлены авторы", title.URL)
	}
	if !title.ParsedTags {
		err = db.UpdateTitleMeta(ctx, title.ID, db.TagsMetaType, p.ParseTags())
		if err != nil {
			return err
		}
		log.Println("обновлены теги", title.URL)
	}
	if !title.ParsedCharacters {
		err = db.UpdateTitleMeta(ctx, title.ID, db.CharactersMetaType, p.ParseCharacters())
		if err != nil {
			return err
		}
		log.Println("обновлены персонажи", title.URL)
	}
	if !title.ParsedCategories {
		err = db.UpdateTitleMeta(ctx, title.ID, db.CategoriesMetaType, p.ParseCategories())
		if err != nil {
			return err
		}
		log.Println("обновлены категории", title.URL)
	}
	if !title.ParsedGroups {
		err = db.UpdateTitleMeta(ctx, title.ID, db.GroupsMetaType, p.ParseGroups())
		if err != nil {
			return err
		}
		log.Println("обновлены группы", title.URL)
	}
	if !title.ParsedLanguages {
		err = db.UpdateTitleMeta(ctx, title.ID, db.LanguagesMetaType, p.ParseLanguages())
		if err != nil {
			return err
		}
		log.Println("обновлены языки", title.URL)
	}
	if !title.ParsedParodies {
		err = db.UpdateTitleMeta(ctx, title.ID, db.ParodiesMetaType, p.ParseParodies())
		if err != nil {
			return err
		}
		log.Println("обновлены пародии", title.URL)
	}
	if !title.ParsedPage {
		pp := true
		pages := p.ParsePages()
		for _, page := range pages {
			if db.InsertPage(ctx, title.ID, page.Ext, page.URL, page.Number) != nil {
				pp = false
			}
		}
		err = db.UpdateTitleParsedPage(ctx, title.ID, len(pages), pp && (len(pages) > 0))
		if err != nil {
			return err
		}
		log.Println("обновлены страницы", title.URL)
	}
	log.Println("завершена обработка", title.URL)
	return nil
}
