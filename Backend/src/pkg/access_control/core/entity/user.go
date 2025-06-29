package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"id"`
	SirName     string    `json:"sir_name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	UserType    string    `json:"user_type"`
	Gender      Gender    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Gender string

var (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
)
