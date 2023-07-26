package stopwatch

import (
	"app/storage/schema"
	"app/super"
	"app/system"
	"context"
)

const depth = 0

type stopwatch struct {
	storage super.Storage
}

func WithStopwatch(storage super.Storage) super.Storage {
	return &stopwatch{storage: storage}
}

func (s *stopwatch) GetPage(ctx context.Context, id int, page int) (*schema.PageFullInfo, error) {
	defer system.StopwatchWithDepth(ctx, "DB - GetPage", depth)()

	return s.storage.GetPage(ctx, id, page)
}

func (s *stopwatch) GetTitle(ctx context.Context, id int) (schema.Title, error) {
	defer system.StopwatchWithDepth(ctx, "DB - GetTitle", depth)()

	return s.storage.GetTitle(ctx, id)
}

func (s *stopwatch) GetTitles(ctx context.Context, offset int, limit int) []schema.Title {
	defer system.StopwatchWithDepth(ctx, "DB - GetTitles", depth)()

	return s.storage.GetTitles(ctx, offset, limit)
}

func (s *stopwatch) GetUnloadedTitles(ctx context.Context) []schema.Title {
	defer system.StopwatchWithDepth(ctx, "DB - GetUnloadedTitles", depth)()

	return s.storage.GetUnloadedTitles(ctx)
}

func (s *stopwatch) GetUnsuccessedPages(ctx context.Context) []schema.PageFullInfo {
	defer system.StopwatchWithDepth(ctx, "DB - GetUnsuccessedPages", depth)()

	return s.storage.GetUnsuccessedPages(ctx)
}

func (s *stopwatch) NewTitle(ctx context.Context, name string, URL string, loaded bool) (int, error) {
	defer system.StopwatchWithDepth(ctx, "DB - NewTitle", depth)()

	return s.storage.NewTitle(ctx, name, URL, loaded)
}

func (s *stopwatch) PagesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - PagesCount", depth)()

	return s.storage.PagesCount(ctx)
}

func (s *stopwatch) TitlesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - TitlesCount", depth)()

	return s.storage.TitlesCount(ctx)
}

func (s *stopwatch) UnloadedPagesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - UnloadedPagesCount", depth)()

	return s.storage.UnloadedPagesCount(ctx)
}

func (s *stopwatch) UnloadedTitlesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - UnloadedTitlesCount", depth)()

	return s.storage.UnloadedTitlesCount(ctx)
}

func (s *stopwatch) UpdatePageRate(ctx context.Context, id int, page int, rate int) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdatePageRate", depth)()

	return s.storage.UpdatePageRate(ctx, id, page, rate)
}

func (s *stopwatch) UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdatePageSuccess", depth)()

	return s.storage.UpdatePageSuccess(ctx, id, page, success)
}

func (s *stopwatch) UpdateTitleAuthors(ctx context.Context, id int, authors []string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleAuthors", depth)()

	return s.storage.UpdateTitleAuthors(ctx, id, authors)
}

func (s *stopwatch) UpdateTitleCategories(ctx context.Context, id int, categories []string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleCategories", depth)()

	return s.storage.UpdateTitleCategories(ctx, id, categories)
}

func (s *stopwatch) UpdateTitleCharacters(ctx context.Context, id int, characters []string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleCharacters", depth)()

	return s.storage.UpdateTitleCharacters(ctx, id, characters)
}

func (s *stopwatch) UpdateTitleGroups(ctx context.Context, id int, groups []string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleGroups", depth)()

	return s.storage.UpdateTitleGroups(ctx, id, groups)
}

func (s *stopwatch) UpdateTitleLanguages(ctx context.Context, id int, languages []string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleLanguages", depth)()

	return s.storage.UpdateTitleLanguages(ctx, id, languages)
}

func (s *stopwatch) UpdateTitleName(ctx context.Context, id int, name string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleName", depth)()

	return s.storage.UpdateTitleName(ctx, id, name)
}

func (s *stopwatch) UpdateTitlePages(ctx context.Context, id int, pages []schema.Page) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitlePages", depth)()

	return s.storage.UpdateTitlePages(ctx, id, pages)
}

func (s *stopwatch) UpdateTitleParodies(ctx context.Context, id int, parodies []string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleParodies", depth)()

	return s.storage.UpdateTitleParodies(ctx, id, parodies)
}

func (s *stopwatch) UpdateTitleRate(ctx context.Context, id int, rate int) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleRate", depth)()

	return s.storage.UpdateTitleRate(ctx, id, rate)
}

func (s *stopwatch) UpdateTitleTags(ctx context.Context, id int, tags []string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateTitleTags", depth)()

	return s.storage.UpdateTitleTags(ctx, id, tags)
}
