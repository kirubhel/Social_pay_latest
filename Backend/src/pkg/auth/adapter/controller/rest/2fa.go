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
	// For now, return a default response
	// TODO: Implement proper 2FA status checking from database
	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Enabled bool `json:"enabled"`
		}{
			Enabled: false,
		},
	}, http.StatusOK)
}

// Enable2FA enables 2FA for the authenticated user
func (controller Controller) Enable2FA(w http.ResponseWriter, r *http.Request) {
	// For now, simulate enabling 2FA by sending an OTP
	// TODO: Implement proper 2FA enabling logic

	// Generate a simple token for demo purposes
	token := "demo-2fa-token-" + uuid.New().String()

	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Message string `json:"message"`
			Token   string `json:"token"`
		}{
			Message: "2FA setup initiated. Verification code sent to your phone.",
			Token:   token,
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

	// TODO: Implement proper password verification and 2FA disabling
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

	// TODO: Implement proper OTP verification logic
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

	// For demo purposes, accept any 6-digit code
	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Message string `json:"message"`
		}{
			Message: "2FA has been enabled successfully",
		},
	}, http.StatusOK)
}
