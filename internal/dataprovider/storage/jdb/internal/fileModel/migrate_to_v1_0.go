package fileModel

import "app/internal/dataprovider/storage/jdb/internal/modelV2"

func (dbd *DatabaseData) migrateToV1_0() error {
	dbd.Version = version1_0

	if len(dbd.DataV0) == 0 {
		return nil
	}

	if len(dbd.Data.Books) == 0 {
		dbd.Data.Books = make(map[int]modelV2.RawBook, len(dbd.DataV0))
	}

	for id, oldBook := range dbd.DataV0 {
		dbd.Data.Books[id] = modelV2.RawBookFromDomain(oldBook.Super())
	}

	dbd.DataV0 = nil

	return nil
}
