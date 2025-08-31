package routes

import (
	"net/http"
	"os"
	"time"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thelostleo/CTF-backend/controllers"
	"github.com/thelostleo/CTF-backend/database"
	"github.com/thelostleo/CTF-backend/middleware"
	"github.com/thelostleo/CTF-backend/models"
)

// NewRouter creates and configures the HTTP router
func NewRouter() *gin.Engine {
	// Connect to database
	database.ConnectDatabase()

	// Auto-migrate the database
	if err := models.MigrateAll(database.DB); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

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
	setupAPIRoutes(router)

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

// setupAPIRoutes configures the CTF API routes
func setupAPIRoutes(router *gin.Engine) {
	// Initialize controllers
	userController := &controllers.UserController{}
	challengeController := &controllers.ChallengeController{}
	adminController := &controllers.AdminController{}

	// API v1 group
	api := router.Group("/api/v1")

	// Public routes (no authentication required)
	public := api.Group("/")
	{
		// Public challenge viewing (without flags)
		public.GET("/challenges", challengeController.GetAllChallenges)
		public.GET("/challenges/:id", challengeController.GetChallengeByID)

		// Public leaderboard
		public.GET("/leaderboard", userController.GetLeaderboard)
	}

	// Rate limited public routes for authentication
	auth := api.Group("/")
	auth.Use(middleware.RateLimitMiddleware(10, time.Minute)) // 10 requests per minute
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	// Protected routes (authentication required)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User profile
		protected.GET("/profile", userController.GetProfile)

		// Token refresh
		protected.POST("/refresh-token", userController.RefreshToken)

		// Flag submission with rate limiting
		flagSubmission := protected.Group("/")
		flagSubmission.Use(middleware.FlagSubmissionRateLimit())
		{
			flagSubmission.POST("/challenges/:id/submit", challengeController.SubmitFlag)
		}
	} // Admin routes (authentication + admin privileges required)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware())
	{
		// Challenge management
		admin.POST("/challenges", adminController.CreateChallenge)
		admin.PUT("/challenges/:id", adminController.UpdateChallenge)
		admin.DELETE("/challenges/:id", adminController.DeleteChallenge)

		// User management
		admin.GET("/users", adminController.GetAllUsers)

		// Admin dashboard
		admin.GET("/dashboard", adminController.GetDashboard)
	}
}
