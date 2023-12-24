package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
	"strings"
)

func (uc *UseCase) ParseWithUpdate(ctx context.Context, book hgraber.Book) {
	uc.logger.Info(ctx, "начата обработка", book.ID, book.URL)
	defer uc.logger.Info(ctx, "завершена обработка", book.ID, book.URL)

	p, err := uc.loader.Load(ctx, strings.TrimSpace(book.URL))
	if err != nil {
		uc.logger.Error(ctx, err)

		return
	}

	if !book.Data.Parsed.Name {
		err = uc.storage.UpdateBookName(ctx, book.ID, p.ParseName(ctx))
		if err != nil {
			uc.logger.Error(ctx, err)

			return
		}

		uc.logger.Info(ctx, "обновлено название", book.ID, book.URL)
	}

	for _, attr := range hgraber.AllAttributes {
		if book.Data.Parsed.Attributes[attr] {
			continue
		}

		err = uc.storage.UpdateAttributes(ctx, book.ID, attr, hgraber.ParseAttr(ctx, p, attr))
		if err != nil {
			uc.logger.Error(ctx, err)

			return
		}

		uc.logger.Info(ctx, "обновлен аттрибут "+string(attr), book.ID, book.URL)
	}

	if !book.Data.Parsed.Page {
		pages := p.ParsePages(ctx)
		if len(pages) > 0 {
			pagesDB := make([]hgraber.Page, len(pages))

			for i, page := range pages {
				pagesDB[i] = hgraber.Page{
					BookID:     book.ID,
					PageNumber: page.PageNumber,
					URL:        page.URL,
					Ext:        page.Ext,
				}
			}

			err = uc.storage.UpdateBookPages(ctx, book.ID, pagesDB)
			if err != nil {
				uc.logger.Error(ctx, err)

				return
			}

			uc.logger.Info(ctx, "обновлены страницы", book.ID, book.URL)
		}
	}
}
