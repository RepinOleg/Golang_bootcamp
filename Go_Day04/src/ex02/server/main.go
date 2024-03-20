package main

//#include "cow.h"
import "C"
import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unsafe"
)

type RequestBody struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

type SuccessResponseBody struct {
	Change int    `json:"change"`
	Thanks string `json:"thanks"`
}

type FailResponseBody struct {
	Error string `json:"error"`
}

var CandyPrices = map[string]int{
	"CE": 10,
	"AA": 15,
	"NT": 17,
	"DE": 21,
	"YR": 23,
}

func buyCandy(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Fatal(err)
	}
	if _, ok := CandyPrices[requestBody.CandyType]; !ok {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(FailResponseBody{
			Error: "Invalid candy type",
		})
		return
	}
	if requestBody.CandyCount < 0 {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(FailResponseBody{
			Error: "Negative candy count is forbidden",
		})
		return
	}
	if requestBody.Money >= CandyPrices[requestBody.CandyType]*requestBody.CandyCount {
		w.WriteHeader(http.StatusCreated)
		cs := C.CString("Thank you!")
		defer C.free(unsafe.Pointer(cs))
		askRes := C.ask_cow(cs)
		result := C.GoString(askRes)
		_ = json.NewEncoder(w).Encode(SuccessResponseBody{
			Change: requestBody.Money - CandyPrices[requestBody.CandyType]*requestBody.CandyCount,
			Thanks: result,
		})
		return
	} else {
		w.WriteHeader(402)
		diff := CandyPrices[requestBody.CandyType]*requestBody.CandyCount - requestBody.Money
		_ = json.NewEncoder(w).Encode(FailResponseBody{
			Error: fmt.Sprintf("You need %d more money!", diff),
		})
		return
	}
}

func getServer() *http.Server {
	server := &http.Server{
		Addr: ":3333",
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
		},
	}
	return server
}

func main() {
	http.HandleFunc("/buy_candy", buyCandy)

	server := getServer()
	err := server.ListenAndServeTLS("minica.pem", "minica-key.pem")
	if err != nil {
		log.Fatal(err)
	}
}
