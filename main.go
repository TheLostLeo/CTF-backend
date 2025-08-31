package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thelostleo/CTF-backend/api/routes"
	"github.com/thelostleo/CTF-backend/database"
	"github.com/thelostleo/CTF-backend/models"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Connect to PostgreSQL database
	database.ConnectDatabase()

	// Auto-migrate database schemas
	err = database.DB.AutoMigrate(&models.User{}, &models.Challenge{}, &models.Submission{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("✅ Database migration completed successfully!")

	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "6009"
	}

	router := routes.NewRouter()

	// Gin has built-in server, so we can use Run() method
	log.Printf("Starting CTF Backend server on port %s", portString)
	log.Printf("Server running at http://localhost:%s", portString)

	err = router.Run(":" + portString)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
