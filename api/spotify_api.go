package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/JakobEdvardsson/GoWeatherWearGo/types"
	_ "github.com/lib/pq"
)

// Get user info from Spotify API: GET https://api.spotify.com/v1/me
func GetSpotifyUser(w http.ResponseWriter, r *http.Request, token string) (profile *types.SpotifyProfileResponse, error error) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		http.Error(w, "Error getting data from Spotify api", http.StatusInternalServerError)
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return nil, err
	}

	var spotifyProfile types.SpotifyProfileResponse
	err = json.Unmarshal(body, &spotifyProfile)
	if err != nil || spotifyProfile.Email == "" {
		http.Error(w, "Error fetching user spotify profile", http.StatusNotFound)
		return nil, err
	}

	return &spotifyProfile, nil
}
