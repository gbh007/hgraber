package migrator

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

// Стандартные провайдеры.
var (
	PostgreSQLProvider Provider = &simpleProvider{innerMigration: techMigrationPostgreSQL}
	MySQLProvider      Provider = &simpleProvider{innerMigration: techMigrationMariaDB}
	ClickHouseProvider Provider = &simpleProvider{innerMigration: techMigrationClickHouse}
	Sqlite3Provider    Provider = &simpleProvider{innerMigration: techMigrationSqlite3}
)

type Provider interface {
	// ApplyInnerMigration - применяет внутренние миграции для работы самого мигратора.
	ApplyInnerMigration(ctx context.Context, tx *sqlx.Tx) error
	// GetAppliedMigration - возвращает примененные миграции, упорядоченные по номеру.
	GetAppliedMigration(ctx context.Context, tx *sqlx.Tx) ([]Migration, error)
	// ApplyMigration - применяет миграцию.
	ApplyMigration(ctx context.Context, tx *sqlx.Tx, id int, body, filename, hash string) error
}

type simpleProvider struct {
	innerMigration string
}

func (p *simpleProvider) ApplyInnerMigration(ctx context.Context, tx *sqlx.Tx) error {
	_, err := tx.ExecContext(ctx, p.innerMigration)
	if err != nil {
		return err
	}

	return nil
}

func (p *simpleProvider) GetAppliedMigration(ctx context.Context, tx *sqlx.Tx) ([]Migration, error) {
	migrations := make([]Migration, 0)

	err := tx.SelectContext(ctx, &migrations, `SELECT * FROM migrations ORDER BY id;`)
	if err != nil {
		return nil, err
	}

	return migrations, nil
}

func (p *simpleProvider) ApplyMigration(ctx context.Context, tx *sqlx.Tx, id int, body, filename, hash string) error {
	_, err := tx.ExecContext(ctx, body)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO migrations (id, filename, hash, applied) VALUES (?, ?, ?, ?);`,
		id, filename, hash, time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	return nil
}
