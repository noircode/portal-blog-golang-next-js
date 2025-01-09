package conv

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash from the provided password.
//
// Parameters:
//   - password: The plain-text password to be hashed.
//
// Returns:
//   - string: The bcrypt hash of the password.
//   - error: An error if the hashing process fails, or nil if successful.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPasswordHash compares a plain-text password with a bcrypt hash.
//
// Parameters:
//   - password: The plain-text password to be checked.
//   - hash: The bcrypt hash to compare against.
//
// Returns:
//   - bool: True if the password matches the hash, false otherwise.
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// GenerateSlug creates a URL-friendly slug from a given title.
//
// Parameters:
//   - title: The original title string to be converted into a slug.
//
// Returns:
//   - string: A lowercase string with spaces replaced by hyphens, suitable for use in URLs.
func GenerateSlug(title string) string {
    slug := strings.ToLower(title)
    slug = strings.ReplaceAll(slug, " ", "-")

    return slug
}
