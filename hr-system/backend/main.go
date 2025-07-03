package main

import (
	"log"
	"net/http"
	"os"

	"hr-backend/config"
	"hr-backend/routes"
	"hr-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// Initialize database
	config.InitDB()

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Apply middleware
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "service": "HR Backend"})
	})

	// Setup routes
	api := router.Group("/api/v1")
	routes.SetupEmployeeRoutes(api)
	routes.SetupDepartmentRoutes(api)
	routes.SetupLeaveRoutes(api)
	routes.SetupAuthRoutes(api)

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HR Backend server on port %s", port)
	log.Fatal(router.Run(":" + port))
}