package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the HTTP router
func NewRouter() *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://*",
			"https://*",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           5 * time.Minute,
	}))

	setupBasicRoutes(router)

	return router
}

// setupBasicRoutes configures the basic application routes
func setupBasicRoutes(router *gin.Engine) {
	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to CTF Backend API",
			"version": "1.0.0",
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"port":   os.Getenv("PORT"),
		})
	})
}
