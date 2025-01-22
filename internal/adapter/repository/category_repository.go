package repository

import (
	"context"
	"errors"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategoryById(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategoryById(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *gorm.DB
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	panic("unimplemented")
}

// DeleteCategoryById implements CategoryRepository.
func (c *categoryRepository) DeleteCategoryById(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// EditCategoryById implements CategoryRepository.
func (c *categoryRepository) EditCategoryById(ctx context.Context, req entity.CategoryEntity) error {
	panic("unimplemented")
}

// GetCategories retrieves a list of categories from the database.
// It orders the categories by creation date in descending order and preloads the associated user.
//
// ctx: The context for the request.
//
// Returns:
// - A slice of CategoryEntity representing the retrieved categories.
// - An error if any occurred during the retrieval process.
func (c *categoryRepository) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
    var modelCategories []*model.Category

    err = c.db.Order("created_at desc").Preload("User").Find(&modelCategories).Error
    if err != nil {
        code = "[REPOSITORY] GetCategories - 1"
        log.Errorw(code, err)
        return nil, err
    }

    if len(modelCategories) == 0 {
        code = "[REPOSITORY] GetCategories - 2"
        err = errors.New("data not found")
        log.Errorw(code, err)
        return nil, err
    }

    var resps []entity.CategoryEntity
    for _, v := range modelCategories {
        resps = append(resps, entity.CategoryEntity{
            ID:    v.ID,
            Title: v.Title,
            Slug:  v.Slug,
            User: entity.UserEntity{
                ID:       v.User.ID,
                Name:     v.User.Name,
                Email:    v.User.Email,
                Password: v.User.Password,
            },
        })
    }

    return resps, nil
}

// GetCategoryById implements CategoryRepository.
func (c *categoryRepository) GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	panic("unimplemented")
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
