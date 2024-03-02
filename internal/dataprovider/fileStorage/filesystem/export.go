package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
)

func (s *Storage) CreateExportFile(ctx context.Context, name string, body io.Reader) error {
	if s.readOnly {
		return readOnlyModeError
	}

	f, err := os.Create(path.Join(s.exportPath, name))
	if err != nil {
		return fmt.Errorf("export file: %w", err)
	}

	_, err = io.Copy(f, body)
	if err != nil {
		fileCloseErr := f.Close()
		if fileCloseErr != nil {
			s.logger.ErrorContext(ctx, fileCloseErr.Error())
		}

		return fmt.Errorf("export file: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("export file: %w", err)
	}

	return nil
}
