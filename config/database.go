package config

import (
	"fmt"
	"portal-blog/database/seeds"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

// ConnectionPostgres establishes a connection to a PostgreSQL database using the provided configuration.
//
// This function uses the Config struct to create a connection string, opens a connection to the database,
// and sets up connection pool parameters.
//
// Returns:
//   - *Postgres: A pointer to a Postgres struct containing the established database connection.
//   - error: An error if the connection fails, or nil if successful.
func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Psql.User,
		cfg.Psql.Password,
		cfg.Psql.Host,
		cfg.Psql.Port,
		cfg.Psql.DBName)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("[Connection Postgres 1] Failed to connect to database" + cfg.Psql.Host)
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Error().Err(err).Msg("[Connection Postgres 2] Failed to get database connection")
		return nil, err
	}

	seeds.SeedRoles(db)

	sqlDB.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Psql.DBMaxIdle)

	return &Postgres{DB: db}, nil
}
