package types

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	Name          string         `json:"name" db:"name"`
	Email         string         `json:"email" db:"email"`
	EmailVerified sql.NullTime   `json:"emailVerified" db:"emailVerified"`
	Image         sql.NullString `json:"image" db:"image"`
}
