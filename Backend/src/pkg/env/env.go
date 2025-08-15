package env

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var requiredEnvVars = []string{
	// Database
	"DB_HOST",
	"DB_USER",
	"DB_PASS",
	"DB_NAME",
	"DB_PORT",
	"SSL_MODE",

	// M-PESA
	"SAFARICOM_USERNAME",
	"SAFARICOM_PASSWORD",
	"MPESA_BASE_URL",

	// Telebirr
	"TELEBIRR_SECURITY_CREDENTIAL",
	"TELEBIRR_PASSWORD",
	"TELEBIRR_BASE_URL",
	"TELEBIRR_SHORT_CODE",
	"TELEBIRR_IDENTITY_ID",

	// CBE
	"CBE_MERCHANT_ID",
	"CBE_MERCHANT_KEY",
	"CBE_TERMINAL_ID",
	"CBE_CREDENTIAL_KEY",
	"CBE_BASE_URL",

	// Cybersource
	"CYBERSOURCE_ACCESS_KEY",
	"CYBERSOURCE_PROFILE_ID",
	"CYBERSOURCE_SECRET_KEY",
	"CYBERSOURCE_BASE_URL",

	// SMS/KMI Cloud
	"KMI_ACCESS_KEY",
	"KMI_SECRET_KEY",
	"KMI_SMS_URL",
	"KMI_SMS_FROM",

	// General
	"AFRO_MESSAGE_API_KEY",
	"YIMULU_API_KEY",
	"ALLOWED_ORIGINS",
	"APP_ENV",
	"APP_URL_V2",
}

func LoadEnv() {
	envFilePath := ".env"

	// Check if .env file exists
	if _, err := os.Stat(envFilePath); err == nil {
		// .env file exists, try to load it
		if err := godotenv.Load(envFilePath); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Printf("Successfully loaded .env file from: %s", envFilePath)
		}
	} else {
		log.Printf("No .env file found at: %s, using environment variables from OS", envFilePath)
	}

	// Check required environment variables regardless of .env file
	if err := CheckEnv(); err != nil {
		log.Fatalf("Environment check failed: %v", err)
	}
}

func CheckEnv() error {
	var missingVars []string
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}
	return nil
}
