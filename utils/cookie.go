package utils

import (
	"net/http"
)

func ExtractReqCookie(r *http.Request, name string) *http.Cookie {
	for _, c := range r.Cookies() {
		if c.Name == "auth-token" {
			return c
		}
	}
	return nil
}
