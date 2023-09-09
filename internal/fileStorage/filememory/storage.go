package filememory

import (
	"bytes"
	"sync"
)

type sKey struct {
	book, page int
	ext        string
}

type sValue struct {
	buff      *bytes.Buffer
	buffMutex *sync.Mutex
}

// FIXME: реализовать.
type Storage struct {
	data map[sKey]*sValue

	dataMutex sync.Mutex
}
