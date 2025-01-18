package repository

import (
	"context"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var err error
var code string

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.UserEntity, error)
}

type authRepository struct {
	db *gorm.DB
}

// GetUserByEmail retrieves a user from the database based on their email address.
//
// Parameters:
//   - ctx: A context.Context for handling timeouts and cancellations.
//   - req: An entity.LoginRequest containing the email of the user to retrieve.
//
// Returns:
//   - *entity.UserEntity: A pointer to the UserEntity if found, containing user details.
//   - error: An error if the user is not found or if there's a database error, nil otherwise.
func (a *authRepository) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.UserEntity, error) {
    var modelUser model.User

    err = a.db.Where("email = ?", req.Email).First(&modelUser).Error
    if err != nil {
        code = "[REPOSITORY] GetUserByEmail - 1"
        log.Errorw(code, err)
        return nil, err
    }

    resp := entity.UserEntity{
        ID: modelUser.ID,
        Name: modelUser.Name,
        Email: modelUser.Email,
        Password: modelUser.Password,
    }

    return &resp, nil
}

// NewAuthRepository creates and returns a new instance of AuthRepository.
//
// Parameters:
//   - db: A pointer to a gorm.DB instance representing the database connection.
//
// Returns:
//   - AuthRepository: An interface that provides methods for authentication-related database operations.
func NewAuthRepository(db *gorm.DB) AuthRepository {
    return &authRepository{db: db}
}
