package hgraberagent

import (
	"app/internal/domain/agent"
	"net/http"
)

func (ws *Controller) createBooks() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := new(agent.CreateBooksRequest)
		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, request)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		result, err := ws.useCases.CreateMultipleBook(ctx, request.URLs)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusInternalServerError, err.Error())

			return
		}

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, result)
	})
}
