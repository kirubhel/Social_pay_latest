package service

import (
	"context"
	"mime/multipart"
)

// FileService define the interface for file management
type FileService interface {
	UploadFile(ctx context.Context, file multipart.File, fileHeader multipart.FileHeader) (*string, error)
}
