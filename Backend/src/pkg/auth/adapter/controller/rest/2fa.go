package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/socialpay/socialpay/src/pkg/auth/usecase"
	"github.com/socialpay/socialpay/src/pkg/utils"

	"github.com/google/uuid"
)

func (controller Controller) GetSet2FA(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
		Hint     string `json:"hint"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Usecase
	// Creare 2FA
	_, err = controller.interactor.InitPasswordAuth(req.Token, req.Password, req.Hint)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	//
	// Usecase

	// Create session
	log.Println("requesting for session creation")
	session, at, err := controller.interactor.CreateSession(req.Token)
	if err != nil {
		if err, ok := err.(usecase.Error); ok {
			switch err.Type {
			case "SET_PASSWORD":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "SET_PASSWORD",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			case "CHECK_PASSWORD":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "CHECK_PASSWORD",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			case "SIGN_UP":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "SIGN_UP",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			}
		}
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNSPECIFIED",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Return Response
	SendJSONResponse(w, Response{
		Success: true,
		Data: AuthResponse{
			Token: &struct {
				Active  string "json:\"active\""
				Refresh string "json:\"refresh\""
			}{
				Active:  at,
				Refresh: session.Token,
			},
			User: &struct {
				Id        uuid.UUID "json:\"id\""
				SirName   string    "json:\"sir_name,omitempty\""
				FirstName string    "json:\"first_name\""
				LastName  string    "json:\"last_name,omitempty\""
				UserType  string    "json:\"user_type,omitempty\""
			}{
				Id:        session.User.Id,
				SirName:   session.User.SirName,
				FirstName: session.User.FirstName,
				LastName:  session.User.LastName,
				UserType:  session.User.UserType,
			},
		},
	}, http.StatusOK)
}

func (controller Controller) GetCheck2FA(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Usecase
	// Check Password
	err = controller.interactor.AuthPassword(req.Token, req.Password)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Create session
	log.Println("requesting for session creation")
	session, at, err := controller.interactor.CreateSession(req.Token)
	if err != nil {
		if err, ok := err.(usecase.Error); ok {
			log.Println(err.Type)
			log.Println(err.Message)
			switch err.Type {
			case "SET_PASSWORD":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "SET_PASSWORD",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			case "CHECK_PASSWORD":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "CHECK_PASSWORD",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			case "SIGN_UP":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "SIGN_UP",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			}
		}
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNSPECIFIED",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Return Response
	SendJSONResponse(w, Response{
		Success: true,
		Data: AuthResponse{
			Token: &struct {
				Active  string "json:\"active\""
				Refresh string "json:\"refresh\""
			}{
				Active:  at,
				Refresh: session.Token,
			},
			User: &struct {
				Id        uuid.UUID "json:\"id\""
				SirName   string    "json:\"sir_name,omitempty\""
				FirstName string    "json:\"first_name\""
				LastName  string    "json:\"last_name,omitempty\""
				UserType  string    "json:\"user_type,omitempty\""
			}{
				Id:        session.User.Id,
				SirName:   session.User.SirName,
				FirstName: session.User.FirstName,
				LastName:  session.User.LastName,
				UserType:  session.User.UserType,
			},
		},
	}, http.StatusOK)
}

func (controller Controller) GetDecryptData(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Data string `json:"data"`
	}
	fmt.Println("###################################################")
	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	fmt.Println("################################################### , two")

	decryptedData, err := utils.AesDecription(req.Data)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return

	}

	parts := strings.Split(decryptedData, ",")

	merchantIdString := parts[0]

	parsedUUID, err := uuid.Parse(merchantIdString)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}
	fmt.Println("################################################### , three", parsedUUID)

	user, err := controller.interactor.GetUserById(parsedUUID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}
	fmt.Println("################################################### , four")

	type Response2 struct {
		Data     string `json:"data"`
		FullName string `json:"full_name"`
	}

	var res Response2

	res.FullName = user.SirName + " " + user.FirstName + " " + user.LastName
	res.Data = decryptedData

	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

// 2FA Management endpoints

// Get2FAStatus returns the current 2FA status for the authenticated user
func (controller Controller) Get2FAStatus(w http.ResponseWriter, r *http.Request) {
	// Get user ID from session token
	token := r.Header.Get("Authorization")
	if token == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNAUTHORIZED",
				Message: "Authorization token required",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	// Get session to find user ID
	session, err := controller.interactor.CheckSession(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_TOKEN",
				Message: "Invalid or expired token",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Get 2FA status from database
	status, err := controller.interactor.GetTwoFactorStatus(session.User.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "DATABASE_ERROR",
				Message: "Failed to get 2FA status",
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    status,
	}, http.StatusOK)
}

// Enable2FA enables 2FA for the authenticated user
func (controller Controller) Enable2FA(w http.ResponseWriter, r *http.Request) {
	// Get user ID from session token
	token := r.Header.Get("Authorization")
	if token == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNAUTHORIZED",
				Message: "Authorization token required",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	// Get session to find user ID
	session, err := controller.interactor.CheckSession(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_TOKEN",
				Message: "Invalid or expired token",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Get user's phone number
	user, err := controller.interactor.GetUserWithPhoneById(session.User.Id)
	if err != nil {
		controller.log.Printf("Failed to get user with phone for user ID %s: %v", session.User.Id, err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "USER_NOT_FOUND",
				Message: "User not found",
			},
		}, http.StatusNotFound)
		return
	}

	// Debug logging
	controller.log.Printf("User phone info - Prefix: '%s', Number: '%s', PhoneID: %s", user.Phone.Prefix, user.Phone.Number, user.PhoneID)

	// Check if user has a phone number
	if user.Phone.Prefix == "" || user.Phone.Number == "" {
		controller.log.Printf("Phone number is empty for user %s", session.User.Id)
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "PHONE_NOT_FOUND",
				Message: "Phone number not found for user",
			},
		}, http.StatusNotFound)
		return
	}

	// Format phone number for SMS
	phoneNumber := fmt.Sprintf("+%s%s", user.Phone.Prefix, user.Phone.Number)
	controller.log.Printf("Formatted phone number: %s", phoneNumber)

	// Enable 2FA and send verification code
	err = controller.interactor.EnableTwoFactor(session.User.Id, phoneNumber)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "ENABLE_2FA_FAILED",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "2FA setup initiated. Verification code sent to your phone.",
		},
	}, http.StatusOK)
}

// Disable2FA disables 2FA for the authenticated user
func (controller Controller) Disable2FA(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Password string `json:"password"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_PASSWORD",
				Message: "Password is required",
			},
		}, http.StatusBadRequest)
		return
	}

	// Get user ID from session token
	token := r.Header.Get("Authorization")
	if token == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNAUTHORIZED",
				Message: "Authorization token required",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	// Get session to find user ID
	session, err := controller.interactor.CheckSession(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_TOKEN",
				Message: "Invalid or expired token",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Disable 2FA
	err = controller.interactor.DisableTwoFactor(session.User.Id, req.Password)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "DISABLE_2FA_FAILED",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "2FA has been disabled successfully",
		},
	}, http.StatusOK)
}

// Verify2FASetup verifies the 2FA setup with the provided code
func (controller Controller) Verify2FASetup(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Code string `json:"code"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	if req.Code == "" || len(req.Code) != 6 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_CODE",
				Message: "Please enter a valid 6-digit verification code",
			},
		}, http.StatusBadRequest)
		return
	}

	// Get user ID from session token
	token := r.Header.Get("Authorization")
	if token == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNAUTHORIZED",
				Message: "Authorization token required",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	// Get session to find user ID
	session, err := controller.interactor.CheckSession(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_TOKEN",
				Message: "Invalid or expired token",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Verify the 2FA code
	err = controller.interactor.VerifyTwoFactorCode(session.User.Id, req.Code)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "VERIFICATION_FAILED",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "2FA has been enabled successfully",
		},
	}, http.StatusOK)
}

// Verify2FALogin verifies 2FA code during login
func (controller Controller) Verify2FALogin(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Token string `json:"token"`
		Code  string `json:"code"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	if req.Code == "" || len(req.Code) != 6 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_CODE",
				Message: "Please enter a valid 6-digit verification code",
			},
		}, http.StatusBadRequest)
		return
	}

	// Get session to find user ID
	session, err := controller.interactor.CheckSession(req.Token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_TOKEN",
				Message: "Invalid or expired token",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Verify the 2FA code
	err = controller.interactor.VerifyTwoFactorCode(session.User.Id, req.Code)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "VERIFICATION_FAILED",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Create session after successful 2FA verification
	session, at, err := controller.interactor.CreateSession(req.Token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "FAILED_TO_CREATE_SESSION",
				Message: "Failed to create session after 2FA verification",
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: map[string]interface{}{
			"token": map[string]string{"active": at, "refresh": session.Token},
			"user": map[string]interface{}{
				"id":         session.User.Id,
				"first_name": session.User.FirstName,
				"last_name":  session.User.LastName,
				"user_type":  session.User.UserType,
			},
		},
	}, http.StatusOK)
}

// Resend2FACode sends a new verification code for 2FA
func (controller Controller) Resend2FACode(w http.ResponseWriter, r *http.Request) {
	// Get user ID from session token
	token := r.Header.Get("Authorization")
	if token == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNAUTHORIZED",
				Message: "Authorization token required",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	// Get session to find user ID
	session, err := controller.interactor.CheckSession(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_TOKEN",
				Message: "Invalid or expired token",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Get user's phone number
	user, err := controller.interactor.GetUserWithPhoneById(session.User.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "USER_NOT_FOUND",
				Message: "User not found",
			},
		}, http.StatusNotFound)
		return
	}

	// Check if user has a phone number
	if user.Phone.Prefix == "" || user.Phone.Number == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "PHONE_NOT_FOUND",
				Message: "Phone number not found for user",
			},
		}, http.StatusNotFound)
		return
	}

	// Format phone number for SMS
	phoneNumber := fmt.Sprintf("+%s%s", user.Phone.Prefix, user.Phone.Number)

	// Send new verification code
	err = controller.interactor.SendTwoFactorCode(session.User.Id, phoneNumber)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "SEND_CODE_FAILED",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "New verification code sent to your phone",
		},
	}, http.StatusOK)
}
