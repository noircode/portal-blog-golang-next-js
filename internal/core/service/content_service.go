package service

import (
	"context"
	"portal-blog/config"
	"portal-blog/internal/adapter/cloudflare"
	"portal-blog/internal/adapter/repository"
	"portal-blog/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type ContentService interface {
	GetContents(ctx context.Context) ([]entity.ContentEntity, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
	UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error)
}

type contentService struct {
	contentRepository repository.ContentRepository
	cfg               *config.Config
	r2                cloudflare.CloudflareR2Adapter
}

// CreateContent implements ContentService.
func (c *contentService) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	err = c.contentRepository.CreateContent(ctx, req)
	if err != nil {
		code := "[SERVICE] CreateContent - 1"
    log.Errorw(code, err)
    return err
	}

	return nil
}

// DeleteContent implements ContentService.
func (c *contentService) DeleteContent(ctx context.Context, id int64) error {
	err = c.contentRepository.DeleteContent(ctx, id)
	if err != nil {
		code := "[SERVICE] DeleteContent - 1"
    log.Errorw(code, err)
    return err
	}

	return nil
}

// GetContentByID implements ContentService.
func (c *contentService) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	result, err := c.contentRepository.GetContentByID(ctx, id)

	if err != nil {
		code = "[SERVICE] GetContentByID - 1"
    log.Errorw(code, err)
    return nil, err
	}

	return result, nil
}

// GetContents implements ContentService.
func (c *contentService) GetContents(ctx context.Context) ([]entity.ContentEntity, error) {
	results, err := c.contentRepository.GetContents(ctx)
	if err != nil {
		code = "[SERVICE] GetContents - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, nil
}

// UpdateContent implements ContentService.
func (c *contentService) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	panic("unimplemented")
}

// UploadImageR2 implements ContentService.
func (c *contentService) UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error) {
	panic("unimplemented")
}

func NewContentService(repo repository.ContentRepository, cfg *config.Config, r2 cloudflare.CloudflareR2Adapter) ContentService {
	return &contentService{
		contentRepository: repo,
		cfg:               cfg,
		r2:                r2,
	}
}
