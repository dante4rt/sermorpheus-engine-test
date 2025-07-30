package main

import (
	"log"
	"sermorpheus-engine-test/internal/config"
	"sermorpheus-engine-test/internal/handlers"
	"sermorpheus-engine-test/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	dbService := services.NewDatabaseService(cfg.DatabaseURL)
	defer dbService.Close()

	eventService := services.NewEventService(dbService.DB)
	customerService := services.NewCustomerService(dbService.DB)
	rateService := services.NewRateService(dbService.DB)
	transactionService := services.NewTransactionService(
		dbService.DB,
		eventService,
		rateService,
		cfg.PlatformFeePercent,
	)

	eventHandler := handlers.NewEventHandler(eventService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	transactionHandler := handlers.NewTransactionHandler(transactionService, customerService)

	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "sermorpheus-engine",
		})
	})

	v1 := r.Group("/api/v1")
	{

		events := v1.Group("/events")
		{
			events.POST("", eventHandler.CreateEvent)
			events.GET("", eventHandler.GetEvents)
			events.GET("/:id", eventHandler.GetEventByID)
		}

		customers := v1.Group("/customers")
		{
			customers.GET("/:id", customerHandler.GetCustomer)
		}

		transactions := v1.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.POST("/:id/confirm", transactionHandler.ConfirmPayment)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
