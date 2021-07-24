package api

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/utils"
	"context"
	"net/http"
)

func (s *Server) withJwtAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtAuth, ok := s.auth.(*entities.JwtAuthenticator)
		if !ok {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		authCookie := utils.ExtractReqCookie(r, "auth-token")
		if authCookie == nil {
			respondHTTPErr(w, r, http.StatusUnauthorized)
			return
		}

		token := authCookie.Value
		claim, err := jwtAuth.ValidateToken(token)
		if err != nil {
			respondHTTPErr(w, r, http.StatusUnauthorized)
			return
		}

		jwtClaim, ok := claim.(*entities.JwtClaim)
		if !ok {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "uid", jwtClaim.Uid)
		ctx = context.WithValue(ctx, "csrf", jwtClaim.CsrfToken)
		next(w, r.WithContext(ctx))
	}
}

// withCSRF must be after withJwtAuth
func (s *Server) withCSRF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST", "PUT", "PATCH", "DELETE":
			csrfHeader := r.Header.Get("X-CSRFToken")
			csrfToken := r.Context().Value("csrf")
			if csrfToken != csrfHeader {
				respondHTTPErr(w, r, http.StatusForbidden)
				return
			}
			next(w, r)
		default:
			next(w, r)
		}
	}
}