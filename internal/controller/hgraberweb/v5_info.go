package hgraberweb

import (
	"app/internal/domain/hgraber"
	"app/internal/externalModel"
	"errors"
	"net/http"
	"strconv"
)

func (ws *WebServer) v5Info() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bookID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)

			return
		}

		book, err := ws.useCases.GetBook(ctx, bookID)
		if errors.Is(err, hgraber.BookNotFoundError) {
			ws.webtool.WriteJSON(ctx, w, http.StatusNotFound, err)

			return
		}

		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)

			return
		}

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, externalModel.V5Convert(book))
	})
}
