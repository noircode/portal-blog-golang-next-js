package handler

import (
	"portal-blog/internal/adapter/handler/request"
	"portal-blog/internal/adapter/handler/response"
	"portal-blog/internal/core/domain/entity"
	"portal-blog/internal/core/service"
	valid "portal-blog/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var err error
var code string
var errorResp response.ErrorResponseDefault

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

// Login implements AuthHandler.
func (a *authHandler) Login(c *fiber.Ctx) error {
	req := request.LoginRequest{}
	resp := response.SuccessAuthResponse{}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] Login - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = valid.ValidateStruct(req); err != nil {
		code = "[HANDLER] Login - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqLogin := entity.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := a.authService.GetUserByEmail(c.Context(), reqLogin)
	if err != nil {
		code = "[HANDLER] Login - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		if err.Error() == "invalid password" {
			return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Login Successfull"
	resp.AccessToken = result.AccessToken
	resp.ExpiredAt = result.ExpiredAt

	return c.JSON(resp)

}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}
