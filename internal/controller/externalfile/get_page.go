package externalfile

import (
	"app/internal/domain/externalfile"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (c *Controller) getPage() http.Handler {
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

		rawFile, err := c.fileStorage.OpenPageFile(ctx, bookID, page, ext)
		if err != nil {
			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		defer func() {
			closeErr := rawFile.Close()
			if closeErr != nil {
				c.logger.ErrorContext(ctx, closeErr.Error())
			}
		}()

		rawData, err := io.ReadAll(rawFile)
		if err != nil {
			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		buff := bytes.NewReader(rawData)

		http.ServeContent(w, r, fmt.Sprintf("%d.%s", page, ext), time.Now(), buff)
	})
}
