package auth

import (
	"fmt"
	"portal-blog/config"
	"portal-blog/internal/core/domain/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	GenerateToken(data *entity.JwtData) (string, int64, error)
	VerifyAccessToken(token string) (*entity.JwtData, error)
}

type Options struct {
	SigningKey string
	Issuer     string
}

// GenerateToken implements Jwt.
// GenerateToken creates a new JWT token with the provided JwtData.
//
// This function generates a JWT token using the HS256 signing method. The token
// includes claims from the provided JwtData and additional standard claims such
// as expiration time, issuer, and not-before time.
//
// Parameters:
//   - data: A pointer to a JwtData struct containing the claims to be included in the token.
//
// Returns:
//   - string: The generated JWT token as a string.
//   - int64: The expiration time of the token as a Unix timestamp.
//   - error: An error if token generation fails, or nil if successful.
func (o *Options) GenerateToken(data *entity.JwtData) (string, int64, error) {
    now := time.Now().Local()
    expireAt := now.Add(time.Hour * 24)
    data.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expireAt)
    data.RegisteredClaims.Issuer = o.Issuer
    data.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
    acToken := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
    accessToken, err := acToken.SignedString([]byte(o.SigningKey))
    if err != nil {
        return "", 0, err
    }

    return accessToken, expireAt.Unix(), nil
}

// VerifyAccessToken implements Jwt.
// VerifyAccessToken validates and parses a JWT access token.
//
// It takes a token string as input and attempts to verify its signature and validity.
// If the token is valid, it extracts the claims and returns the user data.
//
// Parameters:
//   - token: A string representing the JWT access token to be verified.
//
// Returns:
//   - *entity.JwtData: A pointer to a JwtData struct containing the extracted user information if the token is valid.
//   - error: An error if the token is invalid, expired, or if there's any issue during the verification process.
func (o *Options) VerifyAccessToken(token string) (*entity.JwtData, error) {
    parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("signing method invalid")
        }

        return []byte(o.SigningKey), nil
    })

    if err != nil {
        return nil, err
    }

    if parsedToken.Valid {
        claim, ok := parsedToken.Claims.(jwt.MapClaims)
        if !ok || !parsedToken.Valid {
            return nil, err
        }

        jwtData := &entity.JwtData{
            UserID : claim["user_id"].(float64),
        }

        return jwtData, nil
    }

    return nil, fmt.Errorf("Token is not valid")
}

func NewJwt(cfg *config.Config) Jwt {
	apt := new(Options)
	apt.SigningKey = cfg.App.JwtSecretKey
	apt.Issuer = cfg.App.JwtIssuer

	return apt
}
