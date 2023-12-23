package webServer

import (
	"app/internal/controller/webServer/internal/rendering"
	"net/http"
)

func (ws *WebServer) routeNewTitle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			URL  string   `json:"url,omitempty"`
			URLs []string `json:"urls,omitempty"`
		}{}

		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, &request)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)

			return
		}

		if len(request.URLs) > 0 {
			data, err := ws.useCases.FirstHandleMultiple(ctx, request.URLs)
			if err != nil {
				ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)

				return
			}

			ws.webtool.WriteJSON(
				ctx, w, http.StatusOK,
				rendering.HandleMultipleResultFromDomain(data),
			)

			return
		}

		err = ws.useCases.FirstHandle(ctx, request.URL)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
		} else {
			ws.webtool.WriteJSON(ctx, w, http.StatusOK, struct{}{})
		}
	})
}
