package titleHandler

import (
	"app/storage/schema"
	"app/system"
	"context"
	"time"
)

const (
	titleInterval      = time.Second * 15
	titleQueueSize     = 10000
	titleHandlersCount = 10
)

func (s *Service) Name() string {
	return "title handler"
}

func (s *Service) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ctx := system.NewSystemContext(parentCtx, "Title-loader")

		s.runFull(ctx)
	}()

	return done, nil
}

func (s *Service) workStarted() {
	s.inWorkRunnersCountMutex.Lock()
	defer s.inWorkRunnersCountMutex.Unlock()

	s.inWorkRunnersCount++
}

func (s *Service) workEnded() {
	s.inWorkRunnersCountMutex.Lock()
	defer s.inWorkRunnersCountMutex.Unlock()

	s.inWorkRunnersCount--
}

func (s *Service) InQueueCount() int {
	return len(s.titleQueue)
}

func (s *Service) InWorkCount() int {
	s.inWorkRunnersCountMutex.Lock()
	defer s.inWorkRunnersCountMutex.Unlock()

	return s.inWorkRunnersCount
}

func (s *Service) handleTitleFromQueue(ctx context.Context, title schema.Title) {
	s.asyncPathWG.Add(1)
	defer s.asyncPathWG.Done()

	s.workStarted()
	defer s.workEnded()

	s.Update(ctx, title)
}

func (s *Service) runQueueHandler(ctx context.Context) {
	defer system.Debug(ctx, "TitleLoader-handler остановлен")

	for page := range s.titleQueue {
		if system.IsAliveContext(ctx) != nil {
			return
		}

		s.handleTitleFromQueue(ctx, page)
	}
}

func (s *Service) runFull(ctx context.Context) {
	for i := 0; i < titleHandlersCount; i++ {
		go s.runQueueHandler(ctx)
	}

	system.Info(ctx, "TitleLoader запущен")
	defer system.Info(ctx, "TitleLoader остановлен")

	timer := time.NewTicker(titleInterval)

	for {
		select {
		case <-ctx.Done():
			// Дожидаемся завершения всех подпроцессов
			s.asyncPathWG.Wait()

			return

		case <-timer.C:
			if len(s.titleQueue) > 0 || s.InWorkCount() > 0 {
				continue
			}

			for _, title := range s.Storage.GetUnloadedTitles(ctx) {
				select {
				case <-ctx.Done():
					// Дожидаемся завершения всех подпроцессов
					s.asyncPathWG.Wait()

					return

				default:
				}

				s.titleQueue <- title
			}
		}
	}
}
