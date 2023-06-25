package async

import (
	"app/service/jdb"
	"app/service/titleHandler"
	"app/system"
	"context"
	"sync"
	"time"
)

var _tl *TitleLoader

type TitleLoader struct {
	queue  chan jdb.Title
	ctx    context.Context
	inWork int
	mutex  *sync.RWMutex
}

func GetTitleLoader() *TitleLoader {
	return _tl
}

func (tl *TitleLoader) workStarted() {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()
	tl.inWork++
}

func (tl *TitleLoader) workEnded() {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()
	tl.inWork--
}

func (tl *TitleLoader) InQueueCount() int {
	return len(tl.queue)
}

func (tl *TitleLoader) InWorkCount() int {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()
	return tl.inWork
}

func (tl *TitleLoader) handle(title jdb.Title) {
	system.AddWaiting(tl.ctx)
	defer system.DoneWaiting(tl.ctx)

	tl.workStarted()
	defer tl.workEnded()

	titleHandler.Update(tl.ctx, title)
}

func (tl *TitleLoader) runQueueHandler() {
	defer system.Debug(tl.ctx, "TitleLoader-handler остановлен")
	for page := range tl.queue {
		if system.IsAliveContext(tl.ctx) != nil {
			return
		}
		tl.handle(page)
	}
}

func (tl *TitleLoader) Run() {
	for i := 0; i < titleHandlersCount; i++ {
		go tl.runQueueHandler()
	}

	system.Info(tl.ctx, "TitleLoader запущен")
	defer system.Info(tl.ctx, "TitleLoader остановлен")

	timer := time.NewTicker(titleInterval)

	for {
		select {
		case <-tl.ctx.Done():
			return
		case <-timer.C:
			if len(tl.queue) > 0 || tl.InWorkCount() > 0 {
				continue
			}
			for _, title := range jdb.Get().GetUnloadedTitles(tl.ctx) {
				select {
				case <-tl.ctx.Done():
					return
				default:
				}
				tl.queue <- title
			}
		}
	}
}
