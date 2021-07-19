package web

import (
	"app/db"
	"app/file"
	"app/handler"
	"net/http"
	"strconv"
)

// GetMainPage возвращает главную страницу
func GetMainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	tmpl.ExecuteTemplate(w, "main", db.SelectTitles())
}

// NewTitle загружает новый тайтл
func NewTitle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	u := r.FormValue("url")
	err := handler.HandleFull(u)
	if err != nil {
		tmpl.ExecuteTemplate(w, "error", err.Error())
	} else {
		tmpl.ExecuteTemplate(w, "success", u+" успешно добавлен")
	}
}

// SaveToZIP загружает новый тайтл
func SaveToZIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	fromRaw := r.FormValue("from")
	toRaw := r.FormValue("to")
	from, err := strconv.Atoi(fromRaw)
	if err != nil {
		tmpl.ExecuteTemplate(w, "error", err.Error())
		return
	}
	to, err := strconv.Atoi(toRaw)
	if err != nil {
		tmpl.ExecuteTemplate(w, "error", err.Error())
		return
	}
	for i := from; i <= to; i++ {
		err = file.LoadToZip(i)
		if err != nil {
			tmpl.ExecuteTemplate(w, "error", err.Error())
			return
		}
	}
	tmpl.ExecuteTemplate(w, "success", "тайтлы успешно загруженны на диск ZIP")
}
