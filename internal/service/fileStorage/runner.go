package fileStorage

import (
	"app/system"
	"context"
)

func (s *Service) Name() string {
	return "page storage"
}

func (s *Service) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ctx := system.NewSystemContext(parentCtx, "Page-loader")

		s.worker.Serve(ctx, handlersCount)
	}()

	return done, nil
}

func (s *Service) handle(ctx context.Context, page qPage) {

	err := downloadTitlePage(ctx, page.TitleID, page.PageNumber, page.URL, page.Ext)
	if err == nil {
		updateErr := s.storage.UpdatePageSuccess(ctx, page.TitleID, page.PageNumber, true)
		if updateErr != nil {
			system.Error(ctx, updateErr)
		}
	}
}

func (s *Service) getter(ctx context.Context) []qPage {
	raw := s.storage.GetUnsuccessPages(ctx)
	data := make([]qPage, 0, len(raw))

	for _, p := range raw {
		data = append(data, qPage{
			TitleID:    p.BookID,
			PageNumber: p.PageNumber,
			URL:        p.URL,
			Ext:        p.Ext,
		},
		)
	}

	return data
}
