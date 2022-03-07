package fileStorage

import (
	"app/service/jdb"
	"app/service/parser"
	"app/system"
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func DownloadTitlePage(ctx context.Context, id, page int, URL, ext string) error {
	defer system.Stopwatch(ctx, "DownloadTitlePage")()
	// создаем папку с тайтлом
	err := os.MkdirAll(fmt.Sprintf("%s/%d", system.GetFileStoragePath(ctx), id), 0666)
	if err != nil && !os.IsExist(err) {
		system.Error(ctx, err)
		return err
	}
	// скачиваем изображение
	data, err := parser.RequestBytes(ctx, URL)
	if err != nil {
		return err
	}
	// создаем файл и загружаем туда изображение
	f, err := os.Create(fmt.Sprintf("%s/%d/%d.%s", system.GetFileStoragePath(ctx), id, page, ext))
	if err != nil {
		system.Error(ctx, err)
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		system.Error(ctx, err)
		return err
	}
	return f.Close()
}

func ExportTitlesToZip(ctx context.Context, from, to int) error {
	for i := from; i <= to; i++ {
		err := SaveToZip(ctx, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveToZip сохраняет тайтлы на диск zip архивом
func SaveToZip(ctx context.Context, id int) error {
	defer system.Stopwatch(ctx, "SaveToZip")()

	titleInfo, err := jdb.Get().GetTitle(ctx, id)
	if err != nil {
		return err
	}

	zipFile, err := os.Create(fmt.Sprintf(
		"%s/%d)_%s.zip",
		system.GetFileExportPath(ctx),
		id,
		escapeFileName(titleInfo.Data.Name),
	))
	if err != nil {
		system.Error(ctx, err)
		return err
	}
	defer system.IfErrFunc(ctx, zipFile.Close)

	zipWriter := zip.NewWriter(zipFile)

	for pageNumber, p := range titleInfo.Pages {
		pageReader, err := os.Open(fmt.Sprintf("%s/%d/%d.%s", system.GetFileStoragePath(ctx), id, pageNumber+1, p.Ext))
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
	err = json.NewEncoder(w).Encode(titleInfo)
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
	const r = " "
	if len([]rune(n)) > 200 {
		n = string([]rune(n)[:200])
	}
	for _, e := range []string{`\`, `/`, `|`, `:`, `"`, `*`, `?`, `<`, `>`} {
		n = strings.ReplaceAll(n, e, r)
	}
	return n
}
