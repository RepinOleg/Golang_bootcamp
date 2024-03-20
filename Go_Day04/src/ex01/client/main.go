package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

var req RequestBody

type RequestBody struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

func init() {
	flag.StringVar(&req.CandyType, "k", "", "candy name")
	flag.IntVar(&req.Money, "m", 0, "money")
	flag.IntVar(&req.CandyCount, "c", 0, "count")
	flag.Parse()
}

func getClient() (client *http.Client) {
	c, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	if err != nil {
		log.Fatal(err)
	}
	certs := []tls.Certificate{c}
	if len(certs) == 0 {
		client = &http.Client{}
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{Certificates: certs, InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
	return client
}

func main() {
	client := getClient()

	jsonValue, _ := json.Marshal(req)
	resp, err := client.Post("https://127.0.0.1:3333/buy_candy", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data))
	resp.Body.Close()
}
