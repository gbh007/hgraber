package bookHandler

import (
	"app/internal/domain"
	"app/internal/service/bookHandler/internal/parser"
	"context"
	"strings"
)

func (s *Service) updateForWorker(parentCtx context.Context, title domain.Book) {
	ctx := context.WithoutCancel(parentCtx)

	s.update(ctx, title)
}

// update обрабатывает данные тайтла (только недостающие)
func (s *Service) update(ctx context.Context, title domain.Book) {
	s.logger.Info(ctx, "начата обработка", title.ID, title.URL)
	defer s.logger.Info(ctx, "завершена обработка", title.ID, title.URL)

	p, ok, err := parser.Load(ctx, s.requester, strings.TrimSpace(title.URL))
	if err != nil {
		s.logger.Error(ctx, err)

		return
	}
	if !ok {
		return
	}

	if !title.Data.Parsed.Name {
		err = s.storage.UpdateBookName(ctx, title.ID, p.ParseName(ctx))
		if err != nil {
			s.logger.Error(ctx, err)

			return
		}
		s.logger.Info(ctx, "обновлено название", title.ID, title.URL)
	}

	for _, attr := range domain.AllAttributes {
		if !title.Data.Parsed.Attributes[attr] {
			err = s.storage.UpdateAttributes(ctx, title.ID, attr, parser.ParseAttr(ctx, p, attr))
			if err != nil {
				s.logger.Error(ctx, err)

				return
			}
			s.logger.Info(ctx, "обновлен аттрибут "+string(attr), title.ID, title.URL)
		}
	}

	if !title.Data.Parsed.Page {
		pages := p.ParsePages(ctx)
		if len(pages) > 0 {
			pagesDB := make([]domain.Page, len(pages))

			for i, page := range pages {
				pagesDB[i] = domain.Page{
					BookID:     title.ID,
					PageNumber: page.Number,
					URL:        page.URL,
					Ext:        page.Ext,
				}
			}

			err = s.storage.UpdateBookPages(ctx, title.ID, pagesDB)
			if err != nil {
				s.logger.Error(ctx, err)

				return
			}

			s.logger.Info(ctx, "обновлены страницы", title.ID, title.URL)
		}
	}
}
