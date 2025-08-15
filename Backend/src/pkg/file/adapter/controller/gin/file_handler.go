package gin

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	auth_service "github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	"github.com/socialpay/socialpay/src/pkg/file/core/service"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
)

// FileHandler handles file management routes
type FileHandler struct {
	authService auth_service.AuthService
	fileService service.FileService
	logger      logging.Logger
}

// NewFileHandler create new file handler
func NewFileHandler(authService auth_service.AuthService, fileService service.FileService) *FileHandler {
	return &FileHandler{
		authService: authService,
		fileService: fileService,
		logger:      logging.NewStdLogger("file_handler"),
	}
}

// Upload handles file upload
// @Summary Upload single file
// @Description Uploads single file to the system
// @Tags file
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{} "File uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /file/upload [post]
func (h *FileHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No file is received",
		})
		return
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to open file",
		})
		return
	}
	defer file.Close()

	url, err := h.fileService.UploadFile(c, file, *fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upload file",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"data": gin.H{
			"url": url,
		},
	})
}

// RegisterRoutes registers all file routes
func (h *FileHandler) RegisterRoutes(router *gin.RouterGroup) {
	jwtConfig := ginn.JWTAuthMiddlewareConfig{
		AuthService: h.authService,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Public:      false,
	}

	auth := router.Group("/file")
	{
		auth.POST("/upload",
			ginn.JWTAuthMiddleware(jwtConfig),
			h.Upload)
	}
}
