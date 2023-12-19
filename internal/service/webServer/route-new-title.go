package webServer

import (
	"app/internal/service/webServer/rendering"
	"app/pkg/webtool"
	"net/http"
)

func (ws *WebServer) routeNewTitle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			URL  string   `json:"url,omitempty"`
			URLs []string `json:"urls,omitempty"`
		}{}

		ctx := r.Context()

		err := webtool.ParseJSON(r, &request)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		if len(request.URLs) > 0 {
			webtool.WriteJSON(
				ctx, w, http.StatusOK,
				rendering.HandleMultipleResultFromDomain(ws.title.FirstHandleMultiple(ctx, request.URLs)),
			)

			return
		}

		err = ws.title.FirstHandle(ctx, request.URL)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
		} else {
			webtool.WriteJSON(ctx, w, http.StatusOK, struct{}{})
		}
	})
}
