package webServer

import (
	"app/db"
	"app/service/fileStorage"
	"app/service/titleHandler"
	"app/service/webServer/base"
	"net/http"
)

func MainInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base.SetResponse(r, map[string]interface{}{
			"count":               db.SelectTitlesCount(r.Context()),
			"not_load_count":      db.SelectUnloadTitlesCount(r.Context()),
			"page_count":          db.SelectPagesCount(r.Context()),
			"not_load_page_count": db.SelectUnloadPagesCount(r.Context()),
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
			base.SetError(r, err)
			return
		}

		err = titleHandler.FirstHandle(r.Context(), request.URL)
		if err != nil {
			base.SetError(r, err)
		} else {
			base.SetResponse(r, struct{}{})
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
			base.SetError(r, err)
			return
		}
		data := db.SelectTitles(r.Context(), request.Offset, request.Count)
		base.SetResponse(r, data)
	})
}

func TitleInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID int `json:"id"`
		}{}
		err := base.ParseJSON(r, &request)
		if err != nil {
			base.SetError(r, err)
			return
		}
		data, err := db.SelectTitleByID(r.Context(), request.ID)
		if err != nil {
			base.SetError(r, err)
			return
		}
		base.SetResponse(r, data)
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
			base.SetError(r, err)
			return
		}
		data, err := db.SelectPagesByTitleIDAndNumber(r.Context(), request.ID, request.Page)
		if err != nil {
			base.SetError(r, err)
			return
		}
		base.SetResponse(r, data)
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
			base.SetError(r, err)
			return
		}
		for i := request.From; i <= request.To; i++ {
			err = fileStorage.SaveToZip(r.Context(), i)
			if err != nil {
				base.SetError(r, err)
				return
			}
		}
		base.SetResponse(r, struct{}{})
	})
}
