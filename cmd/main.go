package main

import (
	"log"

	"e-wallet-api/internal/config"
	"e-wallet-api/internal/database"
	"e-wallet-api/internal/handlers"
	"e-wallet-api/internal/middleware"
	"e-wallet-api/internal/services"
	"e-wallet-api/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	gin.SetMode(cfg.Server.GinMode)

	database.InitDatabase(cfg)

	// Initialize RabbitMQ
	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitMQ.Close()

	if err := rabbitMQ.DeclareQueue(cfg.RabbitMQ.TransferQueue); err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	authService := services.NewAuthService()
	walletService := services.NewWalletService(rabbitMQ)
	queueService := services.NewQueueService(rabbitMQ, walletService)

	if err := queueService.StartTransferWorker(); err != nil {
		log.Fatal("Failed to start transfer worker:", err)
	}

	authHandler := handlers.NewAuthHandler(authService, cfg)
	topUpHandler := handlers.NewTopUpHandler(walletService)
	paymentHandler := handlers.NewPaymentHandler(walletService)
	transferHandler := handlers.NewTransferHandler(walletService)
	transactionHandler := handlers.NewTransactionHandler(walletService)
	profileHandler := handlers.NewProfileHandler(walletService)

	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("/topup", topUpHandler.TopUp)
		protected.POST("/payments", paymentHandler.Payment)
		protected.POST("/transfers", transferHandler.Transfer)
		protected.GET("/transactions", transactionHandler.GetTransactions)
		protected.PUT("/update-profile", profileHandler.UpdateProfile)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
