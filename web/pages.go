package web

import (
	"app/db"
	"app/file"
	"app/handler"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// PageLimit ограничение на выдачу
var PageLimit int

func applyTemplate(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
	}
}

// GetMainPage возвращает главную страницу
func GetMainPage(w http.ResponseWriter, r *http.Request) {
	applyTemplate(w, "main", map[string]interface{}{
		"Count":           db.SelectTitlesCount(),
		"UnloadCount":     db.SelectUnloadTitlesCount(),
		"PageCount":       db.SelectPagesCount(),
		"UnloadPageCount": db.SelectUnloadPagesCount(),
	})
}

// GetListPage возвращает страницу со списком тайтлов
func GetListPage(w http.ResponseWriter, r *http.Request) {
	titleCount := db.SelectTitlesCount()
	offset := 0
	limit := PageLimit
	pageCount := titleCount / limit
	if titleCount%limit > 0 {
		pageCount++
	}
	pages := []int{}
	for i := 1; i <= pageCount; i++ {
		pages = append(pages, i)
	}
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	if l := r.FormValue("count"); l != "" {
		limit, err = strconv.Atoi(r.FormValue("count"))
		if err != nil {
			applyTemplate(w, "error", err)
			return
		}
	}
	offset = (page - 1) * limit
	titles := db.SelectTitles(offset, limit)
	applyTemplate(w, "title-list", map[string]interface{}{
		"Count":  titleCount,
		"Titles": titles,
		"Offset": offset,
		"Limit":  limit,
		"Pages":  pages,
		"Page":   page,
	})
}

// NewTitle загружает новый тайтл
func NewTitle(w http.ResponseWriter, r *http.Request) {
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
	tid, err := strconv.Atoi(r.URL.Query().Get("title"))
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	pid, err := strconv.Atoi(r.URL.Query().Get("page"))
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
		data.Prev = fmt.Sprintf("/title/page?title=%d&page=%d", page.TitleID, page.PageNumber-1)
	}
	if page.PageNumber < title.PageCount {
		data.Next = fmt.Sprintf("/title/page?title=%d&page=%d", page.TitleID, page.PageNumber+1)
	}
	applyTemplate(w, "title-page", data)
}

// ReloadTitlePage перезагружает страницу из тайтла
func ReloadTitlePage(w http.ResponseWriter, r *http.Request) {
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
	err = db.InsertPage(tid, ext, u, pid)
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	err = file.Load(tid, pid, u, ext)
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	err = db.UpdatePageSuccess(tid, pid, err == nil)
	if err != nil {
		applyTemplate(w, "error", err)
		return
	}
	applyTemplate(w, "success", "страница успешно перезакачана")
}
