package main

import (
	"loaderData/internal/services"
	"log"
	"net/http"
)

func main() {
	http.Handle("/api/recommend/", http.HandlerFunc(services.HandlerNearestRest))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal(err)
	}
}
