package externalfile

import (
	"app/internal/dto"
	"app/pkg/webtool"
	"app/system"
	"io"
	"net/http"
	"strconv"
)

func (c *Controller) setPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bookID, err := strconv.Atoi(r.Header.Get(dto.ExternalFileBookID))
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		page, err := strconv.Atoi(r.Header.Get(dto.ExternalFilePageNumber))
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		ext := r.Header.Get(dto.ExternalFilePageExtension)

		pageFileToWrite, err := c.fileStorage.CreatePageFile(ctx, bookID, page, ext)
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		_, err = io.Copy(pageFileToWrite, r.Body)
		if err != nil {
			system.IfErrFunc(ctx, pageFileToWrite.Close)
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
