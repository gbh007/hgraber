package titleHandler

import (
	"app/service/parser"
	"app/storage/schema"
	"app/system"
	"context"
	"strings"
)

// Update обрабатывает данные тайтла (только недостающие)
func (s *Service) Update(ctx context.Context, title schema.Title) {
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
		err = s.Storage.UpdateTitleName(ctx, title.ID, p.ParseName(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлено название", title.ID, title.URL)
	}

	if !title.Data.Parsed.Authors {
		err = s.Storage.UpdateTitleAuthors(ctx, title.ID, p.ParseAuthors(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены авторы", title.ID, title.URL)
	}

	if !title.Data.Parsed.Tags {
		err = s.Storage.UpdateTitleTags(ctx, title.ID, p.ParseTags(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены теги", title.ID, title.URL)
	}

	if !title.Data.Parsed.Characters {
		err = s.Storage.UpdateTitleCharacters(ctx, title.ID, p.ParseCharacters(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены персонажи", title.ID, title.URL)
	}

	if !title.Data.Parsed.Categories {
		err = s.Storage.UpdateTitleCategories(ctx, title.ID, p.ParseCategories(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены категории", title.ID, title.URL)
	}

	if !title.Data.Parsed.Groups {
		err = s.Storage.UpdateTitleGroups(ctx, title.ID, p.ParseGroups(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены группы", title.ID, title.URL)
	}

	if !title.Data.Parsed.Languages {
		err = s.Storage.UpdateTitleLanguages(ctx, title.ID, p.ParseLanguages(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены языки", title.ID, title.URL)
	}

	if !title.Data.Parsed.Parodies {
		err = s.Storage.UpdateTitleParodies(ctx, title.ID, p.ParseParodies(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлены пародии", title.ID, title.URL)
	}

	if !title.Data.Parsed.Page {
		pages := p.ParsePages(ctx)
		if len(pages) > 0 {
			pagesDB := make([]schema.Page, len(pages))

			for i, page := range pages {
				pagesDB[i] = schema.Page{
					URL: page.URL,
					Ext: page.Ext,
				}
			}

			err = s.Storage.UpdateTitlePages(ctx, title.ID, pagesDB)
			if err != nil {
				system.Error(ctx, err)

				return
			}

			system.Info(ctx, "обновлены страницы", title.ID, title.URL)
		}
	}
}
