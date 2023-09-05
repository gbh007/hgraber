package migrator

import (
	"context"
	"io/fs"

	"github.com/jmoiron/sqlx"
)

// MigrateAll - производит накат всех доступных миграций
func MigrateAll( //nolint:cyclop // требуется рефакторинг в будущем
	ctx context.Context, migrationsDir fs.FS, db *sqlx.DB, checkHash bool, provider Provider,
) error {
	err := New().
		WithFS(migrationsDir).
		WithProvider(provider).
		MigrateAll(ctx, db, checkHash)
	if err != nil {
		return err
	}

	return nil
}
