package types

import "time"

type VerificationToken struct {
	Identifier string    `json:"identifier" db:"identifier"`
	Token      string    `json:"token" db:"token"`
	Expires    time.Time `json:"expires" db:"expires"`
}
