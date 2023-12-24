package externalfile

import (
	"app/internal/domain/externalfile"
	"net/http"
	"strconv"
)

func (c *Controller) setPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bookID, err := strconv.Atoi(r.Header.Get(externalfile.HeaderBookID))
		if err != nil {
			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		page, err := strconv.Atoi(r.Header.Get(externalfile.HeaderPageNumber))
		if err != nil {
			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		ext := r.Header.Get(externalfile.HeaderPageExtension)

		err = c.fileStorage.CreatePageFile(ctx, bookID, page, ext, r.Body)
		if err != nil {
			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		c.webtool.WriteNoContent(ctx, w)
	})
}
