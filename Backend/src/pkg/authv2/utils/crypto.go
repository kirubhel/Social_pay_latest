package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/socialpay/socialpay/src/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash checks if a password matches a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateOTPCode generates a 6-digit OTP code
func GenerateOTPCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += digit.String()
	}
	return code, nil
}

// GenerateRandomToken generates a random token
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// JWTClaims represents JWT claims for auth v2
type JWTClaims struct {
	UserID     string `json:"user_id"`
	UserType   string `json:"user_type"`
	MerchantID string `json:"merchant_id,omitempty"`
	SessionID  string `json:"session_id"`
}

// GenerateJWT generates a JWT token using the existing socialpay JWT implementation
func GenerateJWT(userID, userType, merchantID, sessionID, secret string, expiryHours int) (string, error) {
	payload := jwt.Payload{
		Public: JWTClaims{
			UserID:     userID,
			UserType:   userType,
			MerchantID: merchantID,
			SessionID:  sessionID,
		},
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix(),
	}

	return jwt.Encode(payload, secret), nil
}

// ValidateJWT validates a JWT token and returns claims
func ValidateJWT(tokenString, secret string) (*JWTClaims, error) {
	payload, err := jwt.Decode(tokenString, secret)
	if err != nil {
		return nil, err
	}

	// Extract claims from public field
	if payload.Public == nil {
		return nil, fmt.Errorf("invalid token: missing claims")
	}

	// Convert public field to claims
	claimsMap, ok := payload.Public.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid token: claims format error")
	}

	claims := &JWTClaims{}
	if userID, ok := claimsMap["user_id"].(string); ok {
		claims.UserID = userID
	}
	if userType, ok := claimsMap["user_type"].(string); ok {
		claims.UserType = userType
	}
	if merchantID, ok := claimsMap["merchant_id"].(string); ok {
		claims.MerchantID = merchantID
	}
	if sessionID, ok := claimsMap["session_id"].(string); ok {
		claims.SessionID = sessionID
	}

	return claims, nil
}

// GenerateOTPToken generates a secure token for OTP verification
func GenerateOTPToken() (string, error) {
	return GenerateRandomToken(32)
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken() (string, error) {
	return GenerateRandomToken(64)
}
