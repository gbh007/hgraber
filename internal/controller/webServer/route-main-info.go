package webServer

import (
	"app/internal/controller/webServer/internal/rendering"
	"app/pkg/webtool"
	"net/http"
)

func (ws *WebServer) routeMainInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		info, err := ws.useCases.Info(ctx)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)

			return
		}

		webtool.WriteJSON(ctx, w, http.StatusOK, map[string]interface{}{
			"count":               info.BookCount,
			"not_load_count":      info.NotLoadBookCount,
			"page_count":          info.PageCount,
			"not_load_page_count": info.NotLoadPageCount,
			"monitor":             rendering.MonitorFromDomain(ws.monitor.Info()),
		})
	})
}