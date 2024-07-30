package api

import (
	"net/http"
	"time"

	"github.com/JakobEdvardsson/GoWeatherWearGo/storage"
)

func AddCorsHeaderMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next(w, r)
	}
}

func SpotifyAuthMiddleware(next func(http.ResponseWriter, *http.Request), storage storage.Storage) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		cookieToken := cookie.Value

		// Check if DB session is expired
		session, err := storage.GetUserSession(cookieToken)
		if err != nil || session == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if time.Now().UTC().After(session.Expires.UTC()) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		account, err := storage.GetAccount(session.UserID.String())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//Check if spotify token is expired and if so, refresh it
		if time.Now().UTC().After(time.Unix(account.ExpiresAt, 0).UTC()) {
			if account.RefreshToken.Valid {
				response, err := RefreshSpotifyToken(account.RefreshToken.String, w, r)
				if err != nil || response == nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				date := time.Now().UTC().Add(time.Second * time.Duration(response.ExpiresIn))

				err = storage.UpdateSpotifySession(response.RefreshToken, response.AccessToken, date, account.UserID.String())
				if err != nil {
					http.Error(w, "Failed to refresh spotify session", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}
