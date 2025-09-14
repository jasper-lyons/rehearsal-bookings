package handlers

import (
	"encoding/json"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	templates "rehearsal-bookings/web/templates"
	"time"
	"strings"
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
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		if err != nil {
			http.Error(w, err.Error(), code)
		} else {
			http.Error(w, "", code)
		}
		return nil
	})
}

func Template(name string, data any) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		templates.Render(w, name, data)
		return nil
	})
}

func Redirect(url string) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		http.Redirect(w, r, url, http.StatusSeeOther)
		return nil
	})
}

func JSON[T any](object T) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(object)
		return nil
	})
}

// Middleware is code that should be run on every request
type LoggedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *LoggedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LoggedResponseWriter) Status() int {
	if w.status == 0 {
		return 200
	}
	return w.status
}

func Logging(next http.Handler) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		loggedWriter := &LoggedResponseWriter{
			ResponseWriter: w,
		}

		startTime := time.Now()

		next.ServeHTTP(loggedWriter, r)

		log.Println(
			time.Since(startTime).String(),
			r.Method,
			r.URL.Path,
			loggedWriter.Status(),
		)

		return nil
	})
}

func CreateStaticAuthMiddleware(username string, password string, loginPath string) func(Handler) Handler {
	return func(next Handler) Handler {
		return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
			cookie, err := r.Cookie("StaticAuth")
			if err != nil {
				return Redirect(loginPath)
			}

			decodedValue, err := base64.StdEncoding.DecodeString(cookie.Value)
			if err != nil {
				return Redirect(loginPath)
			}

			credentials := strings.Split(string(decodedValue), ":")
			incommingUsername := credentials[0]
			incommingPassword := credentials[1]

			if incommingUsername != username || incommingPassword != password {
				return Redirect(loginPath)
			}

			return next
		})
	}
}

func CreateBasicAuthMiddleware(username string, password string) func(Handler) Handler {
	return func(next Handler) Handler {
		return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
			user, pass, ok := r.BasicAuth()
			if ok && user == username && pass == password {
				return next
			} else {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				return Error(errors.New("Auth Failed!"), http.StatusUnauthorized)
			}
		})
	}
}
