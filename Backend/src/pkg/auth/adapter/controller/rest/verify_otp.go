package rest

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"
	"github.com/socialpay/socialpay/src/pkg/jwt"

	"github.com/google/uuid"
)

func (controller Controller) VerifySignUpOTP(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Token string `json:"token"`
		Code  string `json:"code"`
	}

	type AuthResponse struct {
		Token *struct {
			Active  string `json:"active"`
			Refresh string `json:"refresh"`
		} `json:"token,omitempty"`
		User *struct {
			Id        uuid.UUID `json:"id"`
			SirName   string    `json:"sir_name,omitempty"`
			FirstName string    `json:"first_name"`
			LastName  string    `json:"last_name,omitempty"`
			UserType  string    `json:"user_type,omitempty"`
		} `json:"user,omitempty"`
	}

	type Response struct {
		Success bool          `json:"success"`
		Message string        `json:"message,omitempty"`
		Data    *AuthResponse `json:"data,omitempty"`
		Error   *entity.Error `json:"error,omitempty"`
	}

	// Parse request
	var req Request
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    entity.ErrInvalidRequest,
				Message: entity.MsgInvalidRequest,
			},
		}, http.StatusBadRequest)
		return
	}

	// Find the phone auth record
	phoneAuth, err := controller.repo.FindPhoneAuth(req.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &entity.Error{
					Type:    "INVALID_OTP_REQUEST",
					Message: "No OTP verification request found for this session",
				},
			}, http.StatusBadRequest)
			return
		}

		controller.log.Printf("Failed to find phone auth: %v", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    entity.ErrInternalServer,
				Message: entity.MsgInternalServer,
			},
		}, http.StatusInternalServerError)
		return
	}

	if phoneAuth.Status {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    "OTP_ALREADY_VERIFIED",
				Message: "This OTP has already been verified",
			},
		}, http.StatusBadRequest)
		return
	}

	decodedPayload, err := jwt.Decode(phoneAuth.Code, "otp_verification_secret")
	if err != nil {
		controller.log.Printf("Failed to decode OTP: %v", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    "INVALID_OTP_FORMAT",
				Message: "Invalid OTP format",
			},
		}, http.StatusBadRequest)
		return
	}

	decodedCode, ok := decodedPayload.Public.(string)
	if !ok {
		controller.log.Printf("Failed to extract OTP code from payload")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    "INVALID_OTP_FORMAT",
				Message: "Invalid OTP format",
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Printf("OTP Verification - Stored: %s, Input: %s", decodedCode, req.Code)
	if decodedCode != req.Code {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    "INVALID_OTP",
				Message: "The OTP code you entered is incorrect",
			},
		}, http.StatusBadRequest)
		return
	}

	if time.Now().Unix() > decodedPayload.Exp {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    "OTP_EXPIRED",
				Message: "The OTP code has expired",
			},
		}, http.StatusBadRequest)
		return
	}

	// Update phone auth status to verified
	if err := controller.repo.UpdatePhoneAuthStatus(phoneAuth.Id, true); err != nil {
		controller.log.Printf("Failed to update phone auth status: %v", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    entity.ErrInternalServer,
				Message: entity.MsgInternalServer,
			},
		}, http.StatusInternalServerError)
		return
	}

	// Create session after successful OTP verification
	session, activeToken, err := controller.interactor.CreateSession(req.Token)
	if err != nil {
		controller.log.Printf("Failed to create session: %v", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    "SESSION_CREATION_FAILED",
				Message: "Failed to create user session",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Return success response with tokens and user details
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Phone number successfully verified and user session created",
		Data: &AuthResponse{
			Token: &struct {
				Active  string `json:"active"`
				Refresh string `json:"refresh"`
			}{
				Active:  activeToken,
				Refresh: session.Token,
			},
			User: &struct {
				Id        uuid.UUID `json:"id"`
				SirName   string    `json:"sir_name,omitempty"`
				FirstName string    `json:"first_name"`
				LastName  string    `json:"last_name,omitempty"`
				UserType  string    `json:"user_type,omitempty"`
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
