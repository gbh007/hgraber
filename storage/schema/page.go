package schema

import "time"

type Page struct {
	URL      string
	Ext      string
	Success  bool
	LoadedAt time.Time
	Rate     int
}

type PageFullInfo struct {
	TitleID    int
	PageNumber int
	URL        string
	Ext        string
	Success    bool
	LoadedAt   time.Time
	Rate       int
}
