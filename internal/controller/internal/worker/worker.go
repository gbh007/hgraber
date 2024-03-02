package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type Worker[T any] struct {
	titleQueue         chan T
	inWorkRunnersCount *atomic.Int32
	runnersCount       *atomic.Int32

	interval time.Duration

	handler func(context.Context, T)
	getter  func(context.Context) []T

	logger *slog.Logger
}

func New[T any](
	queueSize int,
	interval time.Duration,
	logger *slog.Logger,
	handler func(context.Context, T),
	getter func(context.Context) []T,
) *Worker[T] {
	return &Worker[T]{
		titleQueue:         make(chan T, queueSize),
		inWorkRunnersCount: new(atomic.Int32),
		runnersCount:       new(atomic.Int32),
		interval:           interval,
		handler:            handler,
		getter:             getter,

		logger: logger,
	}
}

func (w *Worker[T]) InQueueCount() int {
	return len(w.titleQueue)
}

func (w *Worker[T]) InWorkCount() int {
	return int(w.inWorkRunnersCount.Load())
}

func (w *Worker[T]) RunnersCount() int {
	return int(w.runnersCount.Load())
}

func (w *Worker[T]) handleOne(ctx context.Context, value T) {
	defer func() {
		if p := recover(); p != nil {
			w.logger.WarnContext(ctx, fmt.Sprintf("panic detected %v", p))
		}
	}()

	w.inWorkRunnersCount.Add(1)
	defer w.inWorkRunnersCount.Add(-1)

	w.handler(ctx, value)
}

func (w *Worker[T]) runQueueHandler(ctx context.Context) {
	defer w.logger.DebugContext(ctx, "handler остановлен")

	for {
		select {
		case value := <-w.titleQueue:
			w.handleOne(ctx, value)
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker[T]) Serve(ctx context.Context, handlersCount int) {
	wg := new(sync.WaitGroup)

	w.runnersCount.Store(int32(handlersCount))

	for i := 0; i < handlersCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			w.runQueueHandler(ctx)
		}()
	}

	w.logger.InfoContext(ctx, "запущен")
	defer w.logger.InfoContext(ctx, "остановлен")

	timer := time.NewTicker(w.interval)

handler:
	for {
		select {
		case <-ctx.Done():

			break handler

		case <-timer.C:
			if len(w.titleQueue) > 0 {
				continue
			}

			for _, title := range w.getter(ctx) {
				select {
				case <-ctx.Done():
					break handler

				case w.titleQueue <- title:
				}

			}
		}
	}

	// Дожидаемся завершения всех подпроцессов
	wg.Wait()
}
