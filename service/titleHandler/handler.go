package titleHandler

import (
	"app/service/jdb"
	"app/service/parser"
	"app/system"
	"context"
	"strings"
)

// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
func FirstHandle(ctx context.Context, u string) error {
	system.Info(ctx, "начата обработка", u)
	defer system.Info(ctx, "завершена обработка", u)
	_, err := parser.Parse(ctx, u)
	if err != nil {
		return err
	}
	_, err = jdb.Get().NewTitle(ctx, "", u, false)
	if err != nil {
		return err
	}
	return nil
}

// Update обрабатывает данные тайтла (только недостающие)
func Update(ctx context.Context, title jdb.Title) {
	system.Info(ctx, "начата обработка", title.ID, title.URL)
	defer system.Info(ctx, "завершена обработка", title.ID, title.URL)

	p, ok, err := parser.Load(ctx, strings.TrimSpace(title.URL))
	if err != nil {
		system.Error(ctx, err)

		return
	}
	if !ok {
		return
	}

	if !title.Data.Parsed.Name {
		err = jdb.Get().UpdateTitleName(ctx, title.ID, p.ParseName(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлено название", title.ID, title.URL)
	}

	if !title.Data.Parsed.Authors {
		err = jdb.Get().UpdateTitleAuthors(ctx, title.ID, p.ParseAuthors(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены авторы", title.ID, title.URL)
	}

	if !title.Data.Parsed.Tags {
		err = jdb.Get().UpdateTitleTags(ctx, title.ID, p.ParseTags(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены теги", title.ID, title.URL)
	}

	if !title.Data.Parsed.Characters {
		err = jdb.Get().UpdateTitleCharacters(ctx, title.ID, p.ParseCharacters(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены персонажи", title.ID, title.URL)
	}

	if !title.Data.Parsed.Categories {
		err = jdb.Get().UpdateTitleCategories(ctx, title.ID, p.ParseCategories(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены категории", title.ID, title.URL)
	}

	if !title.Data.Parsed.Groups {
		err = jdb.Get().UpdateTitleGroups(ctx, title.ID, p.ParseGroups(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены группы", title.ID, title.URL)
	}

	if !title.Data.Parsed.Languages {
		err = jdb.Get().UpdateTitleLanguages(ctx, title.ID, p.ParseLanguages(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены языки", title.ID, title.URL)
	}

	if !title.Data.Parsed.Parodies {
		err = jdb.Get().UpdateTitleParodies(ctx, title.ID, p.ParseParodies(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены пародии", title.ID, title.URL)
	}

	if !title.Data.Parsed.Page {
		pages := p.ParsePages(ctx)
		if len(pages) > 0 {
			pagesDB := make([]jdb.Page, len(pages))

			for i, page := range pages {
				pagesDB[i] = jdb.Page{
					URL: page.URL,
					Ext: page.Ext,
				}
			}

			err = jdb.Get().UpdateTitlePages(ctx, title.ID, pagesDB)
			if err != nil {
				system.Error(ctx, err)

				return
			}

			system.Info(ctx, "обновлены страницы", title.ID, title.URL)
		}
	}
}
