package entity

import (
	"time"

	"github.com/google/uuid"

	rbac_entity "github.com/socialpay/socialpay/src/pkg/rbac/core/entity"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	SirName     string     `json:"sir_name" db:"sir_name"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	Email       string     `json:"email" db:"email"`
	UserType    UserType   `json:"user_type" db:"user_type"`
	Gender      Gender     `json:"gender" db:"gender"`
	DateOfBirth time.Time  `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	PhonePrefix string     `json:"phone_prefix" db:"phone_prefix"`
	PhoneNumber string     `json:"phone_number" db:"phone_number"`
	Permissions []string   `json:"permissions,omitempty"`
	Groups      []Group    `json:"groups,omitempty"`
	MerchantID  *uuid.UUID `json:"merchant_id,omitempty"`
}

// Phone represents a phone number
type Phone struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Prefix    string    `json:"prefix" db:"prefix"`
	Number    string    `json:"number" db:"number"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (p Phone) String() string {
	return "+" + p.Prefix + p.Number
}

// Group represents a user group/role
type Group struct {
	ID          uuid.UUID                `json:"id" db:"id"`
	Title       string                   `json:"title" db:"title"`
	Description *string                  `json:"description" db:"description"`
	MerchantID  *uuid.UUID               `json:"merchant_id" db:"merchant_id"`
	Permissions []rbac_entity.Permission `json:"permissions,omitempty"`
	CreatedAt   time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at" db:"updated_at"`
}

// Permission represents a permission
type Permission struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	ResourceName string      `json:"resource_name" db:"resource_name"`
	Operations   []uuid.UUID `json:"operations" db:"operations"`
	Effect       string      `json:"effect" db:"effect"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

// AuthActivity represents authentication activity logs
type AuthActivity struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	UserID       uuid.UUID              `json:"user_id" db:"user_id"`
	ActivityType AuthActivityType       `json:"activity_type" db:"activity_type"`
	IPAddress    string                 `json:"ip_address" db:"ip_address"`
	UserAgent    string                 `json:"user_agent" db:"user_agent"`
	DeviceName   string                 `json:"device_name" db:"device_name"`
	Success      bool                   `json:"success" db:"success"`
	Details      map[string]interface{} `json:"details" db:"details"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

// Session represents an active user session
type Session struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	DeviceID     uuid.UUID `json:"device_id" db:"device_id"`
	Token        string    `json:"token" db:"token"`
	RefreshToken string    `json:"refresh_token,omitempty" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Device represents a user device
type Device struct {
	ID        uuid.UUID `json:"id" db:"id"`
	IP        string    `json:"ip" db:"ip"`
	Name      string    `json:"name" db:"name"`
	Agent     string    `json:"agent" db:"agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// OTPRequest represents an OTP verification request
type OTPRequest struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Code      string    `json:"code" db:"code"`
	Token     string    `json:"token" db:"token"`
	Method    string    `json:"method" db:"method"`
	Status    bool      `json:"status" db:"status"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// DeviceInfo represents device information
type DeviceInfo struct {
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`
	DeviceName string `json:"device_name"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Title        string      `json:"title" validate:"required"`
	FirstName    string      `json:"first_name" validate:"required"`
	LastName     string      `json:"last_name" validate:"required"`
	Email        string      `json:"email" validate:"required,email"`
	PhonePrefix  string      `json:"phone_prefix" validate:"required"`
	PhoneNumber  string      `json:"phone_number" validate:"required"`
	Password     string      `json:"password" validate:"required,min=8"`
	PasswordHint string      `json:"password_hint"`
	UserType     UserType    `json:"user_type" validate:"required"`
	DeviceInfo   *DeviceInfo `json:"device_info,omitempty"`
}

// UpdateUserRequest represents a request to update user information
// Excludes user_type and device_info as per requirements
type UpdateUserRequest struct {
	Title        *string `json:"title,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Email        *string `json:"email,omitempty" validate:"omitempty,email"`
	PhonePrefix  *string `json:"phone_prefix,omitempty"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
	Password     *string `json:"password,omitempty" validate:"omitempty,min=8"`
	PasswordHint *string `json:"password_hint,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	PhonePrefix string      `json:"phone_prefix" validate:"required"`
	PhoneNumber string      `json:"phone_number" validate:"required"`
	Password    string      `json:"password" validate:"required"`
	DeviceInfo  *DeviceInfo `json:"device_info,omitempty"`
}

// VerifyOTPRequest represents an OTP verification request
type VerifyOTPRequest struct {
	Token      string      `json:"token" validate:"required"`
	Code       string      `json:"code" validate:"required"`
	DeviceInfo *DeviceInfo `json:"device_info,omitempty"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *User      `json:"user"`
	Merchants    []Merchant `json:"merchants,omitempty"`
	Token        string     `json:"token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresAt    int64      `json:"expires_at"`
}

// Merchant represents a merchant in the system
type Merchant struct {
	ID                         uuid.UUID `json:"id" db:"id"`
	UserID                     uuid.UUID `json:"user_id" db:"user_id"`
	LegalName                  string    `json:"legal_name" db:"legal_name"`
	TradingName                string    `json:"trading_name" db:"trading_name"`
	BusinessRegistrationNumber string    `json:"business_registration_number" db:"business_registration_number"`
	TaxIdentificationNumber    string    `json:"tax_identification_number" db:"tax_identification_number"`
	IndustryCategory           string    `json:"industry_category" db:"industry_category"`
	BusinessType               string    `json:"business_type" db:"business_type"`
	IsBettingCompany           bool      `json:"is_betting_company" db:"is_betting_company"`
	LotteryCertificateNumber   string    `json:"lottery_certificate_number" db:"lottery_certificate_number"`
	WebsiteURL                 string    `json:"website_url" db:"website_url"`
	EstablishedDate            time.Time `json:"established_date" db:"established_date"`
	Status                     string    `json:"status" db:"status"`
	CreatedAt                  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateDeviceArgs represents args for creating devices
type CreateDeviceArgs struct {
	Agent string `json:"agent"`
	IP    string `json:"ip"`
	Name  string `json:"name"`
}

// CreateSessionArgs represents args for creating session
type CreateSessionArgs struct {
	DeviceID     uuid.UUID `json:"device_id"`
	ExpiresAt    int64     `json:"expires_at"`
	RefreshToken string    `json:"refresh_token"`
	Token        string    `json:"token"`
	UserID       uuid.UUID `json:"user_id"`
}

// RequestOTPRequest represents args for requesting otp
type RequestOTPRequest struct {
	PhonePrefix string `json:"phone_prefix" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

// UpdatePasswordRequest represents args for updating password
type UpdatePasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required"`
	OTPToken    string `json:"otp_token" validate:"required"`
}
