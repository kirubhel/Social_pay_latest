package rest

import (
	"encoding/json"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// UpdatePassword handles password update requests
func (controller Controller) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
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

	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_INPUT",
				Message: "Current password and new password are required",
			},
		}, http.StatusBadRequest)
		return
	}

	// Password policy validation
	if len(req.NewPassword) < 8 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_PASSWORD",
				Message: "New password must be at least 8 characters long",
			},
		}, http.StatusBadRequest)
		return
	}

	// Check for uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(req.NewPassword) {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_PASSWORD",
				Message: "Password must contain at least one uppercase letter (A-Z)",
			},
		}, http.StatusBadRequest)
		return
	}

	// Check for lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(req.NewPassword) {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_PASSWORD",
				Message: "Password must contain at least one lowercase letter (a-z)",
			},
		}, http.StatusBadRequest)
		return
	}

	// Check for number
	if !regexp.MustCompile(`\d`).MatchString(req.NewPassword) {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_PASSWORD",
				Message: "Password must contain at least one number (0-9)",
			},
		}, http.StatusBadRequest)
		return
	}

	// Check for special character
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(req.NewPassword) {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_PASSWORD",
				Message: "Password must contain at least one special character (!@#$%^&*(),.?\":{}|&lt;&gt;)",
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
	if len(token) > 7 && token[:7] == "Bearer " {
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

	// Get current password identity
	passwordIdentity, err := controller.repo.FindPasswordIdentityByUser(session.User.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "DATABASE_ERROR",
				Message: "Failed to retrieve password information",
			},
		}, http.StatusInternalServerError)
		return
	}

	if passwordIdentity == nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "NO_PASSWORD",
				Message: "No password found for user",
			},
		}, http.StatusNotFound)
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(passwordIdentity.Password), []byte(req.CurrentPassword))
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INCORRECT_PASSWORD",
				Message: "Current password is incorrect",
			},
		}, http.StatusBadRequest)
		return
	}

	// Check if new password is same as current
	err = bcrypt.CompareHashAndPassword([]byte(passwordIdentity.Password), []byte(req.NewPassword))
	if err == nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "SAME_PASSWORD",
				Message: "New password must be different from current password",
			},
		}, http.StatusBadRequest)
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "HASH_ERROR",
				Message: "Failed to hash new password",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Update password in database
	err = controller.repo.UpdatePasswordIdentity(string(hashedPassword), session.User.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "DATABASE_ERROR",
				Message: "Failed to update password",
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: map[string]interface{}{
			"message": "Password updated successfully",
		},
	}, http.StatusOK)
}
