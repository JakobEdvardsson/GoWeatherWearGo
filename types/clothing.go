package types

import "github.com/google/uuid"

type Clothing struct {
	ID                     uuid.UUID `json:"id" db:"id"`
	Owner                  uuid.UUID `json:"owner" db:"owner"`
	ClothingType           string    `json:"clothingType" db:"clothing_type"`
	UsableTemperatureRange int       `json:"usableTemperatureRange" db:"usable_temperature_range"`
	Name                   string    `json:"name" db:"name"`
	IsPrecipitationProof   bool      `json:"isPrecipitationProof" db:"is_precipitation_proof"`
	IconPath               string    `json:"iconPath" db:"icon_path"`
}
