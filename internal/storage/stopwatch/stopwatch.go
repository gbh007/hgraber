package stopwatch

import (
	"app/internal/domain"
	"app/system"
	"context"
)

const depth = 0

type storage interface {
	GetPage(ctx context.Context, id int, page int) (*domain.PageFullInfo, error)
	GetTitle(ctx context.Context, id int) (domain.Title, error)
	GetTitles(ctx context.Context, offset int, limit int) []domain.Title
	GetUnloadedTitles(ctx context.Context) []domain.Title
	GetUnsuccessedPages(ctx context.Context) []domain.PageFullInfo
	NewTitle(ctx context.Context, name string, URL string, loaded bool) (int, error)
	PagesCount(ctx context.Context) int
	TitlesCount(ctx context.Context) int
	UnloadedPagesCount(ctx context.Context) int
	UnloadedTitlesCount(ctx context.Context) int

	UpdatePageRate(ctx context.Context, id int, page int, rate int) error
	UpdateTitleRate(ctx context.Context, id int, rate int) error

	UpdateTitleName(ctx context.Context, id int, name string) error

	UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error
	UpdateTitlePages(ctx context.Context, id int, pages []domain.Page) error

	UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error
}

type Stopwatch struct {
	storage storage
}

func WithStopwatch(storage storage) *Stopwatch {
	return &Stopwatch{storage: storage}
}

func (s *Stopwatch) UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateAttributes ("+string(attr)+")", depth)()

	return s.storage.UpdateAttributes(ctx, id, attr, data)
}

func (s *Stopwatch) GetPage(ctx context.Context, id int, page int) (*domain.PageFullInfo, error) {
	defer system.StopwatchWithDepth(ctx, "DB - GetPage", depth)()

	return s.storage.GetPage(ctx, id, page)
}

func (s *Stopwatch) GetTitle(ctx context.Context, id int) (domain.Title, error) {
	defer system.StopwatchWithDepth(ctx, "DB - GetTitle", depth)()

	return s.storage.GetTitle(ctx, id)
}

func (s *Stopwatch) GetTitles(ctx context.Context, offset int, limit int) []domain.Title {
	defer system.StopwatchWithDepth(ctx, "DB - GetTitles", depth)()

	return s.storage.GetTitles(ctx, offset, limit)
}

func (s *Stopwatch) GetUnloadedTitles(ctx context.Context) []domain.Title {
	defer system.StopwatchWithDepth(ctx, "DB - GetUnloadedTitles", depth)()

	return s.storage.GetUnloadedTitles(ctx)
}

func (s *Stopwatch) GetUnsuccessedPages(ctx context.Context) []domain.PageFullInfo {
	defer system.StopwatchWithDepth(ctx, "DB - GetUnsuccessedPages", depth)()

	return s.storage.GetUnsuccessedPages(ctx)
}

func (s *Stopwatch) NewTitle(ctx context.Context, name string, URL string, loaded bool) (int, error) {
	defer system.StopwatchWithDepth(ctx, "DB - NewTitle", depth)()

	return s.storage.NewTitle(ctx, name, URL, loaded)
}

func (s *Stopwatch) PagesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - PagesCount", depth)()

	return s.storage.PagesCount(ctx)
}

func (s *Stopwatch) TitlesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - TitlesCount", depth)()

	return s.storage.TitlesCount(ctx)
}

func (s *Stopwatch) UnloadedPagesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - UnloadedPagesCount", depth)()

	return s.storage.UnloadedPagesCount(ctx)
}

func (s *Stopwatch) UnloadedTitlesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - UnloadedTitlesCount", depth)()

	return s.storage.UnloadedTitlesCount(ctx)
}

func (s *Stopwatch) UpdatePageRate(ctx context.Context, id int, page int, rate int) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdatePageRate", depth)()

	return s.storage.UpdatePageRate(ctx, id, page, rate)
}

func (s *Stopwatch) UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdatePageSuccess", depth)()

	return s.storage.UpdatePageSuccess(ctx, id, page, success)
}
func (s *Stopwatch) UpdateTitleName(ctx context.Context, id int, name string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleName", depth)()

	return s.storage.UpdateTitleName(ctx, id, name)
}

func (s *Stopwatch) UpdateTitlePages(ctx context.Context, id int, pages []domain.Page) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitlePages", depth)()

	return s.storage.UpdateTitlePages(ctx, id, pages)
}

func (s *Stopwatch) UpdateTitleRate(ctx context.Context, id int, rate int) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleRate", depth)()

	return s.storage.UpdateTitleRate(ctx, id, rate)
}
