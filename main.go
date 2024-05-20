package main

import (
    "fmt"
    "log"
    "os"

    "github.com/jalong4/stock-service-go/config"
    "github.com/jalong4/stock-service-go/routes"
    "github.com/gin-gonic/gin"
)

func main() {

    config.LoadEnv()    // Load environment variables
    config.ConnectDB()  // Establish MongoDB connection

	router := gin.Default()

	// Setup routes
	routes.Setup(router)

    GenerateReadme(router) // Generate README.md file

    port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port for development
	}
	
    log.Printf("Starting server on port %s", port)
    router.Run(fmt.Sprintf(":%s", port))

}
