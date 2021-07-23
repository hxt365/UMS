package api

import (
	"Shopee_UMS/entities"
	"Shopee_UMS/usecases"
	"Shopee_UMS/utils"
	"net/http"
)

func NewTestServer(au usecases.AuthUsecaser, uu usecases.UserUsecaser) *Server {
	jwtAuth := entities.JwtAuthenticator{
		SecretKey:  utils.SampleRSAPrivateKey(),
		PublicKey:  utils.SampleRSAPublicKey(),
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
