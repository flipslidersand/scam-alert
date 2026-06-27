package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting Scam Alert API on port %s...\n", port)

	// TODO: Initialize HTTP server
	// TODO: Setup database connection
	// TODO: Setup Vision API client
	// TODO: Register routes

	log.Println("Server started")
}
