package hgraber

import (
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
	titleInfo, err := uc.storage.GetBook(ctx, id)
	if err != nil {
		return fmt.Errorf("export book: %w", err)
	}

	// FIXME: заменить на работу с временным файлом
	zipFile := new(bytes.Buffer)

	zipWriter := zip.NewWriter(zipFile)

	for _, p := range titleInfo.Pages {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Пропускаем не скачанные книги
		if !p.Success {
			continue
		}

		pageReader, err := uc.files.OpenPageFile(ctx, id, p.PageNumber, p.Ext)
		if err != nil {
			return fmt.Errorf("export book: %w", err)
		}

		defer func() {
			closeErr := pageReader.Close()
			if closeErr != nil {
				uc.logger.ErrorContext(ctx, closeErr.Error())
			}
		}()

		w, err := zipWriter.Create(fmt.Sprintf("%d.%s", p.PageNumber, p.Ext))
		if err != nil {
			return fmt.Errorf("export book: %w", err)
		}

		_, err = io.Copy(w, pageReader)
		if err != nil {
			return fmt.Errorf("export book: %w", err)
		}
	}

	w, err := zipWriter.Create("info.txt")
	if err != nil {
		return fmt.Errorf("export book: %w", err)
	}

	_, err = fmt.Fprintf(
		w,
		"URL:%s\nNAME:%s\nPAGE-COUNT:%d\nINNER-ID:%d",
		titleInfo.URL,
		titleInfo.Data.Name,
		len(titleInfo.Pages),
		titleInfo.ID,
	)
	if err != nil {
		return fmt.Errorf("export book: %w", err)
	}

	w, err = zipWriter.Create("data.json")
	if err != nil {
		return fmt.Errorf("export book: %w", err)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	err = enc.Encode(TitleFromStorageWrap(titleInfo))
	if err != nil {
		return fmt.Errorf("export book: %w", err)
	}

	err = zipWriter.Close()
	if err != nil {
		return fmt.Errorf("export book: %w", err)
	}

	err = uc.files.CreateExportFile(
		ctx,
		fmt.Sprintf("%d)_%s.zip", id, escapeFileName(titleInfo.Data.Name)),
		zipFile,
	)
	if err != nil {
		return fmt.Errorf("export book: %w", err)
	}

	return nil
}
