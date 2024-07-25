package types

type SpotifyProfileResponse struct {
	Country     string `json:"country"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	//ExplicitContent ExplicitContent `json:"explicit_content"`
	ExternalUrls ExternalUrls `json:"external_urls"`
	Followers    Followers    `json:"followers"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Images       []Image      `json:"images"`
	Product      string       `json:"product"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

/*
type ExplicitContent struct {
	FilterEnabled bool `json:"filter_enabled"`
	FilterLocked  bool `json:"filter_locked"`
}
*/

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

type Followers struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
