package webServer

import (
	"app/pkg/webtool"
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
			webtool.WritePlain(ctx, w, http.StatusBadRequest, "not 2 pair in first split")

			return
		}

		bookID, err := strconv.Atoi(first[0])
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		second := strings.Split(first[1], ".")
		if len(second) != 2 {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, "not 2 pair in second split")

			return
		}

		page, err := strconv.Atoi(second[0])
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		info, err := ws.storage.GetPage(ctx, bookID, page)
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		rawFile, err := ws.files.OpenPageFile(ctx, bookID, page, second[1])
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		defer ws.logger.IfErrFunc(ctx, rawFile.Close)

		rawData, err := io.ReadAll(rawFile)
		if err != nil {
			webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		buff := bytes.NewReader(rawData)

		http.ServeContent(w, r, info.Fullname(), info.LoadedAt, buff)
	})
}
