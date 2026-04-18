package main

import (
	"time"

	"github.com/gin-contrib/cors"

	"bulls-lab-be/internal/adapters/handler"
	"bulls-lab-be/internal/adapters/repository"
	"bulls-lab-be/internal/core/services"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migrationPostgres "github.com/golang-migrate/migrate/v4/database/postgres" // ← Aliased
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres" // ← No alias needed
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}

	// Run migrations first
	// Connect to database
	db, err := ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := runMigrations(db); err != nil {
		log.Printf("⚠️  Migration warning: %v", err)
	}

	// Initialize repository
	repo := repository.NewPostgresRepo(db)
	orderRepo := repository.NewOrderRepository(db)

	// Initialize service
	service := services.NewUserService(repo)
	orderService := services.NewOrderService(orderRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(service)
	healthHandler := handler.NewHealthHandler()
	orderHandler := handler.NewOrderHandler(orderService)

	// Setup Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://your-frontend-domain.com"}, // Update with your frontend's URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	r.GET("/health", healthHandler.HealthCheck)

	// register api
	publicUserApi := r.Group("/api/v1/users")
	publicUserApi.POST("/register", userHandler.Register)
	publicUserApi.POST("/login", userHandler.Login)

	// Protected API routes
	api := r.Group("/api/v1")
	api.Use(handler.AuthMiddleware())
	{
		users := api.Group("/users")
		{
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateProfile)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("", userHandler.ListUsers)
		}
		orders := api.Group("/orders")
		{
			orders.POST("/create", orderHandler.CreateOrder)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func ConnectDB() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "bulls_lab")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	log.Println("✅ Successfully connected to the database!")
	return db, nil
}

func runMigrations(gormDB *gorm.DB) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	driver, err := migrationPostgres.WithInstance(sqlDB, &migrationPostgres.Config{}) // ← Use alias
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Migrations applied successfully!")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
