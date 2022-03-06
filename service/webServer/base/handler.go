package base

import (
	"app/system"
	"encoding/json"
	"io"
	"net/http"
)

func StaticFileHandler(filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		system.Debug(r.Context(), "Запрос статического файла", r.URL.Path)
		http.ServeFile(w, r, filename)
	})
}

func FileWriteHandler(filename string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/*")
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
		if r.Method == http.MethodOptions {
			return
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
		if err := GetError(r.Context()); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			switch err {
			case ErrForbidden:
				w.WriteHeader(http.StatusForbidden)
			case ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
			case ErrNotAuth:
				w.WriteHeader(http.StatusUnauthorized)
			case ErrParseData:
				w.WriteHeader(http.StatusBadRequest)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			if _, err := io.WriteString(w, err.Error()); err != nil {
				system.Error(r.Context(), err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		data := GetResponse(r.Context())
		if data == nil {
			return
		}
		type XLSX interface {
			Write(w io.Writer) error
		}
		switch f := data.(type) {
		case XLSX:
			if err := f.Write(w); err != nil {
				system.Error(r.Context(), err)
			}
			return
		case io.Reader:
			if _, err := io.Copy(w, f); err != nil {
				system.Error(r.Context(), err)
			}
			return
		}
	})
}

func JSONWriteHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.Header().Set("Allow", "POST, OPTIONS")
			return
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
		enc := json.NewEncoder(w)
		if err := GetError(r.Context()); err != nil {
			switch err {
			case ErrForbidden:
				w.WriteHeader(http.StatusForbidden)
			case ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
			case ErrNotAuth:
				w.WriteHeader(http.StatusUnauthorized)
			case ErrParseData:
				w.WriteHeader(http.StatusBadRequest)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			if err := enc.Encode(err.Error()); err != nil {
				system.Error(r.Context(), err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		data := GetResponse(r.Context())
		if data == nil {
			return
		}
		if err := enc.Encode(data); err != nil {
			system.Error(r.Context(), err)
		}
	})
}

func PanicDefender(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			p := recover()
			if p != nil {
				system.Warning(r.Context(), "обнаружена паника", p)
				SetError(r, ErrPanicDetected)
			}
		}()
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func AddHandler(mux *http.ServeMux, uri string, next http.Handler) {
	mux.Handle(uri, JSONWriteHandler(PanicDefender(next)))
}
