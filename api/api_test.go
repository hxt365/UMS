package api

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/test_utils"
	"Shopee_UMS/usecases"
	"net/http"
)

func NewTestServer(au usecases.AuthUsecaser, uu usecases.UserUsecaser) *Server {
	jwtAuth := entities.JwtAuthenticator{
		SecretKey:  test_utils.GenRSAPrivateKey(),
		PublicKey:  test_utils.GenRSAPublicKey(),
		ExpSeconds: 10,
		Issuer:     "Test Server",
	}
	s := &Server{
		mux: http.NewServeMux(),
		u: &usecases.Usecases{
			Auth: au,
			User: uu,
		},
		auth: &jwtAuth,
	}
	s.routes()
	return s
}
