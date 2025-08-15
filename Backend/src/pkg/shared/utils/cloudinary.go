package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

func InitCloudinary() (*cloudinary.Cloudinary, error) {
	envFilePath := ".env"

	err := godotenv.Overload(envFilePath)
	if err != nil {
		log.Println("Error loading .env file:", err)
		return nil, fmt.Errorf("Error loading .env file: %w", err)
	}

	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)

	if err != nil {
		log.Println("Error initializing Cloudinary:", err)
		return nil, fmt.Errorf("Error initialization Cloudinary: %w", err)
	}

	return cld, nil
}
