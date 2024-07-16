package types

import "github.com/google/uuid"

type Location struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Owner        uuid.UUID `json:"owner" db:"owner"`
	LocationName string    `json:"locationName" db:"location_name"`
	Latitude     float64   `json:"latitude" db:"latitude"`
	Longitude    float64   `json:"longitude" db:"longitude"`
}
