package web

import (
	"app/db"
	"app/file"
	"app/handler"
	"io"
	"log"
	"net/http"
	"strconv"
)

func applyTemplate(w io.Writer, name string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
	}
}

// GetMainPage возвращает главную страницу
func GetMainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	applyTemplate(w, "main", db.SelectTitles())
}

// NewTitle загружает новый тайтл
func NewTitle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	u := r.FormValue("url")
	err := handler.HandleFull(u)
	if err != nil {
		applyTemplate(w, "error", err.Error())
	} else {
		applyTemplate(w, "success", u+" успешно добавлен")
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
		applyTemplate(w, "error", err.Error())
		return
	}
	to, err := strconv.Atoi(toRaw)
	if err != nil {
		applyTemplate(w, "error", err.Error())
		return
	}
	for i := from; i <= to; i++ {
		err = file.LoadToZip(i)
		if err != nil {
			applyTemplate(w, "error", err.Error())
			return
		}
	}
	applyTemplate(w, "success", "тайтлы успешно загруженны на диск ZIP")
}
