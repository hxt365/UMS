package utils

import (
	"net/http"
)

func ExtractReqAuthCookie(r *http.Request, name string) *http.Cookie {
	for _, c := range r.Cookies() {
		if c.Name == AuthCookieKey {
			return c
		}
	}
	return nil
}
