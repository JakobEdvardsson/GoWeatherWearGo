package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/JakobEdvardsson/GoWeatherWearGo/storage"
	"github.com/JakobEdvardsson/GoWeatherWearGo/types"
	"github.com/JakobEdvardsson/GoWeatherWearGo/util"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var (
	spotifyOauthConfig *oauth2.Config
	oauthStateString   = "kebab"
)

type contextKey string

const spotifyClientKey = contextKey("spotifyClient")

func init() {
	// Load environment variables if they are not set
	clientId := os.Getenv("AUTH_SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("AUTH_SPOTIFY_CLIENT_SECRET")

	if clientId == "" || clientSecret == "" {
		err := util.LoadEnvFile(".env")
		if err != nil {
			log.Fatal("No env file or env vars provided!")
		}
	}

	// Initialize the OAuth2 configuration
	initOauthConfig()
}

func initOauthConfig() {
	spotifyOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("AUTH_SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_SPOTIFY_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"user-read-email"},
		Endpoint:     spotify.Endpoint,
	}
}

func handleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(spotifyOauthConfig)
	fmt.Println("spotifyOauthConfig")
	url := spotifyOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleSpotifyCallback(w http.ResponseWriter, r *http.Request, storage storage.Storage) {
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		return
	}

	token, err := spotifyOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Could not get token", http.StatusInternalServerError)
		return
	}

	// Get user info from Spotify API: GET https://api.spotify.com/v1/me

	profile, err := GetSpotifyUser(w, r, token.AccessToken)
	if err != nil || profile == nil {
		fmt.Println("profile, err := GetSpotifyUser(w, r):", err)
		return
	}
	fmt.Println(profile)

	// Check if user exists in DB, if not add user to DB

	user, err := storage.GetUser(profile.Email)
	if err == sql.ErrNoRows {
		fmt.Println("User does not exist in DB")
		user, err = storage.AddUser(profile)
		if err != nil {
			http.Error(w, "Error when adding user to DB", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Error when adding user to DB", http.StatusInternalServerError)
		return
	}
	fmt.Println("User: ", user)

	err = storage.UpdateSpotifySession(token.RefreshToken, token.AccessToken, token.Expiry, user.ID.String())
	if err != nil {
		http.Error(w, "Error updating account session in DB", http.StatusInternalServerError)
		return
	}

	session, err := storage.CreateUserSession(token, user)
	if err != nil || session.SessionToken == "" {
		http.Error(w, "Error creating session in DB", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:       "session_token",
		Value:      session.SessionToken,
		Path:       "",
		Domain:     "",
		Expires:    session.Expires,
		RawExpires: "",
		Secure:     false,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
	})
}

func RefreshSpotifyToken(refreshToken string, w http.ResponseWriter, r *http.Request) (response *types.RefreshTokenResponse, err error) {
	const refreshUrl = "https://accounts.spotify.com/api/token"
	clientID := os.Getenv("AUTH_SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("AUTH_SPOTIFY_CLIENT_SECRET")

	credentials := clientID + ":" + clientSecret
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	basicAuth := "Basic " + encoded

	formData := url.Values{}

	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", refreshToken)
	formData.Set("Authorization", basicAuth)
	formData.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequest("POST", refreshUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return nil, err
	}
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		http.Error(w, "Error getting data from WeatherAPI", http.StatusInternalServerError)
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return nil, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
