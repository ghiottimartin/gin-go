package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/cookit/backend/controllers"
	"github.com/cookit/backend/helpers"
	"github.com/cookit/backend/repositories"
	"github.com/cookit/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func main() {
	loadEnv()

	DB_USER := os.Getenv("DB_USER")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	dbString := "host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable"
	dbUrl := fmt.Sprintf(dbString, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error while connecting to database %s", err)
	}
	defer db.Close()

	var userRepository repositories.UserRepository
	userRepository = repositories.NewPostgresUserRepository(db)

	userService := services.NewUserService(userRepository)

	userController := controllers.NewUserController(userService)

	router := gin.Default()

	// Use the JWTMiddleware for the routes that require authorization
	router.Use(helpers.JWTMiddleware())

	router.POST("/login", userController.Login)
	router.GET("/user/:id", userController.GetUser)
	router.POST("/signup", userController.CreateUser)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("port is required")
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error while running server")
	}
}
