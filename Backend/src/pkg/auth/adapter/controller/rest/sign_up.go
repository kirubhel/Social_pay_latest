package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"
	"github.com/socialpay/socialpay/src/pkg/jwt"

	"github.com/google/uuid"
)

func (controller Controller) GetSignUp(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Title           string `json:"title"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		PhonePrefix     string `json:"phone_prefix"`
		PhoneNumber     string `json:"phone_number"`
		Password        string `json:"password"`
		PasswordHint    string `json:"password_hint,omitempty"`
		ConfirmPassword string `json:"confirm_password"`
	}

	type PhoneResponse struct {
		Prefix  string    `json:"prefix"`
		Number  string    `json:"number"`
		PhoneID uuid.UUID `json:"phone_id"`
	}

	type UserResponse struct {
		ID        uuid.UUID     `json:"id"`
		Title     string        `json:"title"`
		FirstName string        `json:"first_name"`
		LastName  string        `json:"last_name"`
		UserType  string        `json:"user_type"`
		Phone     PhoneResponse `json:"phone"`
		Auth      struct {
			Token string `json:"token"`
		} `json:"auth"`
		CreatedAt time.Time `json:"created_at"`
	}

	type Response struct {
		Success bool          `json:"success"`
		Message string        `json:"message,omitempty"`
		Data    *UserResponse `json:"data,omitempty"`
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

	// Validate phone number format
	if !isValidEthiopianPhoneNumber(req.PhoneNumber) {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    "INVALID_PHONE_NUMBER_FORMAT",
				Message: "Invalid phone number. Ethiopian numbers must start with 9 and be 9 digits long (e.g., 911234567)",
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate password match
	if req.Password != req.ConfirmPassword {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    entity.ErrPasswordMismatch,
				Message: entity.MsgPasswordMismatch,
			},
		}, http.StatusBadRequest)
		return
	}

	// Generate pre-session token
	preSessionID := uuid.New()
	token := jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + 1800,
		Public: preSessionID,
	}, "pre_session_secret")

	preSession := entity.PreSession{
		Id:        preSessionID,
		Token:     token,
		CreatedAt: time.Now(),
	}

	// Store pre-session
	if err := controller.repo.StorePreSession(preSession); err != nil {
		controller.log.Printf("Failed to store pre-session: %v", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    "FAILED_TO_CREATE_PRE_SESSION",
				Message: "Failed to create authentication session",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Process user creation
	user, userErrors := controller.interactor.CreateUser(
		req.Title,
		req.FirstName,
		req.LastName,
		req.PhonePrefix,
		req.PhoneNumber,
		req.Password,
		req.PasswordHint,
		"merchant",
	)

	if userErrors != nil {
		var apiErr *entity.Error
		if errors.As(userErrors, &apiErr) {
			SendJSONResponse(w, Response{
				Success: false,
				Error:   apiErr,
			}, http.StatusBadRequest)
			return
		}

		SendJSONResponse(w, Response{
			Success: false,
			Error: &entity.Error{
				Type:    entity.ErrInternalServer,
				Message: entity.MsgInternalServer,
			},
		}, http.StatusInternalServerError)
		return
	}

	// Generate and send OTP
	phoneAuth := entity.PhoneAuth{
		Id:      uuid.New(),
		Token:   token,
		Phone:   entity.Phone{Id: user.Phone.Id},
		Method:  "SMS",
		Length:  6,
		Timeout: 120,
	}

	otp := rand.Intn(999999-100000) + 100000
	otpStr := fmt.Sprint(otp)
	phoneAuth.Code = jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + 30*60,
		Public: otpStr,
	}, "otp_verification_secret")

	if err := controller.repo.StorePhoneAuth(phoneAuth); err != nil {
		controller.log.Printf("Failed to store phone auth: %v", err)
	}

	go func() {
		fullPhone := user.Phone.String()
		message := fmt.Sprintf("Your SocialPay verification code is %s. Do not share this code with anyone.", strconv.Itoa(otp))

		if err := controller.sms.SendSMS(fullPhone, message); err != nil {
			controller.log.Printf("Failed to send OTP SMS: %v", err)
		} else {
			controller.log.Printf("OTP sent successfully to %s", fullPhone)
		}
	}()

	userResponse := UserResponse{
		ID:        user.Id,
		Title:     user.SirName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserType:  user.UserType,
		Phone: PhoneResponse{
			Prefix:  user.Phone.Prefix,
			Number:  user.Phone.Number,
			PhoneID: user.Phone.Id,
		},
		Auth: struct {
			Token string `json:"token"`
		}{
			Token: token,
		},
		CreatedAt: user.CreatedAt,
	}

	SendJSONResponse(w, Response{
		Success: true,
		Message: "User registered successfully. OTP verification required.",
		Data:    &userResponse,
	}, http.StatusCreated)
}

func isValidEthiopianPhoneNumber(number string) bool {
	if len(number) != 9 {
		return false
	}

	if number[0] != '9' {
		return false
	}

	for _, c := range number {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}
