package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thelostleo/CTF-backend/api/routes"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

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
