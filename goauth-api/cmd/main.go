package main

import (
	"fmt"
	"go-chat/internals/adapters/cache"
	"go-chat/internals/adapters/handler"
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
	authService *services.AuthService
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

	redisCache, err := cache.NewRedisCache("127.0.0.1:6379", "")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&domain.User{}, &domain.Book{})

	store := repository.NewDB(db, redisCache)

	authService = services.NewAuthService(store)
	userService = services.NewUserService(store)
	bookService = services.NewBookService(store)
	InitRoutes()
}

func InitRoutes() {
	app := fiber.New()
	app.Use(cors.New())

	middlewareHandler := handler.NewAuthHandlers(authService)
	userHandler := handler.NewUserHandlers(userService)
	bookHandler := handler.NewBookHandlers(bookService)

	router := app.Group("/api")
	authRouter := router.Group("/auth")
	authRouter.Post("/register", userHandler.Register)
	authRouter.Post("/login", userHandler.LoginUser)
	authRouter.Get("/logout", userHandler.LogoutUser)
	authRouter.Get("/refresh", userHandler.RefreshTokens)

	router.Get("/books", middlewareHandler.Middleware, bookHandler.GetBooks)
	router.Post("/books", middlewareHandler.Middleware, bookHandler.CreateBook)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
