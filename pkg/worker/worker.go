package worker

import (
	"app/system"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Worker[T any] struct {
	titleQueue         chan T
	inWorkRunnersCount *atomic.Int32

	interval time.Duration

	handler func(context.Context, T)
	getter  func(context.Context) []T
}

func New[T any](
	queueSize int,
	interval time.Duration,
	handler func(context.Context, T),
	getter func(context.Context) []T,
) *Worker[T] {
	return &Worker[T]{
		titleQueue:         make(chan T, queueSize),
		inWorkRunnersCount: new(atomic.Int32),
		interval:           interval,
		handler:            handler,
		getter:             getter,
	}
}

func (w *Worker[T]) InQueueCount() int {
	return len(w.titleQueue)
}

func (w *Worker[T]) InWorkCount() int {
	return int(w.inWorkRunnersCount.Load())
}

func (w *Worker[T]) handleOne(ctx context.Context, value T) {
	defer func() {
		if p := recover(); p != nil {
			system.Info(ctx, fmt.Sprintf("panic detected %v", p))
		}
	}()

	w.inWorkRunnersCount.Add(1)
	defer w.inWorkRunnersCount.Add(-1)

	w.handler(ctx, value)
}

func (w *Worker[T]) runQueueHandler(ctx context.Context) {
	defer system.Debug(ctx, "handler остановлен")

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

	for i := 0; i < handlersCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			w.runQueueHandler(ctx)
		}()
	}

	system.Info(ctx, "запущен")
	defer system.Info(ctx, "остановлен")

	timer := time.NewTicker(w.interval)

handler:
	for {
		select {
		case <-ctx.Done():

			break handler

		case <-timer.C:
			if len(w.titleQueue) > 0 || w.InWorkCount() > 0 {
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
