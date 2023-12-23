package main

import (
	"flag"
	"net/http"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "Адрес веб сервера")
	dir := flag.String("dir", ".", "Директория для раздачи файлов")

	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*dir)))
	_ = http.ListenAndServe(*addr, nil)
}
