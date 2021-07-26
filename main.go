package main

import (
	"Shopee_UMS/api"
	database "Shopee_UMS/db"
	"Shopee_UMS/entities"
	"Shopee_UMS/reposistory"
	"Shopee_UMS/storage"
	"Shopee_UMS/usecases"
	"Shopee_UMS/utils"
	"log"
	"net/http"
	"strconv"
)

func main() {
	db, err := database.NewMySQL(10, 10)
	if err != nil {
		log.Fatal("could not connect to DB", err)
	}
	s3, err := storage.NewS3()
	if err != nil {
		log.Fatal("could not connect to AWS S3", err)
	}
	rdb, err := database.NewRedis()
	if err != nil {
		log.Fatal("could not connect to Redis", err)
	}

	accounts := reposistory.NewAccounts(db, rdb)
	users := reposistory.NewUsers(db, rdb)
	statics := reposistory.NewStatics(s3)

	u := usecases.New(accounts, users, statics)
	auth := newJwtAuthenticator()
	s := api.NewServer(u, auth)

	log.Println("start listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", s))
}

func newJwtAuthenticator() *entities.JwtAuthenticator {
	privateKey := utils.ReadRSAPrivateKey("./jwtRS256.key")
	publicKey := utils.ReadRSAPublicKey("./jwtRS256.key.pub")
	expSec, _ := strconv.Atoi(utils.MustEnv("AUTH_TOKEN_EXPIRATION_SECONDS"))

	return &entities.JwtAuthenticator{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		ExpSeconds: expSec,
		Issuer:     "Shopee UMS",
	}
}
