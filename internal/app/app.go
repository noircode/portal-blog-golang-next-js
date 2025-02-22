package app

import (
	"context"
	"os"
	"os/signal"
	"portal-blog/config"
	"portal-blog/internal/adapter/cloudflare"
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

	err = os.MkdirAll("./temp/content/", 0775)
	if err != nil {
		log.Fatal().Msgf("Error creating temp directory: %v", err)
		return
	}

	// Cloudflare R2
	cdfR2 := cfg.LoadAwsConfig()
	s3Client := s3.NewFromConfig(cdfR2)
	r2Adapter := cloudflare.NewCloudflareR2Adapter(s3Client, cfg)

	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(cfg)

	_ = pagination.NewPagination()

	// Repository
	authRepo := repository.NewAuthRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	contentRepo := repository.NewContentRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)

	// Service
	authService := service.NewAuthService(authRepo, cfg, jwt)
	categoryService := service.NewCategoryService(categoryRepo)
	contentService := service.NewContentService(contentRepo, cfg, r2Adapter)
	userService := service.NewUserService(userRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	contentHandler := handler.NewContentHandler(contentService)
	userHandler := handler.NewUserHandler(userService)

	// Fiber App
	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] %{ip} %{status} - %{latency} %{method} %{path}\n",
	}))

	// Group API
	api := app.Group("/api")
	api.Post("/login", authHandler.Login)

	// Group Admin
	adminApp := api.Group("/admin")
	adminApp.Use(middlewareAuth.CheckToken())

	// Group Category
	categoryApp := adminApp.Group("/category")
	categoryApp.Get("/", categoryHandler.GetCategories)
	categoryApp.Post("/", categoryHandler.CreateCategory)
	categoryApp.Get("/:categoryID", categoryHandler.GetCategoryById)
	categoryApp.Put("/:categoryID", categoryHandler.EditCategoryById)
	categoryApp.Delete("/:categoryID", categoryHandler.DeleteCategoryById)

	// Group Content
	contentApp := adminApp.Group("/content")
	contentApp.Get("/", contentHandler.GetContents)
	contentApp.Get("/:contentID", contentHandler.GetContentByID)
	contentApp.Post("/", contentHandler.CreateContent)
	contentApp.Put("/:contentID", contentHandler.UpdateContent)
	contentApp.Post("/upload-image", contentHandler.UploadImageR2)
	contentApp.Delete("/:contentID", contentHandler.DeleteContent)

	// User
	userApp := adminApp.Group("/user")
	userApp.Get("/profile", userHandler.GetUserByID)
	userApp.Put("/update-password", userHandler.UpdatePassword)

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
