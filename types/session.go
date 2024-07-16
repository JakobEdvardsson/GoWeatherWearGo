package types

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"userId" db:"userId"`
	SessionToken string    `json:"sessionToken" db:"sessionToken"`
	Expires      time.Time `json:"expires" db:"expires"`
}
