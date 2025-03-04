package repository

import (
	"context"
	"fmt"
	"math"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/internal/core/domain/model"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContentRepository interface {
	GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error)
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
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
	}

	err = c.db.Create(&modelContent).Error
	if err != nil {
		code := "[REPOSITORY] CreateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteContent implements ContentRepository.
func (c *contentRepository) DeleteContent(ctx context.Context, id int64) error {
	err = c.db.Where("id = ?", id).Delete(&model.Content{}).Error
	if err != nil {
		code := "[REPOSITORY] DeleteContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetContentByID implements ContentRepository.
func (c *contentRepository) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	var modelContent model.Content

	err = c.db.Where("id = ?", id).Preload(clause.Associations).First(&modelContent).Error
	if err != nil {
		code := "[REPOSITORY] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	content := entity.ContentEntity{
		ID:          modelContent.ID,
		Title:       modelContent.Title,
		Excerpt:     modelContent.Excerpt,
		Description: modelContent.Description,
		Image:       modelContent.Image,
		Tags:        strings.Split(modelContent.Tags, ","),
		Status:      modelContent.Status,
		CategoryID:  modelContent.CategoryID,
		CreatedByID: modelContent.CreatedByID,
		CreatedAt:   modelContent.CreatedAt,
		Category: entity.CategoryEntity{
			ID:    modelContent.Category.ID,
			Title: modelContent.Category.Title,
			Slug:  modelContent.Category.Slug,
		},
		User: entity.UserEntity{
			ID:   modelContent.User.ID,
			Name: modelContent.User.Name,
		},
	}

	return &content, nil
}

// GetContents implements ContentRepository.
func (c *contentRepository) GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error) {
	var modelContents []*model.Content
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit
	status := ""
	if query.Status != "" {
		status = query.Status
	}

	sqlMain := c.db.Preload(clause.Associations).
		Where("title ILIKE ? or excerpt ILIKE ? OR description ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%").
		Where("status LIKE ?", "%"+status+"%")

	if query.CategoryID > 0 {
		sqlMain = sqlMain.Where("category_id = ?", query.CategoryID)
	}

	err = sqlMain.Model(&modelContents).Count(&countData).Error
	if err != nil {
		code := "[REPOSITORY] GetContents - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))

	err = sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelContents).Error

	if err != nil {
		code := "[REPOSITORY] GetContents - 3"
		log.Errorw(code, err)
		return nil, 0, 0, err
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
			Category: entity.CategoryEntity{
				ID:    v.Category.ID,
				Title: v.Category.Title,
				Slug:  v.Category.Slug,
			},
			User: entity.UserEntity{
				ID:   v.User.ID,
				Name: v.User.Name,
			},
		}
		contentsEntity = append(contentsEntity, content)
	}

	return contentsEntity, countData, int64(totalPages), nil

}

// UpdateContent implements ContentRepository.
func (c *contentRepository) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
	}

	err = c.db.Where("id = ?", req.ID).Updates(&modelContent).Error
	if err != nil {
		code := "[REPOSITORY] UpdateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}
