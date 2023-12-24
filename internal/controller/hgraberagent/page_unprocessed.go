package hgraberagent

import (
	"app/internal/domain/agent"
	"net/http"
)

func (ws *Controller) pageUnprocessed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := new(agent.UnprocessedRequest)
		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, request)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		toHandle, err := ws.useCases.UnprocessedPages(ctx, request.Prefixes, request.Limit)
		if err != nil {
			ws.webtool.WritePlain(ctx, w, http.StatusInternalServerError, err.Error())

			return
		}

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, agent.UnprocessedResponse[agent.PageToHandle]{
			ToHandle: toHandle,
		})
	})
}
