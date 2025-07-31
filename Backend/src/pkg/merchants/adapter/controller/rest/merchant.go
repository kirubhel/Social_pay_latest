package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/socialpay/socialpay/src/pkg/merchants/core/entity"
	"github.com/socialpay/socialpay/src/pkg/merchants/errors"

	"github.com/google/uuid"
)

func (controller Controller) extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		err := errors.Error{
			Type:    "UNAUTHORIZED",
			Message: "Please provide a valid header token",
			Code:    http.StatusUnauthorized,
		}
		return "", err
	}
	return parts[1], nil
}

func (controller Controller) CreateKey(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [GetSignIn] ")

	token, err := controller.extractToken(r)

	if err != nil {

		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}
	fmt.Println("session")
	fmt.Println(session.User.Id)
	var req entity.Merchant
	// Parse request
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		controller.log.Println(" unable to decode user input::", err.Error())
		err = errors.Error{
			Type:    "BINDING_ERR",
			Message: "Unable to marshal user input",
			Code:    http.StatusBadRequest,
		}

		// Send error response
		SendJSONResponse(w,
			MerchantResponse{
				Success: false,
				Error:   err,
			}, errors.MapErrorToHTTPStatus(err))

		return
	}

	// merchantID, service, expiryDate, store string, commissionFrom bool
	req.UserID = session.User.Id
	apiKey, err := controller.interactor.CreateMerchant(req)
	if err != nil {
		SendJSONResponse(w, MerchantResponse{
			Success: false,
			Error:   err,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	SendJSONResponse(w, MerchantResponse{
		Data:    apiKey,
		Error:   nil,
		Success: true,
	}, http.StatusCreated)
}

func (controller Controller) AddMerchantInfo(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [ADD Merchant INFO] ")

	var merchantInfo entity.MerchantAdditionalInfo

	token, err := controller.extractToken(r)
	if err != nil {

		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}
	fmt.Println("session")
	fmt.Println(session.User.Id)

	// binding the request body to the struct
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&merchantInfo)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w,
			MerchantResponse{
				Success: false,
				Error: errors.Error{
					Type:    "BAD_REQUEST",
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				},
			}, http.StatusBadRequest)
		return
	}

	// usecase layer
	if err := controller.interactor.AddMerchantInfo(r.Context(),
		session.User.Id, merchantInfo); err != nil {
		SendJSONResponse(w, MerchantResponse{
			Data:    "",
			Error:   err,
			Success: false,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}
	SendJSONResponse(w, MerchantResponse{
		Data:    "Merchant info added successfully",
		Error:   nil,
		Success: true,
	}, http.StatusCreated)
	// Send success response

}

func (controller Controller) GetMerchantByUserID(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [GetSignIn] ")

	token, err := controller.extractToken(r)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}
	fmt.Println("session")
	fmt.Println(session.User.Id)

	var req entity.Merchant
	req.UserID = session.User.Id
	apiKey, err := controller.interactor.GetMerchantByUserID(req.UserID)
	if err != nil {
		SendJSONResponse(w, MerchantResponse{
			Data:    "",
			Error:   err,
			Success: false,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	SendJSONResponse(w, MerchantResponse{
		Data:    apiKey,
		Error:   nil,
		Success: true,
	}, http.StatusOK)
}

func (controller Controller) UploadDocument(w http.ResponseWriter, r *http.Request) {
	var token string
	var doc entity.MerchantDocument
	token, err := controller.extractToken(r)

	if err != nil {

		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, http.StatusUnauthorized)
		return
	}

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Println("UNAUTHORIZED TOKEN")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}
	fmt.Printf("SESSION ::%v ", session.User.Id)

	m, err := controller.interactor.GetMerchantByUserID(session.User.Id)

	if err != nil {
		controller.log.Println("Could not get merchant by user ID")
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, http.StatusUnauthorized)
		return
	}

	// r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		controller.log.Println("UNABLE TO PARSE FILE :: FILE TOO LARGE::", err.Error())
		err = errors.Error{
			Type:    "BAD_REQUEST",
			Message: "file size should be < 10MB",
			Code:    http.StatusBadRequest,
		}

		SendJSONResponse(
			w, MerchantResponse{
				Data:    "",
				Error:   err,
				Success: false,
			}, http.StatusBadRequest)

		return
	}

	doc.MerchantID, err = uuid.Parse(m.MerchantID)

	if err != nil {
		controller.log.Println("unable to parse merchant ID::", err.Error())

		SendJSONResponse(
			w, MerchantResponse{
				Data: "",
				Error: errors.Error{
					Type:    "INTERNAL_SERVER_ERROR",
					Message: "Unable to parse merchant Id",
					Code:    http.StatusInternalServerError,
				},
				Success: false,
			}, http.StatusInternalServerError,
		)
		return

	}

	doc.Status = "pending"
	// helper function
	handleUpload := func(fieldName, docType, errorMsg string, m *entity.Merchant) bool {
		allowedExtensions := map[string]bool{
			".pdf":  true,
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}

		allowedMimeTypes := map[string]bool{
			"application/pdf": true,
			"image/jpeg":      true,
			"image/png":       true,
		}

		file, fileHeader, err := r.FormFile(fieldName)
		if err == http.ErrMissingFile {
			if m.IsBettingClient && fieldName == "betting_certificate" {
				SendJSONResponse(w, MerchantResponse{
					Success: false,
					Data:    "",
					Error: &errors.Error{
						Type:    "VALIDATION_ERR",
						Message: "Betting certificate is required for betting clients",
						Code:    http.StatusBadRequest,
					},
				}, http.StatusBadRequest)
				return false
			}
			return true // Optional file, skip silently
		}

		if err != nil {
			controller.log.Printf("[UserID: %s] FILE UPLOAD ERR:: %s", session.User.Id, err.Error())
			SendJSONResponse(w, MerchantResponse{
				Success: false,
				Data:    "",
				Error: &errors.Error{
					Type:    "INTERNAL_SERVER_ERR",
					Message: errorMsg,
					Code:    http.StatusInternalServerError,
				},
			}, http.StatusInternalServerError)
			return false
		}
		defer file.Close()

		if fileHeader.Size > 10*1024*1024 {
			SendJSONResponse(w, MerchantResponse{
				Success: false,
				Data:    "",
				Error: &errors.Error{
					Type:    "VALIDATION_ERR",
					Message: "File too large. Must be < 10MB",
					Code:    http.StatusBadRequest,
				},
			}, http.StatusBadRequest)
			return false
		}

		// Read first 512 bytes to detect MIME type
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			controller.log.Printf("Unable to read file for MIME type detection: %v", err)
			SendJSONResponse(w, MerchantResponse{
				Success: false,
				Data:    "",
				Error: &errors.Error{
					Type:    "INTERNAL_SERVER_ERR",
					Message: "Failed to read file content",
					Code:    http.StatusInternalServerError,
				},
			}, http.StatusInternalServerError)
			return false
		}

		// Reset file reader
		_, _ = file.Seek(0, io.SeekStart)

		// Detect content type
		mimeType := http.DetectContentType(buffer)
		log.Println("Detected MIME type:", mimeType)

		// Extract extension
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		log.Println("File extension:", ext)

		// Validate either by extension or MIME type
		if ext == "" || !allowedExtensions[ext] {
			if !allowedMimeTypes[mimeType] {
				SendJSONResponse(w, MerchantResponse{
					Success: false,
					Data:    "",
					Error: &errors.Error{
						Type:    "VALIDATION_ERR",
						Message: "Invalid or missing file extension or MIME type",
						Code:    http.StatusBadRequest,
					},
				}, http.StatusBadRequest)
				return false
			}
		}

		doc.DocumentType = docType
		doc.CreatedAt = time.Now()

		if err := controller.interactor.AddDocument(r.Context(), file, *fileHeader, doc); err != nil {
			SendJSONResponse(w, MerchantResponse{
				Success: false,
				Data:    "",
				Error:   err,
			}, errors.MapErrorToHTTPStatus(err))
			return false
		}

		return true
	}
	if !handleUpload("tin", "tin", "ERROR UPLOADING TIN DOCUMENT", m) {
		return
	}
	if !handleUpload("business_license", "business_license", "ERROR UPLOADING BUSINESS LICENSE DOCUMENT", m) {
		return
	}
	if m.IsBettingClient {
		if !handleUpload("betting_license", "betting_license", "ERROR UPLOADING BETTING LICENSE DOCUMENT", m) {
			return
		}
	}

	SendJSONResponse(w, MerchantResponse{
		Data:    "file uploaded successfully",
		Error:   nil,
		Success: true,
	}, http.StatusCreated)
}

func (controller Controller) GetMerchants(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [GetMerchants] ")

	// Authentication check
	token, err := controller.extractToken(r)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	// Session validation
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Println("Authentication failed:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}

	controller.log.Printf("Fetching merchants for user %s", session.User.Id)

	// Get merchants from interactor
	merchants, err := controller.interactor.GetMerchants()
	if err != nil {
		controller.log.Printf("Error fetching merchants: %v", err)
		SendJSONResponse(w, MerchantResponse{
			Data:    nil,
			Error:   err,
			Success: false,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	// Successful response
	SendJSONResponse(w, MerchantResponse{
		Data:    merchants,
		Error:   nil,
		Success: true,
	}, http.StatusOK)
}

// controller/merchant.go
func (controller Controller) UpdateMerchantStatus(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[MERCHANT] [CONTROLLER] [UpdateStatus] ")

	// Authentication
	token, err := controller.extractToken(r)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	// Authorization - check if user has permission
	_, err = controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Println("Authorization failed:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "UNAUTHORIZED",
				Message: "Not authorized to update merchant status",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Parse request
	var request struct {
		MerchantID string `json:"merchant_id"`
		Status     string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate merchant ID
	merchantID, err := uuid.Parse(request.MerchantID)
	if err != nil {
		SendJSONResponse(w, MerchantsResponse{
			Success: false,
			Error: &errors.Error{
				Type:    "INVALID_INPUT",
				Message: "Invalid merchant ID format",
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate status
	status := entity.MerchantStatus(request.Status)
	if status != entity.StatusActive && status != entity.StatusPending && status != entity.StatusSuspended {
		SendJSONResponse(w, MerchantsResponse{
			Success: false,
			Error: &errors.Error{
				Type:    "INVALID_INPUT",
				Message: "Status must be either 'active', 'pending', or 'suspended'",
			},
		}, http.StatusBadRequest)
		return
	}

	// Call usecase
	err = controller.interactor.UpdateMerchantStatus(merchantID, status)
	if err != nil {
		SendJSONResponse(w, MerchantsResponse{
			Success: false,
			Error:   err,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	// Success response
	SendJSONResponse(w, MerchantsResponse{
		Success: true,
		Message: "Merchant status updated successfully",
	}, http.StatusOK)
}

func (controller Controller) UpdateFullMerchant(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[MERCHANT] [CONTROLLER] [UPDATE] ")

	// Authentication
	token, err := controller.extractToken(r)
	if err != nil {
		SendJSONResponse(w, Response{Success: false, Error: err}, errors.MapErrorToHTTPStatus(err))
		return
	}

	// Authorization
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   &errors.Error{Type: "UNAUTHORIZED", Message: "Not authorized"},
		}, http.StatusUnauthorized)
		return
	}

	// Parse request
	var request struct {
		MerchantID uuid.UUID                     `json:"merchant_id"`
		Merchant   entity.Merchant               `json:"merchant"`
		Address    entity.MerchantAdditionalInfo `json:"address"`
		Documents  []entity.MerchantDocument     `json:"documents"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   &errors.Error{Type: "INVALID_REQUEST", Message: "Invalid request body"},
		}, http.StatusBadRequest)
		return
	}

	// Verify merchant exists and belongs to user
	merchant, err := controller.interactor.GetMerchantByUserID(request.MerchantID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   &errors.Error{Type: "NOT_FOUND", Message: "Merchant not found"},
		}, http.StatusNotFound)
		return
	}

	if merchant.UserID != session.User.Id {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   &errors.Error{Type: "FORBIDDEN", Message: "Can only update your own merchants"},
		}, http.StatusForbidden)
		return
	}

	// Perform update (preserve existing status)
	request.Merchant.Status = merchant.Status
	err = controller.interactor.UpdateFullMerchant(r.Context(), request.MerchantID,
		&request.Merchant, &request.Address, request.Documents)
	if err != nil {
		SendJSONResponse(w, Response{Success: false, Error: err}, errors.MapErrorToHTTPStatus(err))
		return
	}

	SendJSONResponse(w, MerchantsResponse{
		Success: true,
		Message: "Merchant updated successfully",
	}, http.StatusOK)
}

func (controller Controller) GetMerchantBusinessInformations(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[MERCHANT] [CONTROLLER] [UPDATE] ")

	// Authentication
	token, err := controller.extractToken(r)
	if err != nil {
		SendJSONResponse(w, Response{Success: false, Error: err}, errors.MapErrorToHTTPStatus(err))
		return
	}

	// Authorization
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   &errors.Error{Type: "UNAUTHORIZED", Message: "Not authorized"},
		}, http.StatusUnauthorized)
		return
	}

	// Get the merchant details using the user ID
	merchantDetails, err := controller.interactor.GetMerchantDetails(session.User.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   &errors.Error{Type: "DATABASE_ERROR", Message: "Failed to retrieve merchant details"},
		}, http.StatusInternalServerError)
		return
	}

	// Check if merchant details is nil (no merchant profile exists)
	if merchantDetails == nil {
		SendJSONResponse(w, MerchantsResponse{
			Success: true,
			Message: "No merchant profile found for this user",
			Data:    nil,
		}, http.StatusOK)
		return
	}

	SendJSONResponse(w, MerchantsResponse{
		Success: true,
		Message: "Merchant detail data fetched successfully",
		Data:    merchantDetails,
	}, http.StatusOK)
}

func (controller Controller) DeleteMerchant(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[MERCHANT] [CONTROLLER] [DELETE] ")

	// Authentication
	token, err := controller.extractToken(r)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	// Authorization - check if user has permission
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Println("Authorization failed:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "UNAUTHORIZED",
				Message: "Not authorized to delete merchants",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Parse request
	var request struct {
		MerchantID string `json:"merchant_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate merchant ID
	merchantID, err := uuid.Parse(request.MerchantID)
	if err != nil {
		SendJSONResponse(w, MerchantsResponse{
			Success: false,
			Error: &errors.Error{
				Type:    "INVALID_INPUT",
				Message: "Invalid merchant ID format",
			},
		}, http.StatusBadRequest)
		return
	}

	// Verify the user has permission to delete this specific merchant
	// Here we will handle this with our RBAC in future
	merchant, err := controller.interactor.GetMerchantByUserID(merchantID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "NOT_FOUND",
				Message: "Merchant not found",
			},
		}, http.StatusNotFound)
		return
	}

	if merchant.UserID != session.User.Id {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &errors.Error{
				Type:    "FORBIDDEN",
				Message: "You can only delete your own information",
			},
		}, http.StatusForbidden)
		return
	}

	// Perform deletion
	err = controller.interactor.DeleteMerchant(r.Context(), merchantID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, errors.MapErrorToHTTPStatus(err))
		return
	}

	// Success response
	SendJSONResponse(w, MerchantsResponse{
		Success: true,
		Message: "Merchant deleted successfully",
	}, http.StatusOK)
}
