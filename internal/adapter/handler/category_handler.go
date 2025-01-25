package handler

import (
	"errors"
	"portal-blog/internal/adapter/handler/request"
	"portal-blog/internal/adapter/handler/response"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/internal/core/service"
	"portal-blog/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var defaultSuccessResponse response.DefaultSucessResponse

type CategoryHandler interface {
	GetCategories(c *fiber.Ctx) error
	GetCategoryById(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	EditCategoryById(c *fiber.Ctx) error
	DeleteCategoryById(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

// CreateCategory implements CategoryHandler.
func (ch *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req request.CategoryRequest
  claims := c.Locals("user").(*entity.JwtData)
  userId := claims.UserID
  if userId == 0 {
    code = "[HANDLER] CreateCategory - 1"
    err = errors.New("user not authorized")
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = "Unauthorized access"

    return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
  }

  if err = c.BodyParser(&req); err != nil {
    code = "[HANDLER] CreateCategory - 2"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusBadRequest).JSON(errorResp)
  }

  if err = validator.ValidateStruct(req); err != nil {
    code = "[HANDLER] CreateCategory - 3"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusBadRequest).JSON(errorResp)
  }

  reqEntity := entity.CategoryEntity{
    Title: req.Title,
    User: entity.UserEntity{
      ID: int64(userId),
    },
  }

  err = ch.categoryService.CreateCategory(c.Context(), reqEntity)
  if err != nil {
    code = "[HANDLER] CreateCategory - 4"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
  }

  defaultSuccessResponse.Meta.Status = true
  defaultSuccessResponse.Meta.Message = "Category created successfully"
  defaultSuccessResponse.Data = nil
  defaultSuccessResponse.Pagination = nil
  return c.JSON(defaultSuccessResponse)
}

// DeleteCategoryById implements CategoryHandler.
func (*categoryHandler) DeleteCategoryById(c *fiber.Ctx) error {
	panic("unimplemented")
}

// EditCategoryById implements CategoryHandler.
func (*categoryHandler) EditCategoryById(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetCategories implements CategoryHandler.
func (ch *categoryHandler) GetCategories(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
  userID := claims.UserID
  if userID == 0 {
    code = "[HANDLER] GetCategories - 1"
    err = errors.New("user not authorized")
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = "Unauthorized access"

    return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
  }

  results, err := ch.categoryService.GetCategories(c.Context())
  if err != nil {
    code = "[HANDLER] GetCategories - 2"
    log.Errorw(code, err)
    errorResp.Meta.Status = false
    errorResp.Meta.Message = err.Error()

    return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
  }

  categoryResponses := []response.SuccessCategoryResponse{}

  for _, result := range results {
    categoryResponse := response.SuccessCategoryResponse{
      ID:         result.ID,
      Title:      result.Title,
      Slug:       result.Slug,
      CreatedByName: result.User.Name,
    }
    categoryResponses = append(categoryResponses, categoryResponse)
  }

  defaultSuccessResponse.Meta.Status = true
  defaultSuccessResponse.Meta.Message = "Categories fetched successfully"
  defaultSuccessResponse.Pagination = nil
  defaultSuccessResponse.Data = categoryResponses

  return c.JSON(defaultSuccessResponse)
}

// GetCategoryById implements CategoryHandler.
func (*categoryHandler) GetCategoryById(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}
