package repository

import (
	"context"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	UpdatePassword(ctx context.Context, newPass string, id int64) error
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userRepository struct {
	db *gorm.DB
}

// GetUserByID retrieves a user from the database based on the provided user ID.
//
// Parameters:
//   - ctx: The context for handling request cancellations and timeouts.
//   - id: The unique identifier of the user to be retrieved.
//
// Returns:
//   - A pointer to an entity.UserEntity struct containing user details if found.
//   - An error if the user is not found or if there is an issue with the database query.
func (u *userRepository) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	var modelUser model.User
	err = u.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		code := "[REPOSITORY] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:    id,
		Name:  modelUser.Name,
		Email: modelUser.Email,
	}, nil
}

// UpdatePassword updates the password of a user in the database.
//
// Parameters:
//   - ctx: The context for handling request cancellations and timeouts.
//   - newPass: The new password to be set.
//   - id: The unique identifier of the user whose password is being updated.
//
// Returns:
//   - An error if the update fails.
func (u *userRepository) UpdatePassword(ctx context.Context, newPass string, id int64) error {
	err = u.db.Model(&model.User{}).Where("id = ?", id).Update("password", newPass).Error
	if err != nil {
		code := "[REPOSITORY] UpdatePassword - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// NewUserRepository creates a new instance of userRepository.
//
// Parameters:
//   - db: The GORM database connection instance.
//
// Returns:
//   - A new UserRepository instance.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
