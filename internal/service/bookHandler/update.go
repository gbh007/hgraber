package bookHandler

import (
	"app/internal/domain"
	"app/internal/service/bookHandler/internal/parser"
	"app/system"
	"context"
	"strings"
)

func (s *Service) updateForWorker(parentCtx context.Context, title domain.Book) {
	ctx := context.WithoutCancel(parentCtx)

	s.update(ctx, title)
}

// update обрабатывает данные тайтла (только недостающие)
func (s *Service) update(ctx context.Context, title domain.Book) {
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
		err = s.storage.UpdateBookName(ctx, title.ID, p.ParseName(ctx))
		if err != nil {
			system.Error(ctx, err)

			return
		}
		system.Info(ctx, "обновлено название", title.ID, title.URL)
	}

	for _, attr := range domain.AllAttributes {
		if !title.Data.Parsed.Attributes[attr] {
			err = s.storage.UpdateAttributes(ctx, title.ID, attr, parser.ParseAttr(ctx, p, attr))
			if err != nil {
				system.Error(ctx, err)

				return
			}
			system.Info(ctx, "обновлен аттрибут "+string(attr), title.ID, title.URL)
		}
	}

	if !title.Data.Parsed.Page {
		pages := p.ParsePages(ctx)
		if len(pages) > 0 {
			pagesDB := make([]domain.Page, len(pages))

			for i, page := range pages {
				pagesDB[i] = domain.Page{
					URL: page.URL,
					Ext: page.Ext,
				}
			}

			err = s.storage.UpdateBookPages(ctx, title.ID, pagesDB)
			if err != nil {
				system.Error(ctx, err)

				return
			}

			system.Info(ctx, "обновлены страницы", title.ID, title.URL)
		}
	}
}
