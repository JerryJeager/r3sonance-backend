package users

import (
	"context"
	"errors"
	"time"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSv interface {
	CreateUser(ctx context.Context, spotifyProfile *models.SpotifyProfile, tokenResp *models.SpotifyTokenResponse) error
}

type UserServ struct {
	repo UserStore
}

func NewUserService(repo UserStore) *UserServ {
	return &UserServ{repo: repo}
}

func (s *UserServ) CreateUser(ctx context.Context, spotifyProfile *models.SpotifyProfile, tokenResp *models.SpotifyTokenResponse) error {
	oldUser, err := s.repo.GetUserByEmail(ctx, spotifyProfile.Email)
	if err == nil && oldUser != nil {
		oldUser.SpotifyID = spotifyProfile.ID
		oldUser.DisplayName = spotifyProfile.DisplayName
		oldUser.Country = spotifyProfile.Country
		oldUser.AccessToken = tokenResp.AccessToken
		oldUser.RefreshToken = tokenResp.RefreshToken
		oldUser.TokenExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		return s.repo.CreateUser(ctx, oldUser)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		user := &models.User{
			ID:             uuid.New(),
			SpotifyID:      spotifyProfile.ID,
			DisplayName:    spotifyProfile.DisplayName,
			Email:          spotifyProfile.Email,
			Country:        spotifyProfile.Country,
			AccessToken:    tokenResp.AccessToken,
			RefreshToken:   tokenResp.RefreshToken,
			TokenExpiresAt: time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
			CreatedAt:      time.Now(),
		}
		return s.repo.CreateUser(ctx, user)
	}
	return errors.New("failed to create or update user")
}
