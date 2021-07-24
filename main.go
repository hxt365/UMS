package main

import (
	"Shopee_UMS/utils"
	"fmt"
)

func main() {
	//port := flag.String("port", ":8000", "port number to accept api requests")
	//flag.Parse()
	//
	//u := usecases.New(nil, nil, nil)
	//s := api.NewServer(u)
	//
	//log.Fatal(http.ListenAndServe(*port, s))

	fmt.Println(utils.RandString(32))
}
