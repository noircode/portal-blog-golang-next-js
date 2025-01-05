package app

import (
	"portal-blog/config"

	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// RunServer initializes and starts the server application.
// It sets up the configuration, establishes a database connection,
// and handles any errors that occur during the process.
func RunServer() {
	cfg := config.NewConfig()

	_, err := cfg.ConnectionPostgres()

	if err != nil {
		log.Fatal().Msgf("Error connecting to database: %v", err)
		return
	}

	// Cloudflare R2
	cdfR2 := cfg.LoadAwsConfig()
	_ = s3.NewFromConfig(cdfR2)

}
