package entity

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	TypeSMS   NotificationType = "SMS"
	TypeEmail NotificationType = "EMAIL"
	TypeInApp NotificationType = "IN_APP"
)

// NotificationStatus represents the delivery status
type NotificationStatus string

const (
	StatusPending   NotificationStatus = "PENDING"
	StatusSent      NotificationStatus = "SENT"
	StatusDelivered NotificationStatus = "DELIVERED"
	StatusFailed    NotificationStatus = "FAILED"
)

// NotificationTemplate represents different message templates
type NotificationTemplate string

const (
	TemplateTransactionSuccess NotificationTemplate = "TRANSACTION_SUCCESS"
	TemplateTransactionFailed  NotificationTemplate = "TRANSACTION_FAILED"
	TemplateTransactionPending NotificationTemplate = "TRANSACTION_PENDING"
	TemplateOTP                NotificationTemplate = "OTP"
	TemplateGeneral            NotificationTemplate = "GENERAL"
)

// Notification represents a notification to be sent
type Notification struct {
	ID        uuid.UUID              `json:"id"`
	Type      NotificationType       `json:"type"`
	Template  NotificationTemplate   `json:"template"`
	Recipient string                 `json:"recipient"` // phone number, email, or user ID
	Subject   string                 `json:"subject,omitempty"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"` // template data
	Status    NotificationStatus     `json:"status"`
	Attempts  int                    `json:"attempts"`
	LastError string                 `json:"last_error,omitempty"`
	SentAt    *time.Time             `json:"sent_at,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// SMSNotification represents SMS-specific notification data
type SMSNotification struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

// EmailNotification represents email-specific notification data
type EmailNotification struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// TransactionNotificationData represents data for transaction notifications
type TransactionNotificationData struct {
	TransactionID   string  `json:"transaction_id"`
	MerchantName    string  `json:"merchant_name"`
	CustomerName    string  `json:"customer_name"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	Reference       string  `json:"reference"`
	Timestamp       string  `json:"timestamp"`
	TransactionType string  `json:"transaction_type"`
	TipAmount       float64 `json:"tip_amount"`
	PhoneNumber     string  `json:"phone_number"`
}
