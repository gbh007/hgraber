package filesystem

import (
	"app/system"
	"context"
	"fmt"
	"io"
	"os"
)

func (s *Storage) CreatePageFile(ctx context.Context, id, page int, ext string) (io.WriteCloser, error) {
	defer system.Stopwatch(ctx, "CreatePageFile")()

	// создаем папку с тайтлом
	err := os.MkdirAll(fmt.Sprintf("%s/%d", s.loadPath, id), os.ModeDir|os.ModePerm)
	if err != nil && !os.IsExist(err) {
		system.Error(ctx, err)

		return nil, err
	}

	// создаем файл для загрузки изображения
	f, err := os.Create(fmt.Sprintf("%s/%d/%d.%s", s.loadPath, id, page, ext))
	if err != nil {
		system.Error(ctx, err)

		return nil, err
	}

	return f, nil
}

func (s *Storage) OpenPageFile(ctx context.Context, id, page int, ext string) (io.ReadCloser, error) {
	defer system.Stopwatch(ctx, "OpenPageFile")()

	f, err := os.Open(fmt.Sprintf("%s/%d/%d.%s", s.loadPath, id, page, ext))
	if err != nil {
		system.Error(ctx, err)

		return nil, err
	}

	return f, nil
}
