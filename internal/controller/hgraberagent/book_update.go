package hgraberagent

import (
	"app/internal/domain/agent"
	"net/http"
)

func (ws *Controller) bookUpdate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		book := agent.BookToUpdate{}
		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, &book)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		err = ws.useCases.UpdateBook(ctx, book)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusInternalServerError, err.Error())

			return
		}

		ws.webtool.WriteNoContent(ctx, w)
	})
}
