package async

import (
	"app/service/jdb"
	"app/system"
	"context"
	"sync"
	"time"
)

const (
	// pageQueueSize максимальный размер очереди для загрузки файлов страницы
	pageQueueSize = 10000
	// pageHandlersCount количество одновременно запущенных загрузчиков страниц
	pageHandlersCount = 10

	titleQueueSize     = 10000
	titleHandlersCount = 10
)

func Init(parentCtx context.Context, dbFilename string) {
	ctx := system.NewSystemContext(parentCtx, "Async")
	system.Info(ctx, "Запуск асинхронных обработчиков")

	_pl = &PageLoader{
		queue: make(chan qPage, pageQueueSize),
		ctx:   system.NewSystemContext(ctx, "Page-loader"),
		mutex: &sync.RWMutex{},
	}
	go _pl.Run()

	go autosaveDB(ctx, dbFilename)

	_tl = &TitleLoader{
		queue: make(chan jdb.Title, titleQueueSize),
		ctx:   system.NewSystemContext(ctx, "Title-loader"),
		mutex: &sync.RWMutex{},
	}
	go _tl.Run()

	system.Info(ctx, "Запущены асинхронные обработчики")
}

func autosaveDB(parentCtx context.Context, filename string) {
	ctx := system.NewSystemContext(parentCtx, "DB-autosave")
	system.Info(ctx, "autosaveDB запущен")
	defer system.Info(ctx, "autosaveDB остановлен")
	timer := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if jdb.Get().Save(ctx, filename) == nil {
				system.Debug(ctx, "Автосохранение прошло успешно")
			}
		}
	}
}
