package fileStorage

import (
	"app/internal/service/parser"
	"app/system"
	"context"
	"fmt"
	"os"
)

func downloadTitlePage(ctx context.Context, id, page int, URL, ext string) error {
	defer system.Stopwatch(ctx, "DownloadTitlePage")()

	// создаем папку с тайтлом
	err := os.MkdirAll(fmt.Sprintf("%s/%d", system.GetFileStoragePath(ctx), id), os.ModeDir|os.ModePerm)
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
