package main

import (
	"time"

	"github.com/gin-contrib/cors"

	"bulls-lab-be/constants"
	"bulls-lab-be/internal/adapters/cron"
	"bulls-lab-be/internal/adapters/cron/jobs"
	"bulls-lab-be/internal/adapters/handler"
	"bulls-lab-be/internal/adapters/market"
	"bulls-lab-be/internal/adapters/repository"
	"bulls-lab-be/internal/core/services"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migrationPostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}

	db, err := ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := runMigrations(db); err != nil {
		log.Printf("⚠️  Migration warning: %v", err)
	}

	// Repositories
	userRepo := repository.NewPostgresRepo(db)
	watchlistRepo := repository.NewPostgresWatchlistRepo(db)
	watchlistStockRepo := repository.NewPostgresWatchlistStockRepo(db)

	// Services
	userService := services.NewUserService(userRepo)
	watchlistService := services.NewWatchlistService(watchlistRepo, watchlistStockRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	watchlistHandler := handler.NewWatchlistHandler(watchlistService)
	// Initialize repository
	orderRepo := repository.NewOrderRepository(db)
	holdingRepo := repository.NewStockHoldingRepository(db)

	// Initialize Market Service
	marketClient := market.NewMarketClient(constants.MARKET_SERVICE_URL)

	// Initialize service
	orderService := services.NewOrderService(orderRepo, holdingRepo, marketClient)
	portfolioService := services.NewPortfolioService(holdingRepo)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	orderHandler := handler.NewOrderHandler(orderService)
	portfolioHandler := handler.NewPortfolioHandler(portfolioService)

	// Initialize and start Cron Scheduler
	cronScheduler := cron.NewScheduler()

	// Register Jobs
	// Cron format: Second | Minute | Hour | Day of Month | Month | Day of Week

	heartbeatJob := jobs.NewHeartbeatJob()
	// Runs every minute at second 0
	if _, err := cronScheduler.RegisterJob("0 * * * * *", heartbeatJob); err != nil {
		log.Printf("⚠️  Failed to register heartbeat job: %v", err)
	}

	cancelExpiredOrdersJob := jobs.NewCancelExpiredOrdersJob(orderService)
	// Runs daily at 15:35:00
	if _, err := cronScheduler.RegisterJob("0 35 15 * * *", cancelExpiredOrdersJob); err != nil {
		log.Printf("⚠️  Failed to register cancel expired orders job: %v", err)
	}

	// Register Limit Order Executor Job - runs every 30 seconds
	limitOrderJob := jobs.NewLimitOrderJob(orderService, marketClient)
	if _, err := cronScheduler.RegisterJob("*/30 * * * * *", limitOrderJob); err != nil {
		log.Printf("⚠️  Failed to register limit order job: %v", err)
	}

	cronScheduler.Start()
	defer cronScheduler.Stop()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://your-frontend-domain.com", "https://bulls-lab-fe.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", healthHandler.HealthCheck)

	publicUserApi := r.Group("/api/v1/users")
	publicUserApi.POST("/register", userHandler.Register)
	publicUserApi.POST("/login", userHandler.Login)

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
			orders.GET("", orderHandler.GetOrdersByTab)
		}

		watchlists := api.Group("/watchlists")
		{
			watchlists.POST("", watchlistHandler.CreateWatchlist)
			watchlists.GET("", watchlistHandler.GetWatchlists)
			watchlists.PUT("/:id", watchlistHandler.UpdateWatchlist)
			watchlists.DELETE("/:id", watchlistHandler.DeleteWatchlist)
			watchlists.POST("/:id/stocks", watchlistHandler.AddStock)
			watchlists.DELETE("/:id/stocks", watchlistHandler.RemoveStock)
		}

		portfolio := api.Group("/portfolio")
		{
			portfolio.GET("/holdings", portfolioHandler.GetHoldings)
		}
	}

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
	driver, err := migrationPostgres.WithInstance(sqlDB, &migrationPostgres.Config{})
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
