package handlers

import (
	"net/http"
	"encoding/base64"
	"fmt"
	"os"
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

		fmt.Println(form.Username, username, form.Password, password)

		cookie := http.Cookie {
			Name: "StaticAuth",
			Value: base64.StdEncoding.EncodeToString([]byte(form.Username + ":" + form.Password)),
			Path: "/",
			HttpOnly: true,
			Secure: os.Getenv("APP_ENV") == "production",
			SameSite: http.SameSiteStrictMode,
			MaxAge: 2419200, // 28 days
		}

		http.SetCookie(w, &cookie)

		return Redirect("/admin/bookings")
	})
}
