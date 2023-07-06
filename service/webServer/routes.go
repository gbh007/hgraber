package webServer

import (
	"app/service/fileStorage"
	"app/service/jdb"
	"app/service/titleHandler"
	"app/service/webServer/base"
	"app/system"
	"net/http"
)

func MainInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base.WriteJSON(r.Context(), w, http.StatusOK, map[string]interface{}{
			"count":               jdb.Get().TitlesCount(r.Context()),
			"not_load_count":      jdb.Get().UnloadedTitlesCount(r.Context()),
			"page_count":          jdb.Get().PagesCount(r.Context()),
			"not_load_page_count": jdb.Get().UnloadedPagesCount(r.Context()),
		})
	})
}

func NewTitle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			URL string `json:"url"`
		}{}

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		err = titleHandler.FirstHandle(r.Context(), request.URL)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusInternalServerError, err)
		} else {
			base.WriteJSON(r.Context(), w, http.StatusOK, struct{}{})
		}
	})
}

func TitleList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			Count  int `json:"count"`
			Offset int `json:"offset,omitempty"`
		}{}

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		data := jdb.Get().GetTitles(r.Context(), request.Offset, request.Count)
		base.WriteJSON(r.Context(), w, http.StatusOK, data)
	})
}

func TitleInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID int `json:"id"`
		}{}
		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		data, err := jdb.Get().GetTitle(r.Context(), request.ID)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSON(r.Context(), w, http.StatusOK, data)
	})
}

func TitlePage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID   int `json:"id"`
			Page int `json:"page"`
		}{}

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		data, err := jdb.Get().GetPage(r.Context(), request.ID, request.Page)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSON(r.Context(), w, http.StatusOK, data)
	})
}

func SaveToZIP() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			From int `json:"from"`
			To   int `json:"to"`
		}{}

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		err = fileStorage.ExportTitlesToZip(r.Context(), request.From, request.To)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSON(r.Context(), w, http.StatusOK, struct{}{})
	})
}

func AppInfo() http.Handler {
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

		base.WriteJSON(r.Context(), w, http.StatusOK, response)
	})
}

func SetTitleRate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID   int `json:"id"`
			Rate int `json:"rate"`
		}{}

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		err = jdb.Get().UpdateTitleRate(r.Context(), request.ID, request.Rate)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSON(r.Context(), w, http.StatusOK, struct{}{})
	})
}

func SetPageRate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID   int `json:"id"`
			Page int `json:"page"`
			Rate int `json:"rate"`
		}{}

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		err = jdb.Get().UpdatePageRate(r.Context(), request.ID, request.Page, request.Rate)
		if err != nil {
			base.WriteJSON(r.Context(), w, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSON(r.Context(), w, http.StatusOK, struct{}{})
	})
}
