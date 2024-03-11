package main

import (
	"fmt"
	"go-chat/internals/adapters/handler"
	"go-chat/internals/adapters/handler/middleware"
	"go-chat/internals/adapters/repository"
	"go-chat/internals/config"
	"go-chat/internals/core/domain"
	"go-chat/internals/core/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	userService *services.UserService
	bookService *services.BookService
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := gorm.Open(postgres.Open(connStr))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&domain.User{}, &domain.Book{})
	store := repository.NewDB(db)
	bookService = services.NewBookService(store)
	userService = services.NewUserService(store)
	InitRoutes()
}

func InitRoutes() {
	app := fiber.New()
	app.Use(cors.New())
	authMiddleware := middleware.AuthMiddleware

	router := app.Group("/api")
	userHandler := handler.NewUserHandlers(userService)
	bookHandler := handler.NewBookHandlers(bookService)

	authRouter := router.Group("/auth")
	authRouter.Post("/login", userHandler.Login)
	authRouter.Post("/register", userHandler.Register)
	authRouter.Get("/refresh", userHandler.RefreshTokens)

	router.Get("/books", authMiddleware, bookHandler.GetBooks)
	router.Post("/books", authMiddleware, bookHandler.CreateBook)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
