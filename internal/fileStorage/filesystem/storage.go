package filesystem

import (
	"app/system"
	"context"
	"errors"
	"fmt"
	"os"
)

type Storage struct {
	loadPath   string
	exportPath string
}

func New(load, export string) *Storage {
	return &Storage{
		loadPath:   load,
		exportPath: export,
	}
}

func (s *Storage) Prepare(ctx context.Context) error {
	err := createDir(ctx, s.loadPath)
	if err != nil {
		return err
	}

	err = createDir(ctx, s.exportPath)
	if err != nil {
		return err
	}

	return nil
}

func createDir(ctx context.Context, dirPath string) error {
	info, err := os.Stat(dirPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		system.Error(ctx, err)

		return err
	}

	if info != nil && !info.IsDir() {
		err = fmt.Errorf("dir path is not dir")
		system.Error(ctx, err)

		return err
	}

	err = os.MkdirAll(dirPath, os.ModeDir|os.ModePerm)
	if err != nil {
		system.Error(ctx, err)

		return err
	}

	return nil
}
