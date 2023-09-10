package pageHandler

import (
	"app/internal/service/webServer/rendering"
	"app/system"
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (s *Service) ExportBooksToZip(ctx context.Context, from, to int) error {
	for i := from; i <= to; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := s.saveToZip(ctx, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// saveToZip сохраняет тайтлы на диск zip архивом
func (s *Service) saveToZip(ctx context.Context, id int) error {
	defer system.Stopwatch(ctx, "saveToZip")()

	titleInfo, err := s.storage.GetBook(ctx, id)
	if err != nil {
		return err
	}

	zipFile, err := s.files.CreateExportFile(ctx, fmt.Sprintf(
		"%d)_%s.zip", id,
		escapeFileName(titleInfo.Data.Name),
	))
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	defer system.IfErrFunc(ctx, zipFile.Close)

	zipWriter := zip.NewWriter(zipFile)

	for pageNumber, p := range titleInfo.Pages {
		pageReader, err := s.files.OpenPageFile(ctx, id, pageNumber+1, p.Ext)
		if err != nil {
			system.Error(ctx, err)
			return err
		}
		defer system.IfErrFunc(ctx, pageReader.Close)

		w, err := zipWriter.Create(fmt.Sprintf("%d.%s", pageNumber+1, p.Ext))
		if err != nil {
			system.Error(ctx, err)
			return err
		}

		_, err = io.Copy(w, pageReader)
		if err != nil {
			system.Error(ctx, err)
			return err
		}
	}

	w, err := zipWriter.Create("info.txt")
	if err != nil {
		system.Error(ctx, err)
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
		system.Error(ctx, err)
		return err
	}

	w, err = zipWriter.Create("data.json")
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	// TODO: на данный момент используется формат веб сервера, это очень плохо, необходимо продублировать и разделить
	err = json.NewEncoder(w).Encode(rendering.TitleFromStorage(titleInfo))
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	err = zipWriter.Close()
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	return nil
}

func escapeFileName(n string) string {
	const replacer = " "

	if len([]rune(n)) > 200 {
		n = string([]rune(n)[:200])
	}

	for _, e := range []string{`\`, `/`, `|`, `:`, `"`, `*`, `?`, `<`, `>`} {
		n = strings.ReplaceAll(n, e, replacer)
	}

	return n
}