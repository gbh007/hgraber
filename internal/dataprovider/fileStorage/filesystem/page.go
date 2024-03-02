package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
)

func (s *Storage) CreatePageFile(ctx context.Context, id, page int, ext string, body io.Reader) error {
	if s.readOnly {
		return readOnlyModeError
	}

	// создаем папку с тайтлом
	err := os.MkdirAll(fmt.Sprintf("%s/%d", s.loadPath, id), os.ModeDir|os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("create page file: %w", err)
	}

	// создаем файл для загрузки изображения
	f, err := os.Create(fmt.Sprintf("%s/%d/%d.%s", s.loadPath, id, page, ext))
	if err != nil {
		return fmt.Errorf("create page file: %w", err)
	}

	_, err = io.Copy(f, body)
	if err != nil {
		fileCloseErr := f.Close()
		if fileCloseErr != nil {
			s.logger.ErrorContext(ctx, fileCloseErr.Error())
		}

		return fmt.Errorf("create page file: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("create page file: %w", err)
	}

	return nil
}

func (s *Storage) OpenPageFile(ctx context.Context, id, page int, ext string) (io.ReadCloser, error) {
	f, err := os.Open(fmt.Sprintf("%s/%d/%d.%s", s.loadPath, id, page, ext))
	if err != nil {
		return nil, err
	}

	return f, nil
}
