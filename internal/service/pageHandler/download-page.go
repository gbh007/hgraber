package pageHandler

import (
	"app/pkg/request"
	"app/system"
	"context"
)

func (s *Service) downloadTitlePage(ctx context.Context, id, page int, URL, ext string) error {
	defer system.Stopwatch(ctx, "DownloadPage")()

	// скачиваем изображение
	data, err := request.RequestBytes(ctx, URL)
	if err != nil {
		return err
	}

	// создаем файл и загружаем туда изображение
	f, err := s.files.CreatePageFile(ctx, id, page, ext)
	if err != nil {
		system.Error(ctx, err)

		return err
	}

	_, err = f.Write(data)
	if err != nil {
		system.Error(ctx, err)
		f.Close()

		return err
	}

	return f.Close()
}
