package utils

import (
	"net/http"
	"time"
)

func RemoveCookies(w http.ResponseWriter, key string) {
	http.SetCookie(w, &http.Cookie{
		Name:    key,
		Expires: time.Unix(0, 0), // Set to past
		MaxAge:  -1,              // Also ensures deletion
		Path:    "/",
		Domain:  "",
	})
}

// set secure cookies
func SetCookies(w http.ResponseWriter, key string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:   key,
		Value:  value,
		MaxAge: 0,
		Path:   "/", Domain: "",
		Secure: false, HttpOnly: true,
	})
}
