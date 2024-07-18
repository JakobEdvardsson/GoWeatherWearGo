package storage

import (
	"github.com/JakobEdvardsson/GoWeatherWearGo/types"
	"golang.org/x/oauth2"
)

// TODO: Add all the things
type Storage interface {
	GetUser(email string) (*types.User, error)
	AddUser(spotifyProfile *types.SpotifyProfileResponse) (*types.User, error)
	RefreshAccountSession(token *oauth2.Token, user *types.User) error
}
