package main

import (
	"Shopee_UMS/api"
	"Shopee_UMS/usecases"
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", ":8000", "port number to accept api requests")
	flag.Parse()

	u := usecases.New(nil, nil, nil)
	s := api.NewServer(u)

	log.Fatal(http.ListenAndServe(*port, s))
}
