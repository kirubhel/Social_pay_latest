package entity

import (
	"time"

	"github.com/google/uuid"
)

type MerchantDocument struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	MerchantID      uuid.UUID  `db:"merchant_id" json:"merchant_id"`
	DocumentType    string     `db:"document_type" json:"document_type"`
	DocumentNumber  *string    `db:"document_number" json:"document_number,omitempty"`
	FileURL         string     `db:"file_url" json:"file_url"`
	VerifiedBy      *uuid.UUID `db:"verified_by" json:"verified_by,omitempty"`
	VerifiedAt      *time.Time `db:"verified_at" json:"verified_at,omitempty"`
	Status          string     `db:"status" json:"status"`
	RejectionReason *string    `db:"rejection_reason" json:"rejection_reason,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}
