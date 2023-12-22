package filesystem

import (
	"context"
	"io"
	"os"
	"path"
)

func (s *Storage) CreateExportFile(ctx context.Context, name string) (io.WriteCloser, error) {
	if s.readOnly {
		return nil, readOnlyModeError
	}

	f, err := os.Create(path.Join(s.exportPath, name))
	if err != nil {
		return nil, err
	}

	return f, nil
}
