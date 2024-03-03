package main

import (
	"loaderData/internal/services"
	"log"
	"net/http"
)

func main() {
	http.Handle("/api/get_token/", http.HandlerFunc(services.HandlerJWT))
	http.Handle("/api/recommend/", http.HandlerFunc(services.VerifyJWT))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal(err)
	}
}
