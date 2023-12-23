package hgraber

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func (uc *UseCase) ExportBooksToZip(ctx context.Context, from, to int) error {
	for i := from; i <= to; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := uc.saveToZip(ctx, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) saveToZip(ctx context.Context, id int) error {
	titleInfo, err := uc.storage.GetBook(ctx, id)
	if err != nil {
		return err
	}

	zipFile, err := uc.files.CreateExportFile(ctx, fmt.Sprintf(
		"%d)_%s.zip", id,
		escapeFileName(titleInfo.Data.Name),
	))
	if err != nil {
		uc.logger.Error(ctx, err)
		return err
	}

	defer uc.logger.IfErrFunc(ctx, zipFile.Close)

	zipWriter := zip.NewWriter(zipFile)

	for pageNumber, p := range titleInfo.Pages {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		pageReader, err := uc.files.OpenPageFile(ctx, id, pageNumber+1, p.Ext)
		if err != nil {
			uc.logger.Error(ctx, err)
			return err
		}
		defer uc.logger.IfErrFunc(ctx, pageReader.Close)

		w, err := zipWriter.Create(fmt.Sprintf("%d.%s", pageNumber+1, p.Ext))
		if err != nil {
			uc.logger.Error(ctx, err)
			return err
		}

		_, err = io.Copy(w, pageReader)
		if err != nil {
			uc.logger.Error(ctx, err)
			return err
		}
	}

	w, err := zipWriter.Create("info.txt")
	if err != nil {
		uc.logger.Error(ctx, err)
		return err
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
		uc.logger.Error(ctx, err)
		return err
	}

	w, err = zipWriter.Create("data.json")
	if err != nil {
		uc.logger.Error(ctx, err)
		return err
	}

	err = json.NewEncoder(w).Encode(TitleFromStorageWrap(titleInfo))
	if err != nil {
		uc.logger.Error(ctx, err)
		return err
	}

	err = zipWriter.Close()
	if err != nil {
		uc.logger.Error(ctx, err)
		return err
	}

	return nil
}
