package filesystem

import (
	"app/system"
	"context"
	"io"
	"os"
	"path"
)

func (s *Storage) CreateExportFile(ctx context.Context, name string) (io.WriteCloser, error) {
	defer system.Stopwatch(ctx, "CreateExportFile")()

	f, err := os.Create(path.Join(s.exportPath, name))
	if err != nil {
		system.Error(ctx, err)

		return nil, err
	}

	return f, nil
}
