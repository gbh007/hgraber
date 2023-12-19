package webServer

import (
	"app/pkg/webtool"
	"app/system"
	"net/http"
)

func (ws *WebServer) routeAppInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Version string `json:"version"`
			Commit  string `json:"commit"`
			BuildAt string `json:"build_at"`
		}{
			Version: system.Version,
			Commit:  system.Commit,
			BuildAt: system.BuildAt,
		}

		webtool.WriteJSON(r.Context(), w, http.StatusOK, response)
	})
}
