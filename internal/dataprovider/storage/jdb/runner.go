package jdb

import (
	"app/pkg/ctxtool"
	"context"
	"fmt"
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

	if db.filename == nil {
		return nil, fmt.Errorf("jdb: nil filename to autosave")
	}

	go func(parentCtx context.Context, filename string) {
		defer close(done)

		ctx := ctxtool.NewSystemContext(parentCtx, "DB-autosave")

		db.logger.InfoContext(ctx, "autosaveDB запущен")
		defer db.logger.InfoContext(ctx, "autosaveDB остановлен")

		timer := time.NewTicker(dbSaveInterval)

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				if !db.needSave {
					continue
				}

				if db.Save(ctx, filename, false) == nil {
					db.logger.DebugContext(ctx, "Автосохранение прошло успешно")
				}
			}
		}
	}(parentCtx, *db.filename)

	return done, nil
}
