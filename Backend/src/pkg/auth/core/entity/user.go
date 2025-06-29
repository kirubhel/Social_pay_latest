package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id            uuid.UUID
	SirName       string
	FirstName     string
	LastName      string
	UserType      string
	Gender        Gender
	DateOfBirth   time.Time
	Nationalities []Nationality
	Addresses     []Address
	Identities    []Identity
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Phone         Phone
	PhoneID       uuid.UUID
}

type Gender string

var (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
)

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Detail  string `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

const (
	ErrInvalidRequest           = "INVALID_REQUEST"
	ErrPasswordMismatch         = "PASSWORD_MISMATCH"
	ErrPhoneAlreadyExists       = "PHONE_ALREADY_EXISTS"
	ErrInvalidPhoneFormat       = "INVALID_PHONE_FORMAT"
	ErrMissingRequiredData      = "MISSING_REQUIRED_DATA"
	ErrAccountCreation          = "ACCOUNT_CREATION_FAILED"
	ErrInternalServer           = "INTERNAL_ERROR"
	ErrInvalidPhoneNumberFormat = "ErrInvalidPhoneNumberFormat"
	ErrInvalidPhoneNumberLength = "ErrInvalidPhoneNumberLength"
)

const (
	MsgInvalidRequest     = "Invalid request format"
	MsgPasswordMismatch   = "Password and confirmation do not match"
	MsgPhoneExists        = "This phone number is already registered"
	MsgInvalidPhoneFormat = "Phone prefix must be 3 digits or less (excluding + sign)"
	MsgMissingData        = "Required information was not provided"
	MsgAccountCreation    = "We couldn't create your account. Please try again later."
	MsgInternalServer     = "An unexpected error occurred"
	ErrInvalidPhoneNumber = "Invalid phone number must be between 9 digits long like +251911234567"
)
