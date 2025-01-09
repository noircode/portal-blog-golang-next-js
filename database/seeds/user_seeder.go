package seeds

import (
	"portal-blog/internal/core/domain/model"
	"portal-blog/lib/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// SeedRoles seeds the admin role into the database using the provided GORM database instance.
// It creates a new admin user with the given email and password, and hashes the password using the conv.HashPassword function.
// If the admin user already exists in the database, it updates the user's information.
// If any error occurs during the seeding process, it logs the error and exits the program.
// Otherwise, it logs a success message.
func SeedRoles(db *gorm.DB) {
    bytes, err := conv.HashPassword("admin123")
    if err != nil {
        log.Fatal().Err(err).Msg("Error creating password hash")
    }

    admin := model.User{
        Name:     "Admin",
        Email:    "admin@example.com",
        Password: string(bytes),
    }

    if err := db.FirstOrCreate(&admin, model.User{Email: "admin@mail.com"}).Error; err != nil {
        log.Fatal().Err(err).Msg("Error seeding admin role")
    } else {
        log.Info().Msg("Admin role seeded successfully")
    }
}