package agent

import (
	"app/internal/domain/agent"
	"app/internal/domain/hgraber"
	"context"
	"fmt"
	"strings"
)

func (uc *UseCase) bookHandle(ctx context.Context, book agent.BookToHandle) error {
	p, err := uc.loader.Load(ctx, strings.TrimSpace(book.URL))
	if err != nil {
		return fmt.Errorf("book handle: %w", err)
	}

	name, err := p.Name(ctx)
	if err != nil {
		return fmt.Errorf("book handle: name: %w", err)
	}

	toUpdate := agent.BookToUpdate{
		ID:         book.ID,
		Name:       name,
		Attributes: make([]agent.Attribute, 0, len(hgraber.AllAttributes)),
	}

	for _, attrCode := range hgraber.AllAttributes {
		values, err := hgraber.ParseBookAttr(ctx, p, attrCode)
		if err != nil {
			return fmt.Errorf("book handle: attributes(%s): %w", string(attrCode), err)
		}

		toUpdate.Attributes = append(toUpdate.Attributes, agent.Attribute{
			Code:   string(attrCode),
			Parsed: true,
			Values: values,
		})
	}

	pages, err := p.Pages(ctx)
	if err != nil {
		return fmt.Errorf("book handle: pages: %w", err)
	}

	if len(pages) > 0 {
		pagesToUpdate := make([]agent.PageToUpdate, len(pages))

		for i, page := range pages {
			pagesToUpdate[i] = agent.PageToUpdate{
				PageNumber: page.PageNumber,
				URL:        page.URL,
				Ext:        page.Ext,
			}
		}

		toUpdate.Pages = pagesToUpdate
	}

	err = uc.agentAPI.UpdateBook(ctx, toUpdate)
	if err != nil {
		return fmt.Errorf("book handle: %w", err)
	}

	return nil
}
