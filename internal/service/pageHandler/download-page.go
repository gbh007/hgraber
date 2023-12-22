package pageHandler

import (
	"context"
)

func (s *Service) downloadTitlePage(ctx context.Context, id, page int, URL, ext string) error {
	// скачиваем изображение
	data, err := s.requester.RequestBytes(ctx, URL)
	if err != nil {
		return err
	}

	// создаем файл и загружаем туда изображение
	f, err := s.files.CreatePageFile(ctx, id, page, ext)
	if err != nil {
		s.logger.Error(ctx, err)

		return err
	}

	_, err = f.Write(data)
	if err != nil {
		s.logger.Error(ctx, err)
		f.Close()

		return err
	}

	return f.Close()
}
