package rest

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/auth/usecase"
)

const ErrFailedToInitPhoneAuth = "FAILED_TO_INIT_PHONE_AUTH"

func (controller Controller) GetInitAuth(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [GetInitAuth] ")
	// Accept req with
	// Device info
	// 		- IP
	// 		- Name
	// 		- Agent
	// Phone
	//		- Code
	// 		- Number

	// Request
	type Request struct {
		Device struct {
			IP    net.IPAddr `json:"ip"`
			Name  string     `json:"name"`
			Agent string     `json:"agent"`
		} `json:"device"`
		Phone struct {
			Prefix string `json:"prefix"`
			Number string `json:"number"`
		} `json:"phone"`
	}

	// Response
	type Response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
		Error   error       `json:"error,omitempty"`
	}

	// Parse Request
	var req Request

	req.Device = struct {
		IP    net.IPAddr "json:\"ip\""
		Name  string     "json:\"name\""
		Agent string     "json:\"agent\""
	}{
		IP: net.IPAddr{
			IP: net.ParseIP(r.RemoteAddr),
		},
		Name:  "name",
		Agent: "agent",
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req.Phone)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Do usecase operations
	// Initiate Authentication
	controller.log.Println("initiating authentication")
	preSession, err := controller.interactor.InitPreSession()
	if err != nil {
		controller.log.Println("failed to initiate authentication")
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Init Phone Authentication
	controller.log.Println("authenticating phone")
	phoneAuth, err := controller.interactor.InitPhoneAuth(preSession.Token, req.Phone.Prefix, req.Phone.Number)
	if err != nil {
		controller.log.Printf("failed authenticating phone : %s\n", err.Error())
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Return Response
	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Token   string `json:"token"`
			Method  string `json:"method"`
			Length  int64  `json:"length"`
			Timeout int64  `json:"timeout"`
		}{
			Token:   preSession.Token,
			Method:  phoneAuth.Method,
			Length:  phoneAuth.Length,
			Timeout: phoneAuth.Timeout,
		},
	}, http.StatusOK)
}

func (controller Controller) GetInitLogin(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [GetInitAuth] ")
	// Accept req with
	// Device info
	// 		- IP
	// 		- Name
	// 		- Agent
	// Phone
	//		- Code
	// 		- Number

	// Request
	type Request struct {
		Device struct {
			IP    net.IPAddr `json:"ip"`
			Name  string     `json:"name"`
			Agent string     `json:"agent"`
		} `json:"device"`
		Phone struct {
			Prefix string `json:"prefix"`
			Number string `json:"number"`
		} `json:"phone"`
	}

	// Response
	type Response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
		Error   error       `json:"error,omitempty"`
	}

	// Parse Request
	var req Request

	req.Device = struct {
		IP    net.IPAddr "json:\"ip\""
		Name  string     "json:\"name\""
		Agent string     "json:\"agent\""
	}{
		IP: net.IPAddr{
			IP: net.ParseIP(r.RemoteAddr),
		},
		Name:  "name",
		Agent: "agent",
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req.Phone)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	phone, err := controller.interactor.LoginFindPhone(req.Phone.Prefix, req.Phone.Number)
	if err != nil {
		controller.log.Println("failed finding phone:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    ErrFailedToInitPhoneAuth,
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	if phone == nil {
		controller.log.Println("phone not found")
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "PHONE_NOT_FOUND",
				Message: "The phone number does not exist.",
			},
		}, http.StatusNotFound)
		return
	}

	// Only if the phone is found, proceed to initiate authentication
	controller.log.Println("initiating authentication")
	preSession, err := controller.interactor.InitPreSession()
	if err != nil {
		controller.log.Println("failed to initiate authentication")
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Authenticate Device
	controller.log.Println("authenticating device")
	err = controller.interactor.AuthDevice(preSession.Token, req.Device.IP, req.Device.Name, req.Device.Agent)
	if err != nil {
		controller.log.Printf("failed authenticating device %s\n", err.Error())
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Init Phone Authentication
	controller.log.Println("authenticating phone")
	phoneAuth, err := controller.interactor.InitPhoneAuth(preSession.Token, req.Phone.Prefix, req.Phone.Number)
	if err != nil {
		controller.log.Printf("failed authenticating phone %s\n", err.Error())
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Return Response
	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Token   string `json:"token"`
			Method  string `json:"method"`
			Length  int64  `json:"length"`
			Timeout int64  `json:"timeout"`
		}{
			Token:   preSession.Token,
			Method:  phoneAuth.Method,
			Length:  phoneAuth.Length,
			Timeout: phoneAuth.Timeout,
		},
	}, http.StatusOK)
}

// Login with phone and password, with optional OTP (2FA)
func (controller Controller) LoginWithPhoneAndPassword(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [LoginWithPhoneAndPassword] ")
	type Request struct {
		Prefix   string `json:"prefix"`
		Number   string `json:"number"`
		Password string `json:"password"`
	}
	type Response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
		Error   interface{} `json:"error,omitempty"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   map[string]string{"type": "INVALID_REQUEST", "message": err.Error()},
		}, http.StatusBadRequest)
		return
	}
	// 1. Find phone
	phone, err := controller.interactor.LoginFindPhone(req.Prefix, req.Number)
	if err != nil || phone == nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   map[string]string{"type": "PHONE_NOT_FOUND", "message": "Phone number not found."},
		}, http.StatusNotFound)
		return
	}
	// 2. Create pre-session
	preSession, err := controller.interactor.InitPreSession()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   map[string]string{"type": "FAILED_TO_CREATE_PRE_SESSION", "message": err.Error()},
		}, http.StatusInternalServerError)
		return
	}
	// 3. Init phone auth (for this session)
	_, err = controller.interactor.InitPhoneAuth(preSession.Token, req.Prefix, req.Number)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   map[string]string{"type": "FAILED_TO_INIT_PHONE_AUTH", "message": err.Error()},
		}, http.StatusInternalServerError)
		return
	}
	// 4. Check password
	err = controller.interactor.AuthPassword(preSession.Token, req.Password)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   map[string]string{"type": "INCORRECT_PASSWORD", "message": err.Error()},
		}, http.StatusUnauthorized)
		return
	}
	// 5. Check if 2FA is enabled for the user
	user, err := controller.interactor.GetUserWithPhoneById(phone.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   map[string]string{"type": "USER_NOT_FOUND", "message": "User not found"},
		}, http.StatusNotFound)
		return
	}

	// Get 2FA status
	twoFactorStatus, err := controller.interactor.GetTwoFactorStatus(user.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   map[string]string{"type": "FAILED_TO_GET_2FA_STATUS", "message": "Failed to check 2FA status"},
		}, http.StatusInternalServerError)
		return
	}

	// If 2FA is enabled, require 2FA verification
	if twoFactorStatus.Enabled {
		// Send 2FA code
		phoneNumber := fmt.Sprintf("+%s%s", user.Phone.Prefix, user.Phone.Number)
		err = controller.interactor.SendTwoFactorCode(user.Id, phoneNumber)
		if err != nil {
			SendJSONResponse(w, Response{
				Success: false,
				Error:   map[string]string{"type": "FAILED_TO_SEND_2FA", "message": "Failed to send 2FA code"},
			}, http.StatusInternalServerError)
			return
		}

		SendJSONResponse(w, Response{
			Success: true,
			Data: map[string]interface{}{
				"next_step": "2FA_REQUIRED",
				"token":     preSession.Token,
				"message":   "2FA verification required. Code sent to your phone.",
			},
		}, http.StatusAccepted)
		return
	}

	// If 2FA is not enabled, check if phone auth is verified
	err = controller.interactor.CheckPhoneAuth(preSession.Token)
	if err == nil {
		// Phone already verified, create session
		session, at, err := controller.interactor.CreateSession(preSession.Token)
		if err != nil {
			SendJSONResponse(w, Response{
				Success: false,
				Error:   map[string]string{"type": "FAILED_TO_CREATE_SESSION", "message": err.Error()},
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
		return
	}
	// If phone not verified, require OTP
	SendJSONResponse(w, Response{
		Success: true,
		Data: map[string]interface{}{
			"next_step": "OTP_REQUIRED",
			"token":     preSession.Token,
		},
	}, http.StatusAccepted)
}
