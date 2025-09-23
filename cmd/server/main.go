package main

import (
	"log"
	"net/http"

	"cloud/internal/api"
	"cloud/internal/gcp"
	"cloud/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize GCP service
	gcpService, err := gcp.NewComputeService()
	if err != nil {
		log.Fatal("Failed to initialize GCP service:", err)
	}

	// Initialize services
	deployService := services.NewDeployService(gcpService)

	// Initialize handlers
	deployHandler := api.NewDeployHandler(deployService)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware for frontend
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "auto-deployer",
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/deploy", deployHandler.Deploy)
		v1.GET("/status/:vm_id", deployHandler.GetStatus)
	}

	// Start server
	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// corsMiddleware adds CORS headers for frontend integration
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
