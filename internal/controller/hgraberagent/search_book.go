package hgraberagent

import (
	"app/internal/domain/agent"
	"net/http"
)

func (ws *Controller) searchBook() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := new(agent.SearchBookIDByURLRequest)
		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, request)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		id, found, err := ws.useCases.SearchBook(ctx, request.URL)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusInternalServerError, err.Error())

			return
		}

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, agent.SearchBookIDByURLResponse{
			ID:    id,
			Found: found,
		})
	})
}
