package main

import (
	"github.com/elastic/go-elasticsearch/v8"
	"loaderData/internal/services"
	"log"
	"os"
)

const path = "testData/data.csv"

func main() {

	// Создаем клиент elastic
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	//считываем схему для маппинга
	schemaData, err := os.ReadFile("api/schema.json")
	if err != nil {
		log.Fatalf("Failed to read schema.json file: %v", err)
	}

	// Забираем данные из файла о заведениях
	places, err := services.GetDataFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	// Создаем индекс в elastic и записываем информацию в этот индекс
	services.WriteDataToElastic(places, client, schemaData)
}
