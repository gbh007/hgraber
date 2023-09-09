package domain

import (
	"fmt"
	"time"
)

type Page struct {
	URL      string
	Ext      string
	Success  bool
	LoadedAt time.Time
	Rate     int
}

type PageFullInfo struct {
	BookID     int
	PageNumber int
	URL        string
	Ext        string
	Success    bool
	LoadedAt   time.Time
	Rate       int
}

func (info PageFullInfo) Fullname() string {
	return fmt.Sprintf("%d.%s", info.PageNumber, info.Ext)
}
