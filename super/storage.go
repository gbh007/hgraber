package super

import (
	"app/storage/schema"
	"context"
)

type Storage interface {
	GetPage(ctx context.Context, id int, page int) (*schema.PageFullInfo, error)
	GetTitle(ctx context.Context, id int) (schema.Title, error)
	GetTitles(ctx context.Context, offset int, limit int) []schema.Title
	GetUnloadedTitles(ctx context.Context) []schema.Title
	GetUnsuccessedPages(ctx context.Context) []schema.PageFullInfo
	NewTitle(ctx context.Context, name string, URL string, loaded bool) (int, error)
	PagesCount(ctx context.Context) int
	TitlesCount(ctx context.Context) int
	UnloadedPagesCount(ctx context.Context) int
	UnloadedTitlesCount(ctx context.Context) int
	UpdatePageRate(ctx context.Context, id int, page int, rate int) error
	UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error
	UpdateTitleAuthors(ctx context.Context, id int, authors []string) error
	UpdateTitleCategories(ctx context.Context, id int, categories []string) error
	UpdateTitleCharacters(ctx context.Context, id int, characters []string) error
	UpdateTitleGroups(ctx context.Context, id int, groups []string) error
	UpdateTitleLanguages(ctx context.Context, id int, languages []string) error
	UpdateTitleName(ctx context.Context, id int, name string) error
	UpdateTitlePages(ctx context.Context, id int, pages []schema.Page) error
	UpdateTitleParodies(ctx context.Context, id int, parodies []string) error
	UpdateTitleRate(ctx context.Context, id int, rate int) error
	UpdateTitleTags(ctx context.Context, id int, tags []string) error
}
