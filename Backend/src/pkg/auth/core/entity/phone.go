package entity

import (
	"time"

	"github.com/google/uuid"
)

type Phone struct {
	Id        uuid.UUID `json:"id"`
	PhoneID   uuid.UUID `json:"phone_id"`
	Prefix    string    `json:"prefix"`
	Number    string    `json:"number"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (phone Phone) String() string {
	return "+" + phone.Prefix + phone.Number
}

type PhoneAuth struct {
	Id        uuid.UUID
	Token     string
	Phone     Phone
	Code      string
	Status    bool
	Method    string
	Length    int64
	Timeout   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
