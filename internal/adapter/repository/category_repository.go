package repository

import (
	"context"
	"errors"
	"fmt"
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
	var countSlug int64
	err = c.db.Table("categories").Where("slug = ?", req.Slug).Count(&countSlug).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		code = "[REPOSITORY] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}

	countSlug += 1

	slug := fmt.Sprintf("%s-%d", req.Slug, countSlug)

	modelCategory := model.Category{
		ID:          0,
		Title:       req.Title,
		Slug:        slug,
		CreatedByID: req.User.ID,
	}

	err = c.db.Create(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] CreateCategory - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteCategoryById implements CategoryRepository.
func (c *categoryRepository) DeleteCategoryById(ctx context.Context, id int64) error {
	var count int64
	err = c.db.Table("contents").Where("category_id = ?", id).Count(&count).Error
	if err != nil {
		code = "[REPOSITORY] DeleteCategoryById - 1"
		log.Errorw(code, err)
		return err
	}

	if count > 0 {
		return errors.New("cannot delete a category that has associate contents")
	}

	err = c.db.Where("id = ?", id).Delete(&model.Category{}).Error
	if err != nil {
		code = "[REPOSITORY] DeleteCategoryById - 2"
		log.Errorw(code, err)
		return err
	}

	return nil

}

// EditCategoryById implements CategoryRepository.
func (c *categoryRepository) EditCategoryById(ctx context.Context, req entity.CategoryEntity) error {
	var countSlug int64
	err = c.db.Table("categories").Where("slug = ?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code = "[REPOSITORY] EditCategoryById - 1"
		log.Errorw(code, err)
		return err
	}

	countSlug += 1
	slug := req.Slug

	if countSlug == 0 {
		slug = fmt.Sprintf("%s-%d", req.Slug, countSlug)
	}

	modelCategory := model.Category{
		Title:       req.Title,
		Slug:        slug,
		CreatedByID: req.User.ID,
	}

	err = c.db.Where("id = ?", req.ID).Updates(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] EditCategoryById - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
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
	var categoryModel model.Category

	err = c.db.Where("id = ?", id).Preload("User").First(&categoryModel).Error
	if err != nil {
		code = "[REPOSITORY] GetCategoryById - 1"
		log.Errorw(code, err)
		return nil, err
	}

	categoryEntity := &entity.CategoryEntity{
		ID:    categoryModel.ID,
		Title: categoryModel.Title,
		Slug:  categoryModel.Slug,
		User: entity.UserEntity{
			ID:    categoryModel.User.ID,
			Name:  categoryModel.User.Name,
			Email: categoryModel.User.Email,
		},
	}

	return categoryEntity, nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
