package api

import (
	"fmt"
	"net/http"

	"github.com/JakobEdvardsson/GoWeatherWearGo/storage"
)

type Server struct {
	listenPort                           string
	storage                              storage.Storage
	apiKey                               string
	BASE_URL_WEATHER_API_CURRENT_WEATHER string
	BASE_URL_WEATHER_API_LOCATION        string
}

func NewServer(listenPort string, storage storage.Storage, apiKey string) *Server {
	return &Server{
		listenPort:                           ":" + listenPort,
		storage:                              storage,
		apiKey:                               apiKey,
		BASE_URL_WEATHER_API_CURRENT_WEATHER: "https://api.weatherapi.com/v1/current.json?key=" + apiKey,
		BASE_URL_WEATHER_API_LOCATION:        "https://api.weatherapi.com/v1/search.json?key=" + apiKey,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("GET /{$}", AddCorsHeaderMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello, World!")

		fmt.Fprintf(w, "Hello, World!")
	}))

	http.HandleFunc("GET /api/geocoding/{city}", AddCorsHeaderMiddleware(s.handleGetGeocodeFromCity))

	http.HandleFunc("GET /api/weather", AddCorsHeaderMiddleware(SpotifyAuthMiddleware(s.handleGetWeatherFromCords, s.storage)))

	http.HandleFunc("GET /api/auth/signin/spotify", AddCorsHeaderMiddleware(handleSpotifyLogin))

	http.HandleFunc("GET /api/auth/callback/spotify", AddCorsHeaderMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleSpotifyCallback(w, r, s.storage)
	}))

	// TODO: add CSRF in application
	//http.HandleFunc("GET /api/auth/csrf", nil)

	return http.ListenAndServe(s.listenPort, nil)
}
