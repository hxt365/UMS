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

		authCookie := utils.ExtractReqAuthCookie(r, utils.AuthCookieKey)
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

		ctx := context.WithValue(r.Context(), utils.UidContextKey, jwtClaim.Uid)
		ctx = context.WithValue(ctx, utils.CsrfContextKey, jwtClaim.CsrfToken)
		next(w, r.WithContext(ctx))
	}
}

// withCSRF must be after withJwtAuth
func (s *Server) withCSRF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST", "PUT", "PATCH", "DELETE":
			csrfHeader := r.Header.Get(utils.CSRFHeaderName)
			csrfToken := r.Context().Value(utils.CsrfContextKey)
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