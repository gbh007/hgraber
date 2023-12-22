package externalfile

import (
	"app/internal/dto"
	"app/pkg/webtool"
	"io"
	"net/http"
)

func (c *Controller) fileExport() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		filename := r.Header.Get(dto.ExternalFileFilename)

		pageFileToWrite, err := c.fileStorage.CreateExportFile(ctx, filename)
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		_, err = io.Copy(pageFileToWrite, r.Body)
		if err != nil {
			c.logger.IfErrFunc(ctx, pageFileToWrite.Close)
			webtool.WritePlain(ctx, w, http.StatusInternalServerError, err.Error())

			return
		}

		err = pageFileToWrite.Close()
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusInternalServerError, err.Error())

			return
		}

		webtool.WriteNoContent(ctx, w)
	})
}
