package api

func (s *Server) routes() {
	s.mux.Handle("/api/login", s.handleTokenLogin())
	s.mux.Handle("/api/logout", s.handleLogout())
}
