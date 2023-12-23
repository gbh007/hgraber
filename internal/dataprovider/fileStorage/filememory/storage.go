package filememory

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

type sValue struct {
	buff      *bytes.Buffer
	buffMutex *sync.Mutex
}

func (v *sValue) Write(p []byte) (n int, err error) {
	return v.buff.Write(p)
}

func (v *sValue) Read(p []byte) (n int, err error) {
	return v.buff.Read(p)
}

func (v *sValue) Close() error {
	v.buffMutex.Unlock()

	return nil
}

type Storage struct {
	data map[string]*sValue

	dataMutex *sync.Mutex
}

func New() *Storage {
	return &Storage{
		data:      make(map[string]*sValue),
		dataMutex: new(sync.Mutex),
	}
}

func (s *Storage) CreateExportFile(ctx context.Context, name string) (io.WriteCloser, error) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	name = "export/" + name

	return s.file(name)
}

func (s *Storage) CreatePageFile(ctx context.Context, id int, page int, ext string) (io.WriteCloser, error) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	return s.file(s.loadFilename(id, page, ext))
}

func (s *Storage) OpenPageFile(ctx context.Context, id int, page int, ext string) (io.ReadCloser, error) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	return s.file(s.loadFilename(id, page, ext))
}

func (s *Storage) loadFilename(id int, page int, ext string) string {
	return fmt.Sprintf("load/%d/%d.%s", id, page, ext)
}

func (s *Storage) file(name string) (*sValue, error) {
	r, found := s.data[name]
	if !found {
		v := &sValue{
			buff:      new(bytes.Buffer),
			buffMutex: new(sync.Mutex),
		}

		v.buffMutex.Lock()

		s.data[name] = v

		return v, nil
	}

	locked := r.buffMutex.TryLock()
	if !locked {
		return nil, os.ErrPermission
	}

	return r, nil
}
