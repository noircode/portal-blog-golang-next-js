package middleware

import (
	"portal-blog/config"
	"portal-blog/internal/adapter/handler/response"
	"portal-blog/lib/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	CheckToken() fiber.Handler
}

type Options struct {
	authJwt auth.Jwt
}

// CheckToken returns a Fiber middleware handler that validates JWT tokens in the request header.
//
// This function creates a closure that checks for the presence of a valid Authorization header,
// extracts the JWT token, verifies it, and sets the resulting claims in the request context.
//
// Parameters:
//   - None directly, but uses the receiver (o *Options) which should have an authJwt field.
//
// Returns:
//   - fiber.Handler: A Fiber middleware handler function that performs token validation.
//     If the token is valid, it calls the next handler in the chain.
//     If the token is missing or invalid, it returns an unauthorized error response.
func (o *Options) CheckToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Missing or invalid Authorization header"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		claims, err := o.authJwt.VerifyAccessToken(tokenString)
		if err != nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Invalid or expired token"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		c.Locals("user", claims)

		return c.Next()
	}
}


// NewMiddleware creates and initializes a new Middleware instance.
//
// It takes a configuration object and sets up the necessary components
// for the middleware, including authentication.
//
// Parameters:
//   - cfg: A pointer to a config.Config struct containing the application configuration.
//
// Returns:
//   - Middleware: An interface that provides middleware functionality,
//     specifically for token checking in this implementation.
func NewMiddleware(cfg *config.Config) Middleware {
	opt := new(Options)
	opt.authJwt = auth.NewJwt(cfg)

	return opt
}
