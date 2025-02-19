package service

import (
	"context"
	"portal-blog/internal/adapter/repository"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/lib/conv"

	"github.com/gofiber/fiber/v2/log"
)

type UserService interface {
	UpdatePassword(ctx context.Context, newPass string, id int64) error
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userService struct {
	userRepository repository.UserRepository
}

// GetUserByID retrieves a user from the service layer by delegating to the repository.
//
// Parameters:
//   - ctx: The context for handling request cancellations and timeouts.
//   - id: The unique identifier of the user to be retrieved.
//
// Returns:
//   - A pointer to an entity.UserEntity struct containing user details if found.
//   - An error if the user is not found or if there is an issue with the repository call.
func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	result, err := u.userRepository.GetUserByID(ctx, id)
	if err != nil {
		code := "[SERVICE] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil

}

// UpdatePassword hashes the new password and updates it in the database.
//
// Parameters:
//   - ctx: The context for handling request cancellations and timeouts.
//   - newPass: The new plain text password to be hashed and updated.
//   - id: The unique identifier of the user whose password is being updated.
//
// Returns:
//   - An error if hashing fails or if the update operation fails.
func (u *userService) UpdatePassword(ctx context.Context, newPass string, id int64) error {
	password, err := conv.HashPassword(newPass)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 2"
		log.Errorw(code, err)
		return err
	}

	err = u.userRepository.UpdatePassword(ctx, password, id)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 3"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// NewUserService creates a new instance of userService.
//
// Parameters:
//   - userRepository: The repository instance to interact with the database.
//
// Returns:
//   - A new UserService instance.
func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}
