package users

import (
	"context"
	"errors"
	"time"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
	"github.com/JerryJeager/r3sonance-backend/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSv interface {
	CreateUser(ctx context.Context, spotifyProfile *models.SpotifyProfile, tokenResp *models.SpotifyTokenResponse) error
	CreateUserMusicSnapshot(ctx context.Context, email string, snapshot *models.BuiltUserMusicSnapchat) error
	ShouldUpdateUserMusicSnapshot(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type UserServ struct {
	repo UserStore
}

func NewUserService(repo UserStore) *UserServ {
	return &UserServ{repo: repo}
}

func (s *UserServ) CreateUser(ctx context.Context, spotifyProfile *models.SpotifyProfile, tokenResp *models.SpotifyTokenResponse) error {
	oldUser, err := s.repo.GetUserByEmail(ctx, spotifyProfile.Email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		user := &models.User{
			ID:             uuid.New(),
			SpotifyID:      spotifyProfile.ID,
			PublicID:       utils.GeneratePublicID(8),
			DisplayName:    spotifyProfile.DisplayName,
			Email:          spotifyProfile.Email,
			Country:        spotifyProfile.Country,
			AccessToken:    tokenResp.AccessToken,
			RefreshToken:   tokenResp.RefreshToken,
			TokenExpiresAt: time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
			CreatedAt:      time.Now(),
		}
		return s.repo.CreateUser(ctx, user)
	} else if err == nil && oldUser != nil {
		oldUser.SpotifyID = spotifyProfile.ID
		oldUser.DisplayName = spotifyProfile.DisplayName
		oldUser.Country = spotifyProfile.Country
		oldUser.AccessToken = tokenResp.AccessToken
		oldUser.RefreshToken = tokenResp.RefreshToken
		oldUser.TokenExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		return s.repo.CreateUser(ctx, oldUser)
	}
	return err
}

func (s *UserServ) ShouldUpdateUserMusicSnapshot(ctx context.Context, email string) (bool, error) {
	userMusicSnapshot, err := s.repo.GetUserMusicSnapshotByEmail(ctx, email)
	if err != nil || userMusicSnapshot == nil {
		return true, nil
	}

	updatedAtPlusADay := userMusicSnapshot.UpdatedAt.Add(24 * time.Hour)

	if time.Now().After(updatedAtPlusADay) {
		return true, nil
	}

	return false, nil
}

func (s *UserServ) CreateUserMusicSnapshot(ctx context.Context, email string, snapshot *models.BuiltUserMusicSnapchat) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	oldSnapshot, err := s.repo.GetUserMusicSnapshotByEmail(ctx, email)

	if err == nil {
		// Update
		oldSnapshot.TopArtists = snapshot.TopArtists
		oldSnapshot.TopTracks = snapshot.TopTracks
		// oldSnapshot.AudioProfile = snapshot.AudioProfile
		oldSnapshot.ListeningPattern = snapshot.ListeningPattern
		// oldSnapshot.TopGenres = snapshot.TopGenre
		oldSnapshot.UpdatedAt = time.Now()

		return s.repo.CreateUserMusicSnapshot(ctx, oldSnapshot)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create
		newSnapshot := &models.UserMusicSnapshot{
			ID:         uuid.New(),
			UserID:     user.ID,
			TopArtists: snapshot.TopArtists,
			TopTracks:  snapshot.TopTracks,
			// TopGenres:        snapshot.TopGenre,
			// AudioProfile:     snapshot.AudioProfile,
			ListeningPattern: snapshot.ListeningPattern,
		}

		return s.repo.CreateUserMusicSnapshot(ctx, newSnapshot)
	}

	return err
}

func (s *UserServ) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
