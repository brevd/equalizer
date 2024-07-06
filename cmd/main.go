package main

import (
	"log"

	"github.com/brevd/equalizer/internal"
	"github.com/brevd/equalizer/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set the Gin mode (release, debug, or test)
	gin.SetMode(gin.ReleaseMode)

	// Initialize the database
	internal.InitDatabase()

	// Initialize the router
	r := router.SetupRouter()

	// Start the HTTP server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
