package web

import (
	"app/config"
	"log"
	"net/http"
)

// Run запускает веб сервер
func Run(addr string) <-chan struct{} {
	mux := http.NewServeMux()
	// обработчик статики
	mux.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(config.DefaultFilePath))))

	mux.HandleFunc("/", GetMainPage)
	mux.HandleFunc("/title/list", GetListPage)
	mux.HandleFunc("/new", NewTitle)
	mux.HandleFunc("/prepare", SaveToZIP)
	mux.HandleFunc("/title/page", GetTitlePage)
	mux.HandleFunc("/reload/page", ReloadTitlePage)
	// создание объекта сервера
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
		// ReadTimeout:  1 * time.Minute,
		// WriteTimeout: 1 * time.Minute,
		// IdleTimeout:  1 * time.Minute,
	}
	done := make(chan struct{})
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
		close(done)
	}()
	return done
}
