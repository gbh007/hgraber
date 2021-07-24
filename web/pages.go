package web

import (
	"app/db"
	"app/file"
	"app/handler"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	err := handler.FirstHandle(u)
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

// GetTitlePage возвращает страницу из тайтла
func GetTitlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	tmp := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(tmp) != 3 || tmp[0] != "title" {
		applyTemplate(w, "error", "ошибка адресации")
		return
	}
	tid, err := strconv.Atoi(tmp[1])
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	pid, err := strconv.Atoi(tmp[2])
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	title, err := db.SelectTitleByID(tid)
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	page, err := db.SelectPagesByTitleIDAndNumber(tid, pid)
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	data := struct {
		TitleID, PageNumber int
		Title               db.TitleShortInfo
		Page                db.Page
		File, Prev, Next    string
	}{
		TitleID:    tid,
		PageNumber: pid,
		Title:      title,
		Page:       page,
		File:       fmt.Sprintf("/file/%d/%d.%s", page.TitleID, page.PageNumber, page.Ext),
		Prev:       "/",
		Next:       "/",
	}
	if page.PageNumber > 1 {
		data.Prev = fmt.Sprintf("/title/%d/%d", page.TitleID, page.PageNumber-1)
	}
	if page.PageNumber < title.PageCount {
		data.Next = fmt.Sprintf("/title/%d/%d", page.TitleID, page.PageNumber+1)
	}
	applyTemplate(w, "title-page", data)
}

// ReloadTitlePage перезагружает страницу из тайтла
func ReloadTitlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	tid, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	pid, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	u := r.FormValue("url")
	ext := r.FormValue("ext")
	db.InsertPage(tid, ext, u, pid)
	err = file.Load(tid, pid, u, ext)
	db.UpdatePageSuccess(tid, pid, err == nil)
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	applyTemplate(w, "success", "страница успешно перезакачана")
}
