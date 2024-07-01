package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/JakobEdvardsson/GoWeatherWearGo/types"
)

func (s *Server) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	user := s.storage.Get(4)

	json.NewEncoder(w).Encode(user)
}

// ____    __    ____  _______     ___   .___________. __    __   _______ .______
// \   \  /  \  /   / |   ____|   /   \  |           ||  |  |  | |   ____||   _  \
//  \   \/    \/   /  |  |__     /  ^  \ `---|  |----`|  |__|  | |  |__   |  |_)  |
//   \            /   |   __|   /  /_\  \    |  |     |   __   | |   __|  |      /
//    \    /\    /    |  |____ /  _____  \   |  |     |  |  |  | |  |____ |  |\  \----.
//     \__/  \__/     |_______/__/     \__\  |__|     |__|  |__| |_______|| _| `._____|

func (s *Server) handleGetWeatherFromCords(w http.ResponseWriter, r *http.Request) {
	latitudeInput := r.URL.Query().Get("latitude")
	if latitudeInput == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required attribute latitude")
		return
	}

	longitudeInput := r.URL.Query().Get("longitude")
	if longitudeInput == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing required attribute longitude")
		return
	}

	lat, err := strconv.ParseFloat(latitudeInput, 64)
	if err != nil || lat < -90 || lat > 90 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		fmt.Fprintln(w, "Missing required attribute latitude")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(longitudeInput, 64)
	if err != nil || lon < -180 || lon > 180 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		fmt.Fprintln(w, "Missing required attribute longitude")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := http.Get(s.BASE_URL_WEATHER_API_CURRENT_WEATHER + "&q=" + latitudeInput + "," + longitudeInput)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		fmt.Fprintln(w, "Error getting data from WeatherAPI")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var weather types.WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	weatherCondition, ok := types.WEATHER_CONDITIONS[weather.Current.Condition.Code]
	if !ok {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
