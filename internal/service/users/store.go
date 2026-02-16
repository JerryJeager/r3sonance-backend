package users

import (
	"context"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
	"gorm.io/gorm"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
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
