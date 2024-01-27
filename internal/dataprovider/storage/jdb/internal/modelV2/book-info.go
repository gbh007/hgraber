package modelV2

import "time"

type RawBookInfo struct {
	Created   time.Time `json:"created"`
	URL       string    `json:"url"`
	Name      string    `json:"name,omitempty"`
	Rating    int       `json:"rating,omitempty"`
	PageCount int       `json:"page_count,omitempty"`
}

func (info RawBookInfo) Copy() RawBookInfo {
	return RawBookInfo{
		Created:   info.Created,
		URL:       info.URL,
		Name:      info.Name,
		Rating:    info.Rating,
		PageCount: info.PageCount,
	}
}
