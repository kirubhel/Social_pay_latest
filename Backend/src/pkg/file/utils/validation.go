package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

func GetMimeType(fileHeader *multipart.FileHeader) (*string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("faled to open file: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return nil, fmt.Errorf("failed to read the file: %w", err)
	}

	mimeType := http.DetectContentType(buffer)
	return &mimeType, nil
}

func ValidateMimeType(fileHeader *multipart.FileHeader, allowedMimeTypes map[string]bool) (*bool, error) {
	mimeType, err := GetMimeType(fileHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to get file mime type: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	isValid := true
	if ext == "" || !allowedMimeTypes[*mimeType] {
		isValid = false
	}
	return &isValid, nil
}
