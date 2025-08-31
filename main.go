package main

import (
	"fmt"
	"log"
	"net/http"
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

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Printf("Starting server on port %s\n", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
