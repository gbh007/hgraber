package system

import (
	"context"
	"sync"
)

var mainWG = &sync.WaitGroup{}

func AddWaiting(ctx context.Context) {
	Debug(ctx, "Добавление ожидания")
	mainWG.Add(1)
}

func DoneWaiting(ctx context.Context) {
	Debug(ctx, "Завершение ожидания")
	mainWG.Done()
}

func WaitingChan(ctx context.Context) <-chan struct{} {
	Debug(ctx, "Получение канала ожидания")
	c := make(chan struct{})
	go func() {
		mainWG.Wait()
		close(c)
	}()
	return c
}
