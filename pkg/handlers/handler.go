package handlers

import (
	"log"
	"net/http"
	templates "rehearsal-bookings/web/templates"
	"encoding/json"
)

// Infrastructure Type for nicer handler writing
type Handler func(w http.ResponseWriter, r *http.Request) Handler

func (handle Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if handler := handle(response, request); handler != nil {
		handler.ServeHTTP(response, request)
	}
}

// Handlers that represent values
func Error(err error, code int) Handler {
	return Handler(func (w http.ResponseWriter, r* http.Request) Handler {
		http.Error(w, err.Error(), code)
		return nil
	})
}

func Template(name string, data any) Handler {
	return Handler(func (w http.ResponseWriter, r* http.Request) Handler {
		templates.Render(w, name, data)
		return nil
	})
}

func Redirect(url string) Handler {
	return Handler(func (w http.ResponseWriter, r* http.Request) Handler {
		http.Redirect(w, r, url, http.StatusSeeOther)
		return nil
	})
}

func JSON(object any) Handler {
	return Handler(func (w http.ResponseWriter, r* http.Request) Handler {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(object)
		return nil
	})
}

// Middleware is code that should be run on every request
func Logging(next http.Handler) Handler {
	return Handler(func (w http.ResponseWriter, r* http.Request) Handler {
		defer func () {
			log.Println(r.Method, r.URL.Path)
		}()
		next.ServeHTTP(w, r)
		return nil
	})
}
