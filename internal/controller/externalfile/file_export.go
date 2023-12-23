package externalfile

import (
	"app/internal/dataprovider/fileStorage/externalfile/dto"
	"net/http"
)

func (c *Controller) fileExport() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		filename := r.Header.Get(dto.ExternalFileFilename)

		err := c.fileStorage.CreateExportFile(ctx, filename, r.Body)
		if err != nil {
			c.logger.Error(ctx, err)

			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		c.webtool.WriteNoContent(ctx, w)
	})
}
