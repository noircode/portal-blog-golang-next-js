package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/rs/zerolog/log"
)

// LoadAwsConfig loads AWS configuration using the provided API key and secret.
// It uses the AWS SDK v2's LoadDefaultConfig function to load the configuration.
// The function sets the credentials provider to use static credentials with the provided API key and secret.
// The region is set to "auto", which allows the AWS SDK to automatically determine the region.
// If an error occurs during the configuration loading process, the function logs the error and fatally exits.
// Otherwise, it logs a success message and returns the loaded AWS configuration.
func (cfg Config) LoadAwsConfig() aws.Config {
    conf, err := awsConfig.LoadDefaultConfig(context.TODO(),
        awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            cfg.R2.ApiKey, cfg.R2.ApiSecret, "",
        )),
        awsConfig.WithRegion("auto"),
    )

    if err != nil {
        log.Fatal().Msgf("Unable to load AWS Config: %v", err)
    }

    log.Info().Msg("Success Loaded AWS Config")

    return conf
}