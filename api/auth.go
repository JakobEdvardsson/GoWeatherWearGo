package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JakobEdvardsson/GoWeatherWearGo/storage"
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

// TODO: Add logic for adding session to DB session table
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

	http.SetCookie(w, &http.Cookie{
		Name:       "spotify_token",
		Value:      token.AccessToken,
		Path:       "",
		Domain:     "",
		Expires:    token.Expiry,
		RawExpires: "",
		Secure:     false,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
	})

	// Get user info from Spotify API: GET https://api.spotify.com/v1/me

	profile, err := GetSpotifyUser(w, r, token.AccessToken)
	if err != nil || profile == nil {
		fmt.Println("profile, err := GetSpotifyUser(w, r):", err)
		return
	}
	fmt.Println(profile)

	// Check if user exists in DB, if not add user to DB

	user, err := storage.GetUser(profile.Email)
	fmt.Println("getUser")
	if err == sql.ErrNoRows {
		fmt.Println("User does not exist in DB")
		user, err = storage.AddUser(profile)
		if err != nil {
			fmt.Println("Error when adding user to DB")
			fmt.Println(err)
			http.Error(w, "Error when adding user to DB", http.StatusBadRequest)
			// DB blew up
			return
		}
	} else if err != nil {
		fmt.Println("Error when getting user from DB")
		fmt.Println(err)
		// DB blew up
		return
	}

	fmt.Println("User: ", user)

	// Add session to DB, and check if session is valid

}

// TODO: Add DB check of session && move to middleware.go
func SpotifyAuthMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("spotify_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := &oauth2.Token{AccessToken: cookie.Value}
		client := spotifyOauthConfig.Client(r.Context(), token)

		// Add the client to the request context
		ctx := context.WithValue(r.Context(), spotifyClientKey, client)

		next(w, r.WithContext(ctx))
	})
}
