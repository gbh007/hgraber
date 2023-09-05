package migrator

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MigrateBuilder struct {
	log           Logger
	provider      Provider
	migrationsDir fs.FS
}

func New() *MigrateBuilder {
	return &MigrateBuilder{
		log: &simpleLogger{},
	}
}

func (b *MigrateBuilder) WithLogger(l Logger) *MigrateBuilder {
	b.log = l

	return b
}

func (b *MigrateBuilder) WithProvider(p Provider) *MigrateBuilder {
	b.provider = p

	return b
}

func (b *MigrateBuilder) WithFS(d fs.FS) *MigrateBuilder {
	b.migrationsDir = d

	return b
}

// MigrateAll - производит накат всех доступных миграций
func (b *MigrateBuilder) MigrateAll(
	ctx context.Context, db *sqlx.DB, checkHash bool,
) error {
	if b.log == nil {
		return fmt.Errorf("%w: %w: log", ErrMigrator, ErrInvalidBuildConfiguration)
	}

	if b.provider == nil {
		return fmt.Errorf("%w: %w: provider", ErrMigrator, ErrInvalidBuildConfiguration)
	}

	if b.migrationsDir == nil {
		return fmt.Errorf("%w: %w: FS", ErrMigrator, ErrInvalidBuildConfiguration)
	}

	list, err := b.getFileList(ctx, b.migrationsDir)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMigrator, err)
	}

	var success bool

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMigrator, err)
	}

	// Функция для финализации транзакции
	defer func() {
		if success {
			logIfErr(tx.Commit())
		} else {
			logIfErr(tx.Rollback())
		}
	}()

	err = b.provider.ApplyInnerMigration(ctx, tx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMigrator, err)
	}

	// Получаем список примененных миграций из БД
	appliedMigrationsList, err := b.provider.GetAppliedMigration(ctx, tx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMigrator, err)
	}

	appliedMigrationsMap := make(map[int]Migration)
	for _, migration := range appliedMigrationsList {
		appliedMigrationsMap[migration.ID] = migration
	}

	for _, migrationFile := range list {
		mteInfo, migrationApplied := appliedMigrationsMap[migrationFile.Number]

		var hash, body string

		// Не применена миграция или нужно сверить ее хеш
		if !migrationApplied || checkHash {
			body, hash, err = b.migrationFromFile(ctx, migrationFile, b.migrationsDir)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrMigrator, err)
			}
		}

		hashEqual := mteInfo.Hash == hash

		switch {
		// Миграция не применена
		case !migrationApplied:
			err = b.provider.ApplyMigration(ctx, tx, migrationFile.Number, body, migrationFile.Name, hash)
			if err != nil {
				b.log.Info(fmt.Sprintf("%s - ERR", migrationFile.Name))

				return fmt.Errorf("%w: %w", ErrMigrator, err)
			}

			b.log.Info(fmt.Sprintf("%s - OK", migrationFile.Name))

		// Миграция применена и хеш не эквивалентен
		case migrationApplied && checkHash && !hashEqual:
			b.log.Info("old - hash >>> " + mteInfo.Hash)
			b.log.Info("new - hash >>> " + hash)
			b.log.Info(fmt.Sprintf("%s - Not Equal HASH", migrationFile.Name))

		// Миграция уже применена
		case migrationApplied:
			b.log.Info(fmt.Sprintf("%s - EXIST %v", migrationFile.Name, mteInfo.Applied))
		}
	}

	success = true

	return nil
}

// getFileList - возвращает список файлов миграций
func (b *MigrateBuilder) getFileList(_ context.Context, migrationsDir fs.FS) ([]migrationFile, error) {
	migrationFileList, err := fs.ReadDir(migrationsDir, ".")
	if err != nil {
		return nil, err
	}

	list := []migrationFile{}

	for _, migrationFileInfo := range migrationFileList {
		filename := migrationFileInfo.Name()

		// Миграции только в sql файлах
		if !strings.HasSuffix(filename, ".sql") {
			b.log.Info("Не sql, пропускаю " + filename)

			continue
		}

		// Первые 4 символа обозначают номер миграции, если не получилось их обработать то игнорируем файл
		number, err := strconv.Atoi(filename[:4])
		if err != nil {
			b.log.Error(err)

			continue
		}

		list = append(list, migrationFile{
			Number: number,
			Path:   filename,
			Name:   filename,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Number < list[j].Number
	})

	b.log.Info("Доступные миграции")

	for _, item := range list {
		b.log.Info(fmt.Sprintf("%4d > %s", item.Number, item.Name))
	}

	return list, nil
}

// migrationFromFile - получает данные для применения миграции из файла
func (b *MigrateBuilder) migrationFromFile(_ context.Context, info migrationFile, migrationsDir fs.FS) (string, string, error) {
	file, err := migrationsDir.Open(info.Path)
	if err != nil {
		return "", "", err
	}

	defer logIfErrFunc(file.Close)

	data, err := io.ReadAll(file)
	if err != nil {
		return "", "", err
	}

	hash := fmt.Sprintf("%x", md5.Sum(data))
	body := string(data)

	return body, hash, nil
}
