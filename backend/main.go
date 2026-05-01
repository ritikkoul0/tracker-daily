package main

import (
	"daily-tracker/database"
	"daily-tracker/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize database
	database.InitDatabase()

	// Create Gin router
	router := gin.Default()

	// Configure CORS - Allow all origins for Vercel deployments
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = false // Must be false when AllowAllOrigins is true
	router.Use(cors.New(config))

	// API routes
	api := router.Group("/api")
	{
		// Activity CRUD
		api.POST("/activities", handlers.CreateActivity)
		api.GET("/activities", handlers.GetActivities)
		api.GET("/activities/:date", handlers.GetActivityByDate)
		api.PUT("/activities/:id", handlers.UpdateActivity)
		api.DELETE("/activities/:id", handlers.DeleteActivity)

		// Analysis endpoints
		api.GET("/analysis/weekly", handlers.GetWeeklyAnalysis)
		api.GET("/analysis/monthly", handlers.GetMonthlyAnalysis)
		api.GET("/analysis/yearly", handlers.GetYearlyAnalysis)
	}

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Made with Bob
