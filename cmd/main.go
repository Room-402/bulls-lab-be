package main

import (
	"bulls-lab-be/internal/adapters/handler"
	"bulls-lab-be/internal/adapters/repository"
	"bulls-lab-be/internal/core/services"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initializing repository
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	// Register Global Audit Callbacks
	repository.RegisterAuditCallbacks(db)

	repo := repository.NewPostgresRepo(db)

	// Initializing service
	service := services.NewUserService(repo)

	// Initializing handlers
	userHandler := handler.NewUserHandler(service)
	healthHandler := handler.NewHealthHandler()
	// Setup Gin router
	r := gin.Default()

	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/users/:id", userHandler.GetUser)
	r.GET("/users", userHandler.ListUsers)
	r.POST("/users", userHandler.CreateUser)

	r.Run(":8080")
}

func ConnectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify the connection is actually alive
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	fmt.Println("✅ Successfully connected to the database!")
	return db, nil
}
