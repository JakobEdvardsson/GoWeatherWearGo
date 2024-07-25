package storage

import (
	"time"

	"github.com/JakobEdvardsson/GoWeatherWearGo/types"
	"golang.org/x/oauth2"
)

// TODO: Add all the things
type Storage interface {
	GetUser(email string) (*types.User, error)
	AddUser(spotifyProfile *types.SpotifyProfileResponse) (*types.User, error)
	UpdateSpotifySession(refreshToken string, accessToken string, expiry time.Time, userId string) error
	CreateUserSession(token *oauth2.Token, user *types.User) (session *types.Session, err error)
	GetUserSession(sessionToken string) (session *types.Session, err error)
	GetAccount(userID string) (account *types.Account, err error)
}
