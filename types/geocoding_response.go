package types

type GeocodingLocation struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	URL     string  `json:"url"`
}

type GeocodingResponse []GeocodingLocation
