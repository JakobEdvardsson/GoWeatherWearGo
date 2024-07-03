package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/JakobEdvardsson/GoWeatherWearGo/types"
)

func (s *Server) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	user := s.storage.Get(4)

	json.NewEncoder(w).Encode(user)
}

//   _______  _______   ______     ______   ______    _______   __  .__   __.   _______
//  /  _____||   ____| /  __  \   /      | /  __  \  |       \ |  | |  \ |  |  /  _____|
// |  |  __  |  |__   |  |  |  | |  ,----'|  |  |  | |  .--.  ||  | |   \|  | |  |  __
// |  | |_ | |   __|  |  |  |  | |  |     |  |  |  | |  |  |  ||  | |  . `  | |  | |_ |
// |  |__| | |  |____ |  `--'  | |  `----.|  `--'  | |  '--'  ||  | |  |\   | |  |__| |
//  \______| |_______| \______/   \______| \______/  |_______/ |__| |__| \__|  \______|

func (s *Server) handleGetGeocodeFromCity(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start")
	city := r.PathValue("city")

	if city == "" {
		http.Error(w, "Missing required attribute city", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, s.BASE_URL_WEATHER_API_LOCATION+"&q="+city, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		http.Error(w, "Error getting data from WeatherAPI", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var geocoding types.GeocodingResponse
	err = json.Unmarshal(body, &geocoding)

	if err != nil || len(geocoding) < 1 {
		http.Error(w, "No search result found", http.StatusNotFound)
		return
	}

	type GeocodingLocation struct {
		Name    string  `json:"name"`
		Lat     float64 `json:"latitude"`
		Lon     float64 `json:"longitude"`
		Region  string  `json:"region"`
		Country string  `json:"country"`
	}

	response := GeocodingLocation{
		Name:    geocoding[0].Name,
		Lat:     geocoding[0].Lat,
		Lon:     geocoding[0].Lon,
		Region:  geocoding[0].Region,
		Country: geocoding[0].Country,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// ____    __    ____  _______     ___    ___________  __    __   _______  ______
// \   \  /  \  /   / |   ____|   /   \  |           ||  |  |  | |   ____||   _  \
//  \   \/    \/   /  |  |__     /  ^  \ `---|  |----`|  |__|  | |  |__   |  |_)  |
//   \            /   |   __|   /  /_\  \    |  |     |   __   | |   __|  |      /
//    \    /\    /    |  |____ /  _____  \   |  |     |  |  |  | |  |____ |  |\  \----.
//     \__/  \__/     |_______/__/     \__\  |__|     |__|  |__| |_______|| _| `._____|

func (s *Server) handleGetWeatherFromCords(w http.ResponseWriter, r *http.Request) {
	latitudeInput := r.URL.Query().Get("latitude")
	if latitudeInput == "" {
		http.Error(w, "Missing required attribute latitude", http.StatusBadRequest)
		return
	}

	longitudeInput := r.URL.Query().Get("longitude")
	if longitudeInput == "" {
		http.Error(w, "Missing required attribute longitude", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latitudeInput, 64)
	if err != nil || lat < -90 || lat > 90 {
		http.Error(w, "Missing required attribute latitude", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(longitudeInput, 64)
	if err != nil || lon < -180 || lon > 180 {
		http.Error(w, "Missing required attribute longitude", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, s.BASE_URL_WEATHER_API_CURRENT_WEATHER+"&q="+latitudeInput+","+longitudeInput, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		http.Error(w, "Error getting data from WeatherAPI", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var weather types.WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	weatherCondition, ok := types.WEATHER_CONDITIONS[weather.Current.Condition.Code]
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	type Response struct {
		Location       string               `json:"location"`
		LocalTime      string               `json:"local_time"`
		Precipitation  float64              `json:"precipitation"`
		Degrees        float64              `json:"degrees"`
		Condition      string               `json:"condition"`
		WeatherKeyword types.WeatherKeyword `json:"weather_keyword"`
		WeatherPicture string               `json:"weather_picture"`
	}

	response := Response{
		Location:       weather.Location.Name,
		LocalTime:      weather.Location.Localtime,
		Precipitation:  weather.Current.PrecipMm,
		Degrees:        weather.Current.TempC,
		Condition:      weather.Current.Condition.Text,
		WeatherKeyword: weatherCondition.WeatherKeyword,
		WeatherPicture: "/images/weather/" + string(weatherCondition.WeatherKeyword) + "svg",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
