package api

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/usecases"
	"net/http"
)

type Server struct {
	mux  *http.ServeMux
	u    *usecases.Usecases
	auth entities.Authenticator
}

func NewServer(u *usecases.Usecases, auth entities.Authenticator) *Server {
	s := &Server{
		mux:  http.NewServeMux(),
		u:    u,
		auth: auth,
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
