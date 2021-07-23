package api

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/utils"
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
		if err := decodeBody(r, &req); err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}

		uid, err := s.u.Auth.Authenticate(req.Username, req.Password)
		if err != nil {
			if _, ok := err.(utils.AuthError); ok {
				respondErr(w, r, http.StatusBadRequest, err)
				return
			}
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		auth, ok := s.auth.(entities.TokenAuthenticator)
		if !ok {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}
		token, err := auth.GenerateToken(map[string]interface{}{
			"uid": uid,
		})
		if err != nil {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		user, err := s.u.User.GetData(uid)
		if err != nil {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		setAuthToken(w, token)
		respond(w, r, http.StatusOK, user)
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
