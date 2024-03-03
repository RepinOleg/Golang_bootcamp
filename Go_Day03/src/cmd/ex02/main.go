package main

import (
	"loaderData/internal/services"
	"log"
	"net/http"
)

func main() {
	http.Handle("/api/places/", http.HandlerFunc(services.HandlerJson))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal(err)
	}
}
