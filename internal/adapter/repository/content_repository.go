package repository

import (
	"context"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/internal/core/domain/model"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ContentRepository interface {
	GetContents(ctx context.Context) ([]entity.ContentEntity, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
}

type contentRepository struct {
	db *gorm.DB
}

// CreateContent implements ContentRepository.
func (c *contentRepository) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	panic("unimplemented")
}

// DeleteContent implements ContentRepository.
func (c *contentRepository) DeleteContent(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetContentByID implements ContentRepository.
func (c *contentRepository) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	
	panic("unimplemented")
}

// GetContents implements ContentRepository.
func (c *contentRepository) GetContents(ctx context.Context) ([]entity.ContentEntity, error) {
	var modelContents []*model.Content

	err := c.db.Order("created_at DESC").Preload("User", "Category").Find(&modelContents).Error
	if err != nil {
		code := "[REPOSITORY] GetContents - 1"
    log.Errorw(code, err)
    return nil, err
	}

	var contentsEntity []entity.ContentEntity

	for _, v := range modelContents {
		content := entity.ContentEntity{
      ID:          v.ID,
      Title:       v.Title,
      Excerpt:     v.Excerpt,
      Description: v.Description,
      Image:       v.Image,
      Tags:        strings.Split(v.Tags, ","),
      Status:      v.Status,
      CategoryID:  v.CategoryID,
			CreatedByID: v.CreatedByID,
      CreatedAt:   v.CreatedAt,
			Category:    entity.CategoryEntity{
				ID:    v.Category.ID,
        Title: v.Category.Title,
        Slug:  v.Category.Slug,
			},
			User: entity.UserEntity{
				ID: v.User.ID,
				Name: v.User.Name,
			},
    }
    contentsEntity = append(contentsEntity, content)
	}


	return contentsEntity, nil

}

// UpdateContent implements ContentRepository.
func (c *contentRepository) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	panic("unimplemented")
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}
