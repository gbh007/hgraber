package schema

import (
	"context"
	"time"
)

type TitleInfoParsed struct {
	Name       bool
	Page       bool
	Tags       bool
	Authors    bool
	Characters bool
	Languages  bool
	Categories bool
	Parodies   bool
	Groups     bool
}

func (tip TitleInfoParsed) IsFullParsed(ctx context.Context) bool {
	return tip.Name &&
		tip.Page &&
		tip.Tags &&
		tip.Authors &&
		tip.Characters &&
		tip.Languages &&
		tip.Categories &&
		tip.Parodies &&
		tip.Groups
}

type TitleInfo struct {
	Parsed     TitleInfoParsed
	Name       string
	Rate       int
	Tags       []string
	Authors    []string
	Characters []string
	Languages  []string
	Categories []string
	Parodies   []string
	Groups     []string
}

type Title struct {
	ID      int
	Created time.Time
	URL     string

	Pages []Page
	Data  TitleInfo
}
