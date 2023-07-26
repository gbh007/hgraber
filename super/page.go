package super

import "context"

type PageHandler interface {
	ExportTitlesToZip(ctx context.Context, from, to int) error
}
