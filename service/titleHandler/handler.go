package titleHandler

import (
	"app/service/fileStorage"
	"app/service/jdb"
	"app/service/parser"
	"app/system"
	"context"
	"fmt"
	"sync"
	"time"
)

// maxQueueSize максимальный размер очереди для загрузки файлов
const maxQueueSize = 100000

type qPage struct {
	TitleID    int
	PageNumber int
	URL        string
	Ext        string
}

// fileQueue очередь для загрузки файлов
var fileQueue chan qPage

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

func Init(parentCtx context.Context) {
	ctx := system.NewSystemContext(parentCtx, "FILE-HANDLE")
	fileQueue = make(chan qPage, maxQueueSize)
	for i := 0; i < maxFileHandlersCount; i++ {
		go handleFileQueue(ctx)
	}
}

// handleFileQueue обработчик файловой очереди
func handleFileQueue(ctx context.Context) {
	for page := range fileQueue {
		fileWG.Add(1)
		err := fileStorage.DownloadTitlePage(ctx, page.TitleID, page.PageNumber, page.URL, page.Ext)
		if err == nil {
			_ = jdb.Get().UpdatePageSuccess(ctx, page.TitleID, page.PageNumber, true)
		}
		fileWG.Done()
	}
}

// AddUnloadedPagesToQueue добавляет незагруженные страницы в очередь
func AddUnloadedPagesToQueue(ctx context.Context) {
	for _, p := range jdb.Get().GetUnsuccessedPages(ctx) {
		fileQueue <- qPage{
			TitleID:    p.TitleID,
			PageNumber: p.PageNumber,
			URL:        p.URL,
			Ext:        p.Ext,
		}
	}
}

// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
func FirstHandle(ctx context.Context, u string) error {
	system.Info(ctx, "начата обработка", u)
	p, ok, err := parser.Load(ctx, u)
	if err != nil {
		return err
	}
	_, err = jdb.Get().NewTitle(ctx, p.ParseName(ctx), u, ok)
	if err != nil {
		return err
	}
	system.Info(ctx, "завершена обработка", u)
	return nil
}

// Update обрабатывает данные тайтла (только недостающие)
func Update(ctx context.Context, title jdb.Title) error {
	system.Info(ctx, "начата обработка", title.ID, title.URL)
	p, ok, err := parser.Load(ctx, title.URL)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("not load")
	}
	if !title.Data.Parsed.Name {
		err = jdb.Get().UpdateTitleName(ctx, title.ID, p.ParseName(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлено название", title.ID, title.URL)
	}
	if !title.Data.Parsed.Authors {
		err = jdb.Get().UpdateTitleAuthors(ctx, title.ID, p.ParseAuthors(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены авторы", title.ID, title.URL)
	}
	if !title.Data.Parsed.Tags {
		err = jdb.Get().UpdateTitleTags(ctx, title.ID, p.ParseTags(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены теги", title.ID, title.URL)
	}
	if !title.Data.Parsed.Characters {
		err = jdb.Get().UpdateTitleCharacters(ctx, title.ID, p.ParseCharacters(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены персонажи", title.ID, title.URL)
	}
	if !title.Data.Parsed.Categories {
		err = jdb.Get().UpdateTitleCategories(ctx, title.ID, p.ParseCategories(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены категории", title.ID, title.URL)
	}
	if !title.Data.Parsed.Groups {
		err = jdb.Get().UpdateTitleGroups(ctx, title.ID, p.ParseGroups(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены группы", title.ID, title.URL)
	}
	if !title.Data.Parsed.Languages {
		err = jdb.Get().UpdateTitleLanguages(ctx, title.ID, p.ParseLanguages(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены языки", title.ID, title.URL)
	}
	if !title.Data.Parsed.Parodies {
		err = jdb.Get().UpdateTitleParodies(ctx, title.ID, p.ParseParodies(ctx))
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены пародии", title.ID, title.URL)
	}
	if !title.Data.Parsed.Page {
		pages := p.ParsePages(ctx)
		pagesDB := make([]jdb.Page, len(pages))
		for i, page := range pages {
			pagesDB[i] = jdb.Page{
				URL: page.URL,
				Ext: page.Ext,
			}
		}

		err = jdb.Get().UpdateTitlePages(ctx, title.ID, pagesDB)
		if err != nil {
			return err
		}
		system.Info(ctx, "обновлены страницы", title.ID, title.URL)
	}
	system.Info(ctx, "завершена обработка", title.ID, title.URL)
	return nil
}
