package system

import (
	"context"
	"errors"
	"fmt"
	"os"
)

var (
	fileStoragePath string
	fileExportPath  string
)

func GetFileStoragePath(ctx context.Context) string {
	return fileStoragePath
}

func SetFileStoragePath(ctx context.Context, dirPath string) error {
	err := createDir(ctx, dirPath)
	if err != nil {
		Warning(ctx, dirPath, "не является директорией для FileStorage")
		return err
	}

	fileStoragePath = dirPath

	return nil
}

func GetFileExportPath(ctx context.Context) string {
	return fileExportPath
}

func SetFileExportPath(ctx context.Context, dirPath string) error {
	err := createDir(ctx, dirPath)
	if err != nil {
		Warning(ctx, dirPath, "не является директорией для FileExport")
		return err
	}

	fileExportPath = dirPath

	return nil
}

func createDir(ctx context.Context, dirPath string) error {
	info, err := os.Stat(dirPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		Error(ctx, err)
		return err
	}

	if info != nil && !info.IsDir() {
		err = fmt.Errorf("dir path is not dir")
		Error(ctx, err)
		return err
	}

	err = os.MkdirAll(dirPath, os.ModeDir|os.ModePerm)
	if err != nil {
		Error(ctx, err)
		return err
	}

	return nil
}
