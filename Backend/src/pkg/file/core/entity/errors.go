package entity

import "fmt"

// FileError represents an authentication error
type FileError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"-"`
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Error types
const (
	ErrInvalidMimeType = "INVALID_MIME_TYPE"
	ErrInvalidFileSize = "INVALID_FILE_SIZE"
)

// Error messages
const (
	MsgInvalidMimeType = "Invalid mime type"
	MsgInvalidFileSize = "Invalid file size"
)

// NewFileError creates a new authentication error
func NewFileError(errType, message string) *FileError {
	return &FileError{
		Type:    errType,
		Message: message,
	}
}

// NewFileErrorWithDetail creates a new authentication error with details
func NewFileErrorWithDetail(errType, message, detail string) *FileError {
	return &FileError{
		Type:    errType,
		Message: message,
		Detail:  detail,
	}
}
