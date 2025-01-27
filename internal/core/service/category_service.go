package service

import (
	"context"
	"portal-blog/internal/adapter/repository"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/lib/conv"

	"github.com/gofiber/fiber/v2/log"
)

type CategoryService interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategoryById(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategoryById(ctx context.Context, id int64) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

// CreateCategory implements CategoryService.
func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	slug := conv.GenerateSlug(req.Title)
	req.Slug = slug

	err = c.categoryRepository.CreateCategory(ctx, req)

	if err != nil {
		code = "[SERVICE] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}

	return nil

}

// DeleteCategoryById implements CategoryService.
func (c *categoryService) DeleteCategoryById(ctx context.Context, id int64) error {
	err = c.categoryRepository.DeleteCategoryById(ctx, id)
	if err != nil {
		code = "[SERVICE] DeleteCategoryById - 1"
    log.Errorw(code, err)
    return err
	}

	return nil
}

// EditCategoryById implements CategoryService.
func (c *categoryService) EditCategoryById(ctx context.Context, req entity.CategoryEntity) error {
	categoryData, err := c.categoryRepository.GetCategoryById(ctx, req.ID)
	if err != nil {
		code = "[SERVICE] EditCategoryById - 1"
		log.Errorw(code, err)
		return err
	}

	slug := conv.GenerateSlug(req.Title)
	if categoryData.Title == req.Title {
		slug = categoryData.Slug
	}

	req.Slug = slug

	err = c.categoryRepository.EditCategoryById(ctx, req)
	if err != nil {
		code = "[SERVICE] EditCategoryById - 2"
		log.Errorw(code, err)
		return err
	}

	return nil

}

// GetCategories implements CategoryService.
func (c *categoryService) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	results, err := c.categoryRepository.GetCategories(ctx)
	if err != nil {
		code = "[SERVICE] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, nil
}

// GetCategoryById implements CategoryService.
func (c *categoryService) GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	result, err := c.categoryRepository.GetCategoryById(ctx, id)
	if err != nil {
		code = "[SERVICE] GetCategoryById - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepo}
}
