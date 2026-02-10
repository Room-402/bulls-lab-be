package main

import (
	"bulls-lab-be/internal/adapters/handler"
	"bulls-lab-be/internal/adapters/repository"
	"bulls-lab-be/internal/core/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initializing repository
	repo := repository.NewMemRepo()

	// Initializing service
	service := services.NewUserService(repo)

	// Initializing handlers
	userHandler := handler.NewUserHandler(service)
	healthHandler := handler.NewHealthHandler()

	// Setup Gin router
	r := gin.Default()

	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/users", userHandler.CreateUser)

	r.Run(":8080")
}
