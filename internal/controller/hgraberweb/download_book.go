package hgraberweb

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (ws *WebServer) downloadBook() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bookID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		rawFile, err := ws.useCases.Archive(ctx, bookID)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		// FIXME: на самом деле там и так возвращается буффер, так что надо отказаться от двойного перекладывания.
		rawData, err := io.ReadAll(rawFile)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		buff := bytes.NewReader(rawData)

		// FIXME: заменить на более экономную реализацию, без буффера
		http.ServeContent(w, r, fmt.Sprintf("%d.zip", bookID), time.Now(), buff)
	})
}
