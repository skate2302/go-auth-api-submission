package main

import (
	"context"	
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"     
	ginSwagger "github.com/swaggo/gin-swagger" 

	_ "github.com/Skate2302/go-auth-api/docs" 
	
	"github.com/Skate2302/go-auth-api/internal/handlers"
	"github.com/Skate2302/go-auth-api/internal/db"
	"github.com/Skate2302/go-auth-api/internal/middleware" 
)

// @title           User Authentication API
// @version         1.0
// @description     This is a simple user authentication service using Go, Gin, and MongoDB.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found. Using default environment variables.")
	}

	mongoClient := db.ConnectToMongo()
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
		log.Println("Disconnected from MongoDB.")
	}()

	router := gin.Default()

	// This route will serve the auto-generated Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	
	{
		v1.POST("/signup", middleware.RateLimiterMiddleware(), handlers.SignUp(mongoClient))
		v1.POST("/login", handlers.Login(mongoClient))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" //use 8080 as a default
	}

	log.Printf("Server is starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}