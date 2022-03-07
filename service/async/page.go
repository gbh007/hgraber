package async

import (
	"app/service/fileStorage"
	"app/service/jdb"
	"app/system"
	"context"
	"sync"
	"time"
)

type qPage struct {
	TitleID    int
	PageNumber int
	URL        string
	Ext        string
}

var _pl *PageLoader

type PageLoader struct {
	queue  chan qPage
	ctx    context.Context
	inWork int
	mutex  *sync.RWMutex
}

func GetPageLoader() *PageLoader {
	return _pl
}

func (pl *PageLoader) workStarted() {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	pl.inWork++
}

func (pl *PageLoader) workEnded() {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	pl.inWork--
}

func (pl *PageLoader) InQueueCount() int {
	return len(pl.queue)
}

func (pl *PageLoader) InWorkCount() int {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	return pl.inWork
}

func (pl *PageLoader) handle(page qPage) {
	system.AddWaiting(pl.ctx)
	defer system.DoneWaiting(pl.ctx)

	pl.workStarted()
	defer pl.workEnded()

	err := fileStorage.DownloadTitlePage(pl.ctx, page.TitleID, page.PageNumber, page.URL, page.Ext)
	if err == nil {
		_ = jdb.Get().UpdatePageSuccess(pl.ctx, page.TitleID, page.PageNumber, true)
	}
}

func (pl *PageLoader) runQueueHandler() {
	defer system.Debug(pl.ctx, "PageLoader-handler остановлен")
	for page := range pl.queue {
		if system.IsAliveContext(pl.ctx) != nil {
			return
		}
		pl.handle(page)
	}
}

func (pl *PageLoader) Run() {
	for i := 0; i < pageHandlersCount; i++ {
		go pl.runQueueHandler()
	}
	system.Info(pl.ctx, "PageLoader запущен")
	defer system.Info(pl.ctx, "PageLoader остановлен")
	timer := time.NewTicker(time.Minute)
	for {
		select {
		case <-pl.ctx.Done():
			return
		case <-timer.C:
			if len(pl.queue) > 0 || pl.InWorkCount() > 0 {
				continue
			}
			for _, p := range jdb.Get().GetUnsuccessedPages(pl.ctx) {
				select {
				case <-pl.ctx.Done():
					return
				default:
				}
				pl.queue <- qPage{
					TitleID:    p.TitleID,
					PageNumber: p.PageNumber,
					URL:        p.URL,
					Ext:        p.Ext,
				}
			}
		}
	}
}
