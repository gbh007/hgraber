package handler

import (
	"app/db"
	"app/file"
	"app/parser"
	"app/system"
	"fmt"
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
	ctx := system.NewSystemContext("FILE-HANDLE")
	fileQueue = make(chan db.Page, maxQueueSize)
	for i := 0; i < maxFileHandlersCount; i++ {
		go handleFileQueue(ctx)
	}
}

// handleFileQueue обработчик файловой очереди
func handleFileQueue(ctx system.Context) {
	for page := range fileQueue {
		fileWG.Add(1)
		err := file.Load(ctx, page.TitleID, page.PageNumber, page.URL, page.Ext)
		if err == nil {
			_ = db.UpdatePageSuccess(ctx, page.TitleID, page.PageNumber, true)
		}
		fileWG.Done()
	}
}

// AddUnloadedPagesToQueue добавляет незагруженные страницы в очередь
func AddUnloadedPagesToQueue(ctx system.Context) {
	for _, p := range db.SelectUnsuccessPages(ctx) {
		fileQueue <- p
	}
}

// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
func FirstHandle(ctx system.Context, u string) error {
	system.Info(ctx, "начата обработка", u)
	p, ok, err := parser.Load(ctx, u)
	if err != nil {
		return err
	}
	_, err = db.InsertTitle(ctx, p.ParseName(ctx), u, ok)
	if err != nil {
		return err
	}
	system.Info(ctx, "завершена обработка", u)
	return nil
}

// Update обрабатывает данные тайтла (только недостающие)
func Update(ctx system.Context, title db.TitleShortInfo) error {
	system.Info(ctx, "начата обработка", title.URL)
	p, ok, err := parser.Load(ctx, title.URL)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("not load")
	}
	if !title.Loaded {
		err = db.UpdateTitle(ctx, title.ID, p.ParseName(ctx), ok)
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлено название", title.URL)
	}
	if !title.ParsedAuthors {
		err = db.UpdateTitleMeta(ctx, title.ID, db.AuthorsMetaType, p.ParseAuthors(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены авторы", title.URL)
	}
	if !title.ParsedTags {
		err = db.UpdateTitleMeta(ctx, title.ID, db.TagsMetaType, p.ParseTags(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены теги", title.URL)
	}
	if !title.ParsedCharacters {
		err = db.UpdateTitleMeta(ctx, title.ID, db.CharactersMetaType, p.ParseCharacters(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены персонажи", title.URL)
	}
	if !title.ParsedCategories {
		err = db.UpdateTitleMeta(ctx, title.ID, db.CategoriesMetaType, p.ParseCategories(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены категории", title.URL)
	}
	if !title.ParsedGroups {
		err = db.UpdateTitleMeta(ctx, title.ID, db.GroupsMetaType, p.ParseGroups(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены группы", title.URL)
	}
	if !title.ParsedLanguages {
		err = db.UpdateTitleMeta(ctx, title.ID, db.LanguagesMetaType, p.ParseLanguages(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены языки", title.URL)
	}
	if !title.ParsedParodies {
		err = db.UpdateTitleMeta(ctx, title.ID, db.ParodiesMetaType, p.ParseParodies(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены пародии", title.URL)
	}
	if !title.ParsedPage {
		pp := true
		pages := p.ParsePages(ctx)
		for _, page := range pages {
			if db.InsertPage(ctx, title.ID, page.Ext, page.URL, page.Number) != nil {
				pp = false
			}
		}
		err = db.UpdateTitleParsedPage(ctx, title.ID, len(pages), pp && (len(pages) > 0))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены страницы", title.URL)
	}
	system.Info(ctx, "завершена обработка", title.URL)
	return nil
}
