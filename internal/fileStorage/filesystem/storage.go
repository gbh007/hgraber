package filesystem

import (
	"context"
	"errors"
	"fmt"
	"os"
)

var readOnlyModeError = errors.New("read only mode")

type Storage struct {
	loadPath   string
	exportPath string
	readOnly   bool
}

func New(load, export string, readOnly bool) *Storage {
	return &Storage{
		loadPath:   load,
		exportPath: export,
		readOnly:   readOnly,
	}
}

func (s *Storage) Prepare(ctx context.Context) error {
	if s.readOnly {
		return nil
	}

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
		return err
	}

	if info != nil && !info.IsDir() {
		err = fmt.Errorf("dir path is not dir")

		return err
	}

	err = os.MkdirAll(dirPath, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
