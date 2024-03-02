package postgresql

import (
	"app/internal/dataprovider/storage/postgresql/internal/migration"
	"context"
	"database/sql"
	"log/slog"

	migrator "gitlab.com/gbh007/go-sql-migrator"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // драйвер для PostgreSQL
)

type loggerWrapper struct { // FIXME: добавить адаптер slog в мигратор
	logger *slog.Logger
}

func (l *loggerWrapper) Info(ctx context.Context, message string) {
	l.logger.InfoContext(ctx, message)
}

func (l *loggerWrapper) Error(ctx context.Context, err error) {
	l.logger.ErrorContext(ctx, err.Error())
}

type Database struct {
	db *sqlx.DB

	logger *slog.Logger
}

// Connect - возвращает соединение с хранилищем данных
func Connect(ctx context.Context, dataSourceName string, logger *slog.Logger) (*Database, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)

	return &Database{db: db, logger: logger}, nil
}

// MigrateAll - производит миграции данных
func (storage *Database) MigrateAll(ctx context.Context) error {
	return migrator.New().
		WithFS(migration.Migrations).
		WithCtxLogger(&loggerWrapper{logger: storage.logger}).
		WithProvider(migrator.PostgreSQLProvider).
		MigrateAll(ctx, storage.db, true)
}

func isApplyWithErr(r sql.Result) (bool, error) {
	c, err := r.RowsAffected()
	if err != nil {
		return false, nil
	}

	return c != 0, nil
}

func (storage *Database) isApply(ctx context.Context, r sql.Result) bool {
	apply, err := isApplyWithErr(r)

	if err != nil {
		storage.logger.ErrorContext(ctx, err.Error())
	}

	return apply
}
