package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"loaderData/internal/types"
	"loaderData/utils/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ElasticClient struct {
	Client *elasticsearch.Client
}

type Data struct {
	Places   []types.Place
	Total    int
	IsPrev   bool
	IsNext   bool
	PrevPage int
	NextPage int
	LastPage int
}

type JsonData struct {
	Name     string        `json:"name"`
	Total    int           `json:"total"`
	Places   []types.Place `json:"places"`
	PrevPage int           `json:"prev_page,omitempty"`
	NextPage int           `json:"next_page,omitempty"`
	LastPage int           `json:"last_page"`
}

type JsonClosestData struct {
	Name   string        `json:"name"`
	Places []types.Place `json:"places"`
}

func (es *ElasticClient) GetPlaces(limit, offset int) ([]types.Place, int, error) {
	res, err := es.Client.Search(
		es.Client.Search.WithIndex("places"),
		es.Client.Search.WithSize(limit),
		es.Client.Search.WithBody(strings.NewReader(`{	
			"sort": {
				"id": "asc"
			}
			}`)),
		es.Client.Search.WithFrom(offset),
		es.Client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	var response types.Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, 0, err
	}

	total := response.Hits.Total.Value
	places := make([]types.Place, len(response.Hits.Hits))

	for i, hit := range response.Hits.Hits {
		places[i] = hit.Source
	}
	return places, total, nil
}

func (es *ElasticClient) PrepareData(request *http.Request) (*Data, error) {
	key := "page"
	pageStr := request.URL.Query().Get(key)
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, errors.ErrorResponse{
			Err: fmt.Sprintf("Invalid '%s' value: '%v'", key, pageStr),
		}
	}

	// limit - количество записей
	limit := 10

	// offset - номер записи с которой нужно вывести
	offset := (page - 1) * limit

	//создаем интерфейс
	myStore := types.Store(&ElasticClient{Client: es.Client})

	// забираем данные из elastic
	places, total, err := myStore.GetPlaces(limit, offset)
	if err != nil {
		return nil, errors.ErrorResponse{
			Err: err.Error(),
		}
	}
	//
	result := Data{
		Places: places,
		Total:  total,
	}

	if offset > 0 {
		result.IsPrev = true
		result.PrevPage = page - 1
	}

	if offset+limit < total {
		result.IsNext = true
		result.NextPage = page + 1
	}
	result.LastPage = (total + limit - 1) / limit
	if page > result.LastPage {
		return nil, err
	}
	return &result, nil
}

func (es *ElasticClient) PrepareNearestRest(r *http.Request) ([]types.Place, error) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		return nil, err
	}
	lat, err := strconv.ParseFloat(u.Query().Get("lat"), 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(u.Query().Get("lon"), 64)
	if err != nil {
		return nil, err
	}
	myStore := types.Store(&ElasticClient{
		Client: es.Client,
	})

	return myStore.GetClosestPlaces(lat, lon)
}

func (es *ElasticClient) GetClosestPlaces(lat, lon float64) ([]types.Place, error) {
	query := map[string]interface{}{
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": lat,
						"lon": lon,
					},
					"order":           "asc",
					"unit":            "km",
					"mode":            "min",
					"distance_type":   "arc",
					"ignore_unmapped": true,
				},
			},
		},
	}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithIndex("places"),
		es.Client.Search.WithSize(3),
		es.Client.Search.WithBody(strings.NewReader(string(queryJson))),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var response types.Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	places := make([]types.Place, len(response.Hits.Hits))

	for i, hit := range response.Hits.Hits {
		places[i] = hit.Source
	}
	return places, nil
}

func CreateClient() (*ElasticClient, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	})
	if err != nil {
		return nil, err
	}

	settings := `{"index": {"max_result_window": 20000}}`
	req := esapi.IndicesPutSettingsRequest{
		Index: []string{"places"},
		Body:  strings.NewReader(settings),
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return &ElasticClient{Client: client}, nil
}
