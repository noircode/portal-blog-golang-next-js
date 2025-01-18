package app

import (
	"context"
	"os"
	"os/signal"
	"portal-blog/config"
	"portal-blog/internal/adapter/handler"
	"portal-blog/internal/adapter/repository"
	"portal-blog/internal/core/service"
	"portal-blog/lib/auth"
	"portal-blog/lib/middleware"
	"portal-blog/lib/pagination"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

// RunServer initializes and starts the server application.
// It sets up the configuration, establishes a database connection,
// and handles any errors that occur during the process.
func RunServer() {
	cfg := config.NewConfig()

	db, err := cfg.ConnectionPostgres()

	if err != nil {
		log.Fatal().Msgf("Error connecting to database: %v", err)
		return
	}

	// Cloudflare R2
	cdfR2 := cfg.LoadAwsConfig()
	_ = s3.NewFromConfig(cdfR2)

	jwt := auth.NewJwt(cfg)
	_ = middleware.NewMiddleware(cfg)

	_ = pagination.NewPagination()

	// Repository
	authRepo := repository.NewAuthRepository(db.DB)


	// Service
	authService := service.NewAuthService(authRepo, cfg, jwt)

	// Handler
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] %{ip} %{status} - %{latency} %{method} %{path}\n",
	}))

	api := app.Group("/api")
	api.Post("/login", authHandler.Login)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatal().Msgf("Error starting server: %v", err)
    }
	}()

	quit := make(chan os.Signal, 3)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Print("server shutdown of 5 seconds\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
