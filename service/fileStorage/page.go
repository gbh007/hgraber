package fileStorage

import (
	"app/system"
	"context"
	"time"
)

const (
	pageInterval = time.Second * 15
	// pageQueueSize максимальный размер очереди для загрузки файлов страницы
	pageQueueSize = 10000
	// pageHandlersCount количество одновременно запущенных загрузчиков страниц
	pageHandlersCount = 10
)

type qPage struct {
	TitleID    int
	PageNumber int
	URL        string
	Ext        string
}

func (s *Service) workStarted() {
	s.inWorkMutex.Lock()
	defer s.inWorkMutex.Unlock()

	s.inWork++
}

func (s *Service) workEnded() {
	s.inWorkMutex.Lock()
	defer s.inWorkMutex.Unlock()

	s.inWork--
}

func (s *Service) InQueueCount() int {
	return len(s.queue)
}

func (s *Service) InWorkCount() int {
	s.inWorkMutex.Lock()
	defer s.inWorkMutex.Unlock()

	return s.inWork
}

func (s *Service) handle(ctx context.Context, page qPage) {
	s.asyncPathWG.Add(1)
	defer s.asyncPathWG.Done()

	s.workStarted()
	defer s.workEnded()

	err := downloadTitlePage(ctx, page.TitleID, page.PageNumber, page.URL, page.Ext)
	if err == nil {
		updateErr := s.Storage.UpdatePageSuccess(ctx, page.TitleID, page.PageNumber, true)
		if updateErr != nil {
			system.Error(ctx, updateErr)
		}
	}
}

func (s *Service) runQueueHandler(ctx context.Context) {
	defer system.Debug(ctx, "PageLoader-handler остановлен")

	for page := range s.queue {
		if system.IsAliveContext(ctx) != nil {
			return
		}

		s.handle(ctx, page)
	}
}

func (s *Service) runFull(ctx context.Context) {
	for i := 0; i < pageHandlersCount; i++ {
		go s.runQueueHandler(ctx)
	}

	system.Info(ctx, "PageLoader запущен")
	defer system.Info(ctx, "PageLoader остановлен")

	timer := time.NewTicker(pageInterval)

	for {
		select {
		case <-ctx.Done():
			// Дожидаемся завершения всех подпроцессов
			s.asyncPathWG.Wait()

			return

		case <-timer.C:
			if len(s.queue) > 0 || s.InWorkCount() > 0 {
				continue
			}

			for _, p := range s.Storage.GetUnsuccessedPages(ctx) {
				select {
				case <-ctx.Done():
					// Дожидаемся завершения всех подпроцессов
					s.asyncPathWG.Wait()

					return

				default:
				}

				s.queue <- qPage{
					TitleID:    p.TitleID,
					PageNumber: p.PageNumber,
					URL:        p.URL,
					Ext:        p.Ext,
				}
			}
		}
	}
}
