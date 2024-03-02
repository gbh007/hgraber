package externalfile

import (
	"app/internal/domain/externalfile"
	"net/http"
)

func (c *Controller) fileExport() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		filename := r.Header.Get(externalfile.HeaderFilename)

		err := c.fileStorage.CreateExportFile(ctx, filename, r.Body)
		if err != nil {
			c.logger.ErrorContext(ctx, err.Error())

			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		c.webtool.WriteNoContent(ctx, w)
	})
}
