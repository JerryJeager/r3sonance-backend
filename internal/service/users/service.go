package users

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
	"github.com/JerryJeager/r3sonance-backend/internal/service/compatibility"
	"github.com/JerryJeager/r3sonance-backend/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSv interface {
	CreateUser(ctx context.Context, spotifyProfile *models.SpotifyProfile, tokenResp *models.SpotifyTokenResponse) error
	CreateUserMusicSnapshot(ctx context.Context, email string, snapshot *models.BuiltUserMusicSnapchat) error
	ShouldUpdateUserMusicSnapshot(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserMusicSnapshot(ctx context.Context, email string) (*models.UserMusicSnapshot, error)
	GetMusicCompatibility(ctx context.Context, email, publicID string) (*models.CompatibilityResult, string, error)
}

type UserServ struct {
	repo UserStore
}

func NewUserService(repo UserStore) *UserServ {
	return &UserServ{repo: repo}
}

func (s *UserServ) CreateUser(ctx context.Context, spotifyProfile *models.SpotifyProfile, tokenResp *models.SpotifyTokenResponse) error {
	oldUser, err := s.repo.GetUserByEmail(ctx, spotifyProfile.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
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
		log.Println("creating user: ", user.ID)
		return s.repo.CreateUser(ctx, user)
	}

	if err != nil {
		return err
	}
	oldUser.SpotifyID = spotifyProfile.ID
	oldUser.DisplayName = spotifyProfile.DisplayName
	oldUser.Country = spotifyProfile.Country
	oldUser.AccessToken = tokenResp.AccessToken
	oldUser.RefreshToken = tokenResp.RefreshToken
	oldUser.TokenExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	log.Println("Updating user: ", oldUser.ID)
	return s.repo.UpdateUser(ctx, oldUser)
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

func (s *UserServ) GetUserMusicSnapshot(ctx context.Context, email string) (*models.UserMusicSnapshot, error) {
	snapshot, err := s.repo.GetUserMusicSnapshotByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	topArtists := snapshot.TopArtists[:5]
	topTracks := snapshot.TopTracks[:10]
	snapshot.TopArtists = topArtists
	snapshot.TopTracks = topTracks
	return snapshot, nil
}

func (s *UserServ) GetMusicCompatibility(ctx context.Context, email, publicID string) (*models.CompatibilityResult, string, error) {
	userASnapshot, err := s.repo.GetUserMusicSnapshotByEmail(ctx, email) //this is the snapshot of the person initiating the compatibility request; they were most likely sent a link by their friend to get a music compatibility result
	if err != nil {
		return nil, "", err
	}

	userProfile, err := s.repo.GetUserByPublicID(ctx, publicID) //this is the public id of the person user a is doing a music compatibility result with, we get the display name of this user that would be passed to the frontend for better ux
	if err != nil {
		return nil, "", err
	}
	userBSnapshot, err := s.repo.GetUserMusicSnapshotByPublicID(ctx, publicID)
	if err != nil {
		return nil, "", err
	}

	compatibilityResult := compatibility.CalculateCompatibility(*userASnapshot, *userBSnapshot)

	return &compatibilityResult, userProfile.DisplayName, nil
}
