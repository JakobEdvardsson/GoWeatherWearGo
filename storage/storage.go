package storage

import "github.com/JakobEdvardsson/GoWeatherWearGo/types"

// TODO: Add all the things
type Storage interface {
	GetUser(email string) (*types.User, error)
	AddUser(spotifyProfile *types.SpotifyProfileResponse) (*types.User, error)
}
