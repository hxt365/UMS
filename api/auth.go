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

		var (
			uid int
			err error
			ok bool
		)

		// check if the user has a valid auth token
		uid, ok = r.Context().Value(utils.UidContextKey).(int)
		if !ok {
			// if no, authenticate him/her
			var req request
			if err := decodeBody(r, &req); err != nil {
				respondErr(w, r, http.StatusBadRequest, "malformed request format")
				return
			}

			uid, err = s.u.Auth.Authenticate(req.Username, req.Password)
			if err != nil {
				if _, ok := err.(utils.AuthError); ok {
					respondErr(w, r, http.StatusBadRequest, err)
					return
				}
				respondHTTPErr(w, r, http.StatusInternalServerError)
				return
			}
		}

		auth, ok := s.auth.(entities.TokenAuthenticator)
		if !ok {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}
		csrfToken := utils.RandString(32)
		authToken, err := auth.GenerateToken(map[string]interface{}{
			"uid":  uid,
			"csrf": csrfToken,
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

		setAuthToken(w, authToken)
		setCsrfToken(w, csrfToken)
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
