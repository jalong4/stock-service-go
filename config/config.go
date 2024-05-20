package config

import (
	"os"
)

// GetBasePath returns the base path for static files.
func GetBasePath() string {
    basePath := os.Getenv("PUBLIC_FILES_PATH")
    if basePath == "" {
        return "./public" // Default path for local development
    }
    return "/app/public" // Path inside the Docker container
}