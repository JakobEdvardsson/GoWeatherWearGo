package api

import (
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
	http.HandleFunc("GET /user", s.handleGetUserById)

	http.HandleFunc("GET /api/weather", s.handleGetWeatherFromCords)

	return http.ListenAndServe(s.listenPort, nil)
}
