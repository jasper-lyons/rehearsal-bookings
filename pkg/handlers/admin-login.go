package handlers

import (
	"net/http"
)

func AdminLogin() Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		return Template("admin-login.html.tmpl", nil)
	})
}
