package users

import (
	"context"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
	"gorm.io/gorm"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserMusicSnapshotByEmail(ctx context.Context, email string) (*models.UserMusicSnapshot, error)
	CreateUserMusicSnapshot(ctx context.Context, snapshot *models.UserMusicSnapshot) error
}

type UserRepo struct {
	client *gorm.DB
}

func NewUserRepo(client *gorm.DB) *UserRepo {
	return &UserRepo{client: client}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return r.client.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.client.WithContext(ctx).First(&user).Where("email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetUserMusicSnapshotByEmail(ctx context.Context, email string) (*models.UserMusicSnapshot, error) {
	var userMusicSnapshot models.UserMusicSnapshot

	err := r.client.WithContext(ctx).
		Table("user_music_snapshots as us").
		Select("us.*").
		Joins("inner join users as u on u.id = us.user_id").
		Where("u.email = ?", email).
		First(&userMusicSnapshot).Error

	if err != nil {
		return nil, err
	}

	return &userMusicSnapshot, nil
}

func (r *UserRepo) CreateUserMusicSnapshot(ctx context.Context, snapshot *models.UserMusicSnapshot) error {
	return r.client.WithContext(ctx).Save(snapshot).Error
}

