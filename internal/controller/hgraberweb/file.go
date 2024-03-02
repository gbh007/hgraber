package hgraberweb

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (ws *WebServer) getFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// /file/100/5.jpg
		ctx := r.Context()

		first := strings.Split(r.URL.Path, "/")
		if len(first) != 2 {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, "not 2 pair in first split")

			return
		}

		bookID, err := strconv.Atoi(first[0])
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		second := strings.Split(first[1], ".")
		if len(second) != 2 {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, "not 2 pair in second split")

			return
		}

		page, err := strconv.Atoi(second[0])
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		info, rawFile, err := ws.useCases.PageWithBody(ctx, bookID, page)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		defer func() {
			closeErr := rawFile.Close()
			if closeErr != nil {
				ws.logger.ErrorContext(ctx, closeErr.Error())
			}
		}()

		rawData, err := io.ReadAll(rawFile)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		buff := bytes.NewReader(rawData)

		// FIXME: заменить на более экономную реализацию, без буффера
		http.ServeContent(w, r, info.Fullname(), info.LoadedAt, buff)
	})
}
