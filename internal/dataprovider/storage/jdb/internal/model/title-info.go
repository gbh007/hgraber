package model

import (
	"app/internal/domain/hgraber"
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

func (tip RawTitleInfoParsed) Copy() RawTitleInfoParsed {
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

func (tip RawTitleInfoParsed) Super() hgraber.BookInfoParsed {
	t := hgraber.BookInfoParsed{
		Name:       tip.Name,
		Page:       tip.Page,
		Attributes: make(map[hgraber.Attribute]bool, len(hgraber.AllAttributes)),
	}

	t.Attributes[hgraber.AttrTag] = tip.Tags
	t.Attributes[hgraber.AttrAuthor] = tip.Authors
	t.Attributes[hgraber.AttrCharacter] = tip.Characters
	t.Attributes[hgraber.AttrLanguage] = tip.Languages
	t.Attributes[hgraber.AttrCategory] = tip.Categories
	t.Attributes[hgraber.AttrParody] = tip.Parodies
	t.Attributes[hgraber.AttrGroup] = tip.Groups

	return t
}

func (tip RawTitleInfoParsed) IsFullParsed() bool {
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

func (ti RawTitleInfo) Copy() RawTitleInfo {
	c := RawTitleInfo{
		Parsed:     ti.Parsed.Copy(),
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

func (ti RawTitleInfo) Super() hgraber.BookInfo {
	c := hgraber.BookInfo{
		Parsed:     ti.Parsed.Super(),
		Name:       ti.Name,
		Rate:       ti.Rate,
		Attributes: make(map[hgraber.Attribute][]string, len(hgraber.AllAttributes)),
	}

	c.Attributes[hgraber.AttrTag] = make([]string, len(ti.Tags))
	c.Attributes[hgraber.AttrAuthor] = make([]string, len(ti.Authors))
	c.Attributes[hgraber.AttrCharacter] = make([]string, len(ti.Characters))
	c.Attributes[hgraber.AttrLanguage] = make([]string, len(ti.Languages))
	c.Attributes[hgraber.AttrCategory] = make([]string, len(ti.Categories))
	c.Attributes[hgraber.AttrParody] = make([]string, len(ti.Parodies))
	c.Attributes[hgraber.AttrGroup] = make([]string, len(ti.Groups))

	copy(c.Attributes[hgraber.AttrTag], ti.Tags)
	copy(c.Attributes[hgraber.AttrAuthor], ti.Authors)
	copy(c.Attributes[hgraber.AttrCharacter], ti.Characters)
	copy(c.Attributes[hgraber.AttrLanguage], ti.Languages)
	copy(c.Attributes[hgraber.AttrCategory], ti.Categories)
	copy(c.Attributes[hgraber.AttrParody], ti.Parodies)
	copy(c.Attributes[hgraber.AttrGroup], ti.Groups)

	return c
}
