package storage

import "github.com/JakobEdvardsson/GoWeatherWearGo/types"

type Storage interface {
	Get(int) *types.User
}
