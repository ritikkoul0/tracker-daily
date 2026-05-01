package handler

import (
	"daily-tracker/database"
	"daily-tracker/handlers"
	"net/http"
	"os"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	router     *gin.Engine
	routerOnce sync.Once
)

// Initialize router once
func initRouter() {
	routerOnce.Do(func() {
		// Set Gin to release mode for production
		if os.Getenv("GO_ENV") == "production" {
			gin.SetMode(gin.ReleaseMode)
		}

		// Initialize database
		database.InitDatabase()

		// Create Gin router
		router = gin.New()
		router.Use(gin.Recovery())

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

		// Health check endpoint
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Daily Tracker API is running",
			})
		})
	})
}

// Handler is the entry point for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	initRouter()
	router.ServeHTTP(w, r)
}

// Made with Bob
