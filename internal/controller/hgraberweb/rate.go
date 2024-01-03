package hgraberweb

import (
	"net/http"
)

func (ws *WebServer) rateUpdate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID   int `json:"id"`
			Page int `json:"page,omitempty"`
			Rate int `json:"rate"`
		}{}

		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, &request)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)

			return
		}

		if request.Page == 0 {
			err = ws.useCases.UpdateBookRate(ctx, request.ID, request.Rate)
		} else {
			err = ws.useCases.UpdatePageRate(ctx, request.ID, request.Page, request.Rate)
		}

		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)

			return
		}

		ws.webtool.WriteNoContent(ctx, w)
	})
}
