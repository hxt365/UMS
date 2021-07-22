package api

import (
	"Shopee_UMS/entities"
	"net/http"
)

func (s *Server) handleTokenLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			respondHTTPErr(w, r, http.StatusMethodNotAllowed)
			return
		}

		var req request
		err := decodeBody(r, &req)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.u.Auth.Authenticate(req.Username, req.Password)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}

		auth, ok := s.auth.(entities.TokenAuthenticator)
		if !ok {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}
		token, err := auth.GenerateToken(map[string]string{
			"username": req.Username,
		})
		if err != nil {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		setAuthToken(w, token)
		respond(w, r, http.StatusOK, nil)
	}
}

func (s *Server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			respondHTTPErr(w, r, http.StatusMethodNotAllowed)
			return
		}

		removeAuthToken(w)
		respond(w, r, http.StatusOK, nil)
	}
}
