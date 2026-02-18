package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	SpotifyID      string    `json:"spotify_id"`
	PublicID       string    `json:"public_id"`
	DisplayName    string    `json:"display_name"`
	Email          string    `json:"email"`
	Country        string    `json:"country"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Country     string `json:"country"`
}
