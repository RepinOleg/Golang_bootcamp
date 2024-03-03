package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"loaderData/internal/types"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const indexName = "places"

func GetDataFromFile(path string) ([]*types.Place, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	info, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var (
		id     int
		lon    float64
		lat    float64
		result []*types.Place
	)

	for _, place := range info[1:] {
		id, err = strconv.Atoi(place[0])
		if err != nil {
			log.Println(err)
		}

		lon, err = strconv.ParseFloat(place[4], 64)
		if err != nil {
			log.Println(err)
		}

		lat, err = strconv.ParseFloat(place[5], 64)
		if err != nil {
			log.Println(err)
		}

		place := &types.Place{
			ID:      id + 1,
			Name:    place[1],
			Address: place[2],
			Phone:   place[3],
			Location: types.Location{
				Lon: lon,
				Lat: lat,
			},
		}
		result = append(result, place)
	}
	return result, nil
}

func WriteDataToElastic(places []*types.Place, client *elasticsearch.Client, schema []byte) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "places",
		Client: client,
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Indices.Delete([]string{indexName}, client.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}
	res.Body.Close()

	res, err = client.Indices.Create(indexName, client.Indices.Create.WithBody(strings.NewReader(string(schema))))
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}
	res.Body.Close()

	for _, place := range places {
		data, err := json.Marshal(place)
		if err != nil {
			log.Fatal(err)
		}

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(place.ID),
				Body:       bytes.NewReader(data),
			})

		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}

func WriteJson(w http.ResponseWriter, data interface{}, statusCode int) {
	resMarsh, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, _ = w.Write(resMarsh)
}
