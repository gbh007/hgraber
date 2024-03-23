package fileModel

import (
	"app/internal/dataprovider/storage/jdb/internal/modelV1"
	"app/internal/dataprovider/storage/jdb/internal/modelV2"
	"fmt"
)

const (
	version1_0     = "v1.0"
	version1_1     = "v1.1"
	currentVersion = version1_1
)

type DatabaseData struct {
	Version string `json:"version,omitempty"`

	// Deprecated: устаревшая версия.
	DataV0 map[int]modelV1.RawTitle `json:"titles,omitempty"`

	Data DataV1 `json:"v1,omitempty"`
}

type DataV1 struct {
	Books map[int]modelV2.RawBook `json:"books,omitempty"`
}

func New() *DatabaseData {
	return &DatabaseData{
		Data: DataV1{
			Books: make(map[int]modelV2.RawBook),
		},
	}
}

func (dbd *DatabaseData) Migrate() (bool, error) {
	if dbd.Version == currentVersion {
		return false, nil
	}

	// Конвертация легаси версии
	if dbd.Version == "" && len(dbd.DataV0) > 0 {
		err := dbd.migrateToV1_0()
		if err != nil {
			return false, fmt.Errorf("migrate: %w", err)
		}
	}

	// С 1.0 до 1.1 изменений не требуется, просто появляются новые опциональные поля
	if dbd.Version == version1_0 {
		dbd.Version = version1_1
	}

	dbd.Version = currentVersion

	return true, nil
}
