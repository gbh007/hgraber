package filememory

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

type Storage struct {
	data map[string][]byte

	dataMutex *sync.Mutex
}

func New() *Storage {
	return &Storage{
		data:      make(map[string][]byte),
		dataMutex: new(sync.Mutex),
	}
}

func (s *Storage) CreateExportFile(ctx context.Context, name string, body io.Reader) error {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	name = "export/" + name

	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("export file: %w", err)
	}

	s.data[name] = data

	return nil
}

func (s *Storage) CreatePageFile(ctx context.Context, id int, page int, ext string, body io.Reader) error {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	name := s.loadFilename(id, page, ext)

	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("create page file: %w", err)
	}

	s.data[name] = data

	return nil
}

func (s *Storage) OpenPageFile(ctx context.Context, id int, page int, ext string) (io.ReadCloser, error) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	name := s.loadFilename(id, page, ext)
	raw, ok := s.data[name]
	if !ok {
		return nil, fmt.Errorf("open page file: %w", os.ErrNotExist)
	}

	return io.NopCloser(bytes.NewReader(raw)), nil
}

func (s *Storage) loadFilename(id int, page int, ext string) string {
	return fmt.Sprintf("load/%d/%d.%s", id, page, ext)
}
