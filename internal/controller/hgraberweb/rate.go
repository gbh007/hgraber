package hgraberweb

import (
	"net/http"
)

func (ws *WebServer) ratingUpdate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID     int `json:"id"`
			Page   int `json:"page,omitempty"`
			Rating int `json:"rating"`
		}{}

		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, &request)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)

			return
		}

		if request.Page == 0 {
			err = ws.useCases.UpdateBookRate(ctx, request.ID, request.Rating)
		} else {
			err = ws.useCases.UpdatePageRate(ctx, request.ID, request.Page, request.Rating)
		}

		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)

			return
		}

		ws.webtool.WriteNoContent(ctx, w)
	})
}
