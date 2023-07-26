package jdb

import (
	"app/system"
	"context"
	"time"
)

const (
	dbSaveInterval = time.Minute
)

func (db *Database) Name() string {
	return "storage - JDB"
}

func (db *Database) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	go func(parentCtx context.Context, filename string) {
		defer close(done)

		ctx := system.NewSystemContext(parentCtx, "DB-autosave")

		system.Info(ctx, "autosaveDB запущен")
		defer system.Info(ctx, "autosaveDB остановлен")

		timer := time.NewTicker(dbSaveInterval)

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				if db.Save(ctx, filename, false) == nil {
					system.Debug(ctx, "Автосохранение прошло успешно")
				}
			}
		}
	}(parentCtx, db.filename)

	return done, nil
}
