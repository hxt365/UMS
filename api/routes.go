package api

func (s *Server) routes() {
	// auth
	s.mux.Handle("/api/auth/login", s.handleTokenLogin())
	s.mux.Handle("/api/auth/logout", s.withJwtAuth(s.withCSRF(s.handleLogout())))
	// profile
	s.mux.Handle("/api/profile", s.withJwtAuth(s.withCSRF(s.handleProfile())))
}
