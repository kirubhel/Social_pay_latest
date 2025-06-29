package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// GetPublicDir returns the absolute path to the public directory
func GetPublicDir() (string, error) {
	// Get the current working directory
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	log.Printf("[PATH] Current working directory: %s", workDir)

	// Construct the public directory path
	publicDir := filepath.Join(workDir, "public")
	log.Printf("[PATH] Constructed public directory path: %s", publicDir)

	// Check if directory exists
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		log.Printf("[PATH] Public directory does not exist, creating it at: %s", publicDir)
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(publicDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create public directory: %w", err)
		}
		log.Printf("[PATH] Successfully created public directory")
	} else {
		log.Printf("[PATH] Public directory exists at: %s", publicDir)
	}

	return publicDir, nil
}

// GetPublicFilePath returns the absolute path for a file in the public directory
func GetPublicFilePath(filename string) (string, error) {
	publicDir, err := GetPublicDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(publicDir, filename)
	log.Printf("[PATH] Constructed file path: %s", filePath)

	// Check if parent directory exists
	parentDir := filepath.Dir(filePath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		log.Printf("[PATH] Parent directory does not exist, creating it at: %s", parentDir)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create parent directory: %w", err)
		}
		log.Printf("[PATH] Successfully created parent directory")
	}

	return filePath, nil
}
