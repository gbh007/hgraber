package hgraber

import (
	"app/internal/domain/hgraber"
	"app/internal/externalModel"
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func (uc *UseCase) ExportBooksToZip(ctx context.Context, from, to int) error {
	for i := from; i <= to; i++ {
		uc.tempStorage.AddExport(ctx, i)
	}

	return nil
}

func (uc *UseCase) ExportBook(ctx context.Context, id int) error {
	book, err := uc.storage.GetBook(ctx, id)
	if err != nil {
		return fmt.Errorf("export book: get book from storage: %w", err)
	}

	body, err := uc.newArchive(ctx, book)
	if err != nil {
		return fmt.Errorf("export book: create archive: %w", err)
	}

	err = uc.files.CreateExportFile(
		ctx,
		fmt.Sprintf("%d)_%s.zip", book.ID, externalModel.EscapeFileName(book.Data.Name)),
		body,
	)
	if err != nil {
		return fmt.Errorf("export book: write archive: %w", err)
	}

	return nil
}

func (uc *UseCase) Archive(ctx context.Context, id int) (io.Reader, error) {
	book, err := uc.storage.GetBook(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("archive: get book from storage: %w", err)
	}

	body, err := uc.newArchive(ctx, book)
	if err != nil {
		return nil, fmt.Errorf("archive: create: %w", err)
	}

	return body, nil
}

func (uc *UseCase) newArchive(ctx context.Context, book hgraber.Book) (io.Reader, error) {
	// FIXME: заменить на работу с временным файлом
	zipFile := new(bytes.Buffer)

	zipWriter := zip.NewWriter(zipFile)

	for _, p := range book.Pages {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Пропускаем не скачанные книги
		if !p.Success {
			continue
		}

		pageReader, err := uc.files.OpenPageFile(ctx, book.ID, p.PageNumber, p.Ext)
		if err != nil {
			return nil, fmt.Errorf("open page file: %w", err)
		}

		defer func() {
			closeErr := pageReader.Close()
			if closeErr != nil {
				uc.logger.ErrorContext(ctx, closeErr.Error())
			}
		}()

		w, err := zipWriter.Create(fmt.Sprintf("%d.%s", p.PageNumber, p.Ext))
		if err != nil {
			return nil, fmt.Errorf("create page %d: %w", p.PageNumber, err)
		}

		_, err = io.Copy(w, pageReader)
		if err != nil {
			return nil, fmt.Errorf("write page %d: %w", p.PageNumber, err)
		}
	}

	w, err := zipWriter.Create("info.txt")
	if err != nil {
		return nil, fmt.Errorf("create text info: %w", err)
	}

	_, err = fmt.Fprintf(
		w,
		"URL:%s\nNAME:%s\nPAGE-COUNT:%d\nINNER-ID:%d",
		book.URL,
		book.Data.Name,
		len(book.Pages),
		book.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("write text info: %w", err)
	}

	w, err = zipWriter.Create("data.json")
	if err != nil {
		return nil, fmt.Errorf("create hg4 info: %w", err)
	}

	encV4 := json.NewEncoder(w)
	encV4.SetIndent("", "  ")

	err = encV4.Encode(externalModel.TitleFromStorageWrap(book))
	if err != nil {
		return nil, fmt.Errorf("encode hg4 info: %w", err)
	}

	w, err = zipWriter.Create("info.json")
	if err != nil {
		return nil, fmt.Errorf("create hg5 info: %w", err)
	}

	encV5 := json.NewEncoder(w)
	encV5.SetIndent("", "  ")

	err = encV5.Encode(externalModel.V5Convert(book))
	if err != nil {
		return nil, fmt.Errorf("encode hg5 info: %w", err)
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("write hg4 info: %w", err)
	}

	return zipFile, nil
}
