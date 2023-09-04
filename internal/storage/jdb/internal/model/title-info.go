package model

import (
	"app/internal/domain"
	"context"
)

type RawTitleInfoParsed struct {
	Name       bool `json:"name,omitempty"`
	Page       bool `json:"page,omitempty"`
	Tags       bool `json:"tags,omitempty"`
	Authors    bool `json:"authors,omitempty"`
	Characters bool `json:"characters,omitempty"`
	Languages  bool `json:"languages,omitempty"`
	Categories bool `json:"categories,omitempty"`
	Parodies   bool `json:"parodies,omitempty"`
	Groups     bool `json:"groups,omitempty"`
}

func (tip RawTitleInfoParsed) Copy(ctx context.Context) RawTitleInfoParsed {
	return RawTitleInfoParsed{
		Name:       tip.Name,
		Page:       tip.Page,
		Tags:       tip.Tags,
		Authors:    tip.Authors,
		Characters: tip.Characters,
		Languages:  tip.Languages,
		Categories: tip.Categories,
		Parodies:   tip.Parodies,
		Groups:     tip.Groups,
	}
}

func (tip RawTitleInfoParsed) Super(ctx context.Context) domain.TitleInfoParsed {
	t := domain.TitleInfoParsed{
		Name:       tip.Name,
		Page:       tip.Page,
		Attributes: make(map[domain.Attribute]bool, len(domain.AllAttributes)),
	}

	t.Attributes[domain.AttrTag] = tip.Tags
	t.Attributes[domain.AttrAuthor] = tip.Authors
	t.Attributes[domain.AttrCharacter] = tip.Characters
	t.Attributes[domain.AttrLanguage] = tip.Languages
	t.Attributes[domain.AttrCategory] = tip.Categories
	t.Attributes[domain.AttrParody] = tip.Parodies
	t.Attributes[domain.AttrGroup] = tip.Groups

	return t
}

func (tip RawTitleInfoParsed) IsFullParsed(ctx context.Context) bool {
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

type RawTitleInfo struct {
	Parsed     RawTitleInfoParsed `json:"parsed,omitempty"`
	Name       string             `json:"name,omitempty"`
	Rate       int                `json:"rate,omitempty"`
	Tags       []string           `json:"tags,omitempty"`
	Authors    []string           `json:"authors,omitempty"`
	Characters []string           `json:"characters,omitempty"`
	Languages  []string           `json:"languages,omitempty"`
	Categories []string           `json:"categories,omitempty"`
	Parodies   []string           `json:"parodies,omitempty"`
	Groups     []string           `json:"groups,omitempty"`
}

func (ti RawTitleInfo) Copy(ctx context.Context) RawTitleInfo {
	c := RawTitleInfo{
		Parsed:     ti.Parsed.Copy(ctx),
		Name:       ti.Name,
		Rate:       ti.Rate,
		Tags:       make([]string, len(ti.Tags)),
		Authors:    make([]string, len(ti.Authors)),
		Characters: make([]string, len(ti.Characters)),
		Languages:  make([]string, len(ti.Languages)),
		Categories: make([]string, len(ti.Categories)),
		Parodies:   make([]string, len(ti.Parodies)),
		Groups:     make([]string, len(ti.Groups)),
	}

	copy(c.Tags, ti.Tags)
	copy(c.Authors, ti.Authors)
	copy(c.Characters, ti.Characters)
	copy(c.Languages, ti.Languages)
	copy(c.Categories, ti.Categories)
	copy(c.Parodies, ti.Parodies)
	copy(c.Groups, ti.Groups)

	return c
}

func (ti RawTitleInfo) Super(ctx context.Context) domain.TitleInfo {
	c := domain.TitleInfo{
		Parsed:     ti.Parsed.Super(ctx),
		Name:       ti.Name,
		Rate:       ti.Rate,
		Attributes: make(map[domain.Attribute][]string, len(domain.AllAttributes)),
	}

	c.Attributes[domain.AttrTag] = make([]string, len(ti.Tags))
	c.Attributes[domain.AttrAuthor] = make([]string, len(ti.Authors))
	c.Attributes[domain.AttrCharacter] = make([]string, len(ti.Characters))
	c.Attributes[domain.AttrLanguage] = make([]string, len(ti.Languages))
	c.Attributes[domain.AttrCategory] = make([]string, len(ti.Categories))
	c.Attributes[domain.AttrParody] = make([]string, len(ti.Parodies))
	c.Attributes[domain.AttrGroup] = make([]string, len(ti.Groups))

	copy(c.Attributes[domain.AttrTag], ti.Tags)
	copy(c.Attributes[domain.AttrAuthor], ti.Authors)
	copy(c.Attributes[domain.AttrCharacter], ti.Characters)
	copy(c.Attributes[domain.AttrLanguage], ti.Languages)
	copy(c.Attributes[domain.AttrCategory], ti.Categories)
	copy(c.Attributes[domain.AttrParody], ti.Parodies)
	copy(c.Attributes[domain.AttrGroup], ti.Groups)

	return c
}
