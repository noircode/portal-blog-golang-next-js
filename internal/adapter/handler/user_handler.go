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

type UserHandler interface {
	UpdatePassword(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

// GetUserByID retrieves a user by their ID.
//
// Input:
//   - c: *fiber.Ctx - The request context containing JWT claims.
//
// Output:
//   - error: Returns an error response if unauthorized or if fetching user details fails.
//
// It extracts the user ID from JWT claims and fetches user details from the service.
// Returns unauthorized status if the user is not authorized or an internal server error if fetching fails.
func (u *userHandler) GetUserByID(c *fiber.Ctx) error {
	// Extract JWT claims from request context
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetContentByID - 1"
		err := errors.New("user not authorized")
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	user, err := u.userService.GetUserByID(c.Context(), int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] GetContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	resp := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	defaultSuccessResponse.Data = resp

	return c.JSON(defaultSuccessResponse)
}

// UpdatePassword updates the user's password.
//
// Input:
//   - c: *fiber.Ctx - The request context containing JWT claims and password update request body.
//
// Output:
//   - error: Returns an error response if unauthorized, if request validation fails, or if updating the password fails.
//
// It extracts the user ID from JWT claims, validates the request body, and updates the user's password in the service.
func (u *userHandler) UpdatePassword(c *fiber.Ctx) error {
	// Extract JWT claims from request context
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetContentByID - 1"
		err := errors.New("user not authorized")
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.UpdatePasswordRequest
	if err = c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdatePassword - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validator.ValidateStruct(&req); err != nil {
		code := "[HANDLER] UpdatePassword - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = u.userService.UpdatePassword(c.Context(), req.NewPassword, int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService}
}
