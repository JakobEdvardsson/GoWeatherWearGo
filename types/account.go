package types

import (
	"database/sql"

	"github.com/google/uuid"
)

type Account struct {
	ID                uuid.UUID      `json:"id" db:"id"`
	UserID            uuid.UUID      `json:"userId" db:"userId"`
	Type              string         `json:"type" db:"type"`
	Provider          string         `json:"provider" db:"provider"`
	ProviderAccountID string         `json:"providerAccountId" db:"providerAccountId"`
	RefreshToken      sql.NullString `json:"refreshToken" db:"refreshToken"`
	AccessToken       sql.NullString `json:"accessToken" db:"accessToken"`
	ExpiresAt         int64          `json:"expiresAt" db:"expiresAt"`
	TokenType         sql.NullString `json:"tokenType" db:"tokenType"`
	Scope             sql.NullString `json:"scope" db:"scope"`
	IDToken           sql.NullString `json:"idToken" db:"idToken"`
	SessionState      sql.NullString `json:"sessionState" db:"sessionState"`
}
