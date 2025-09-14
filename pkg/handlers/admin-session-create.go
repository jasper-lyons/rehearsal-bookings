package handlers

import (
	"net/http"
	"encoding/base64"
)

type AdminSessionCreateForm struct {
	Username string `json:"username"`
	Password string `json:"name"`
}

func AdminSessionCreate(username string, password string, loginPath string) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		form, err := ExtractForm[AdminSessionCreateForm](r)
		if err != nil {
			return Redirect(loginPath)
		}

		if form.Username != username || form.Password != password {
			return Redirect(loginPath)
		}

		cookie := http.Cookie {
			Name: "StaticAuth",
			Value: base64.StdEncoding.EncodeToString([]byte(form.Username + ":" + form.Password)),
			Path: "/admin/",
			HttpOnly: true,
			Secure: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge: 2419200, // 28 days
		}

		http.SetCookie(w, &cookie)

		return Redirect("/admin/bookings")
	})
}
