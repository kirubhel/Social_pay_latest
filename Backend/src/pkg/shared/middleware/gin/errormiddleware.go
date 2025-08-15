package gin

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/socialpay/socialpay/src/pkg/shared/errorxx"
	"github.com/socialpay/socialpay/src/pkg/shared/response"
)

func ErrorMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Next()

		// Check if any errors were set during the request processing
		if len(c.Errors) > 0 {

			// Iterate through all the errors in the context
			for _, err := range c.Errors {

				// Check if the error is an errorx error
				if exErr, ok := err.Err.(*errorx.Error); ok {
					// Extract the error code from errorx properties
					code, _ := exErr.Property(errorxx.ErrorCode)
					statusCode, _ := code.(int)

					// Log the error (optional, can be removed if not needed)

					log.Printf("errorx error: %v", exErr)
					errType := exErr.Type().String()
					parts := strings.Split(errType, ".")
					typeName := parts[len(parts)-1]

					// Send the response
					c.AbortWithStatusJSON(statusCode, response.ErrorResponse{
						Success: false,
						Error: response.ApiError{
							Type:    typeName,
							Message: err.Error(),
						},
					})
					return
				}

				// Handle other error types (non-errorx)
				// Default to HTTP 500 if we don't have a status code in the error
				statusCode := http.StatusInternalServerError
				message := err.Error()

				// Send the response for non-errorx errors
				c.AbortWithStatusJSON(statusCode, gin.H{
					"err":     err,
					"message": message,
				})
				return

			}
		}

	}
}
