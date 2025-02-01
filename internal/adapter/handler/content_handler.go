package handler

import (
	"errors"
	"portal-blog/internal/adapter/handler/request"
	"portal-blog/internal/adapter/handler/response"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/internal/core/service"
	"portal-blog/lib/conv"
	"portal-blog/lib/validator"
	"strings"

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
func (ch *contentHandler) CreateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] CreateContent - 1"
    err := errors.New("user not authorized")
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = "Unauthorized access"

    return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	userID := claims.UserID

	var req request.ContentRequest
	if err = c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreateContent - 2"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validator.ValidateStruct(&req); err != nil {
		code := "[HANDLER] CreateContent - 3"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")

	reqEntity := entity.ContentEntity{
		Title:       req.Title,
    Excerpt:     req.Excerpt,
    Description: req.Description,
    Image:       req.Image,
    Tags:        tags,
    CategoryID:  req.CategoryID,
    Status:      req.Status,
    CreatedByID: int64(userID),
	}

	err = ch.contentService.CreateContent(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] CreateContent - 4"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Content created successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

// DeleteContent implements ContentHandler.
func (ch *contentHandler) DeleteContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetContentByID - 1"
    err := errors.New("user not authorized")
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = "Unauthorized access"

    return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] GetContentByID - 2"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = ch.contentService.DeleteContent(c.Context(), contentID)
	if err != nil {
		code := "[HANDLER] DeleteContent - 1"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Content Deleted Successfully"

	return c.JSON(defaultSuccessResponse)
}

// GetContentByID implements ContentHandler.
func (ch *contentHandler) GetContentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetContentByID - 1"
    err := errors.New("user not authorized")
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = "Unauthorized access"

    return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] GetContentByID - 2"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.contentService.GetContentByID(c.Context(), contentID)
	if err != nil {
		code := "[HANDLER] GetContentByID - 3"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	respContent := response.ContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		Status:       result.Status,
		CategoryID:   result.CategoryID,
		CreatedByID:  result.CreatedByID,
		CreatedAt:    result.CreatedAt.Local().String(),
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
	}
	
	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Data = respContent
	defaultSuccessResponse.Meta.Message = "Success"

	return c.JSON(defaultSuccessResponse)

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
      CreatedAt:    content.CreatedAt.Local().String(),
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
