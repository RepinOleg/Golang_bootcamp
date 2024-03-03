package types

type Store interface {
	GetPlaces(limit, offset int) ([]Place, int, error)
	GetClosestPlaces(lat, lon float64) ([]Place, error)
}

type Place struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

type Location struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Response struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	Total Total `json:"total"`
	Hits  []Hit `json:"hits"`
}

type Total struct {
	Value int `json:"value"`
}

type Hit struct {
	Source Place `json:"_source"`
}

type Token struct {
	Token string `json:"token"`
}
