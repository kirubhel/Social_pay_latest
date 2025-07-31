package entity

import "github.com/google/uuid"

type User struct {
	Id uuid.UUID
}

type User2 struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	SirName     string    `json:"sir_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	UserType    string    `json:"user_type"`
}

type UserProfileUpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	SirName     string `json:"sir_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
