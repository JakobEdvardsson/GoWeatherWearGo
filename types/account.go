package types

import "github.com/google/uuid"

type Account struct {
	ID                uuid.UUID `json:"id" db:"id"`
	UserID            uuid.UUID `json:"userId" db:"userId"`
	Type              string    `json:"type" db:"type"`
	Provider          string    `json:"provider" db:"provider"`
	ProviderAccountID string    `json:"providerAccountId" db:"providerAccountId"`
	RefreshToken      string    `json:"refreshToken" db:"refreshToken"`
	AccessToken       string    `json:"accessToken" db:"accessToken"`
	ExpiresAt         int64     `json:"expiresAt" db:"expiresAt"`
	TokenType         string    `json:"tokenType" db:"tokenType"`
	Scope             string    `json:"scope" db:"scope"`
	IDToken           string    `json:"idToken" db:"idToken"`
	SessionState      string    `json:"sessionState" db:"sessionState"`
}
