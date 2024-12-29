package seeds

import (
	"portal-blog/internal/core/domain/model"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	bytes, err := bcrypt.GenerateFromPassword([]byte("admin123"), 14)
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