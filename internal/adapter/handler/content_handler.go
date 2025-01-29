package handler

import (
	"errors"
	"portal-blog/internal/adapter/handler/response"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/internal/core/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentByID(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	UpdateContent(c *fiber.Ctx) error
	DeleteContent(c *fiber.Ctx) error
	UploadImageR2(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

// CreateContent implements ContentHandler.
func (*contentHandler) CreateContent(c *fiber.Ctx) error {
	panic("unimplemented")
}

// DeleteContent implements ContentHandler.
func (*contentHandler) DeleteContent(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetContentByID implements ContentHandler.
func (*contentHandler) GetContentByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetContents implements ContentHandler.
func (ch *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetContents - 1"
    err := errors.New("user not authorized")
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = "Unauthorized access"

    return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := ch.contentService.GetContents(c.Context())

	if err != nil {
		code := "[HANDLER] GetContents - 2"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Successfully"

	respContents := []response.ContentResponse{}

	for _, content := range results {
		respContent := response.ContentResponse{
			ID:           content.ID,
      Title:        content.Title,
      Excerpt:      content.Excerpt,
      Description:  content.Description,
      Image:        content.Image,
      Tags:         content.Tags,
      Status:       content.Status,
      CategoryID:   content.CategoryID,
      CreatedByID:  content.CreatedByID,
      CreatedAt:    content.CreatedAt.Format(time.RFC3339),
      CategoryName: content.Category.Title,
      Author:       content.User.Name,
		}

		respContents = append(respContents, respContent)
	}

	defaultSuccessResponse.Data = respContents
	return c.JSON(defaultSuccessResponse)
}

// UpdateContent implements ContentHandler.
func (*contentHandler) UpdateContent(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UploadImageR2 implements ContentHandler.
func (*contentHandler) UploadImageR2(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{
		contentService: contentService,
	}
}
