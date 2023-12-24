package hgraber

import (
	"fmt"
	"time"
)

type Page struct {
	BookID     int
	PageNumber int
	URL        string
	Ext        string
	Success    bool
	LoadedAt   time.Time
	Rate       int
}

func (info Page) Fullname() string {
	return fmt.Sprintf("%d.%s", info.PageNumber, info.Ext)
}
