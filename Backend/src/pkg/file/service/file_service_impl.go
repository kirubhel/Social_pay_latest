package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/file/core/entity"
	"github.com/socialpay/socialpay/src/pkg/file/core/service"
	"github.com/socialpay/socialpay/src/pkg/file/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// FileServiceImpl implements the FileService interface
type FileServiceImpl struct {
	cld    *cloudinary.Cloudinary
	logger *log.Logger
}

// NewFileService creates a new file service
func NewFileService(cld *cloudinary.Cloudinary, logger *log.Logger) service.FileService {
	return &FileServiceImpl{
		cld:    cld,
		logger: logger,
	}
}

// UploadFile handles file single upload to cloudinary
func (s *FileServiceImpl) UploadFile(ctx context.Context, file multipart.File, fileHeader multipart.FileHeader) (*string, error) {
	allowedMimeTypes := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
	}

	if fileHeader.Size > 10*1024*1024 {
		return nil, entity.NewFileError(entity.ErrInvalidFileSize, entity.MsgInvalidFileSize)
	}

	isValidFileType, err := utils.ValidateMimeType(&fileHeader, allowedMimeTypes)

	if err != nil {
		return nil, err
	}

	if !*isValidFileType {
		return nil, entity.NewFileError(entity.ErrInvalidMimeType, entity.MsgInvalidMimeType)
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)
	if err != nil {
		s.logger.Printf("ERR WHILE BUFFERING FILE::%v", err.Error())
		return nil, fmt.Errorf("failed to buffer file: %w", err)
	}

	res, err := s.cld.Upload.Upload(ctx, &buf, uploader.UploadParams{
		PublicID: fileHeader.Filename + uuid.NewString(),
		Folder:   "Doc",
	})

	if err != nil {
		s.logger.Printf("ERR WHILE UPLOADING FILE TO CLOUDINARY::%v", err.Error())
		return nil, fmt.Errorf("Failed to upload file to cloudinary : %w", err)
	}

	return &res.SecureURL, nil
}
