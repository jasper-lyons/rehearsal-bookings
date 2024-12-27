package handlers

import (
	"net/http"
	templates "rehearsal-bookings/web/templates"
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
