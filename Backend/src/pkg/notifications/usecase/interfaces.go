package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/notifications/core/entity"
)

// SMSProvider defines the interface for SMS providers
type SMSProvider interface {
	SendSMS(phoneNumber, message string) error
}

// EmailProvider defines the interface for email providers (for future implementation)
type EmailProvider interface {
	SendEmail(to, subject, body string) error
}

// InAppProvider defines the interface for in-app notifications (for future implementation)
type InAppProvider interface {
	SendInApp(userID uuid.UUID, title, message string, data map[string]interface{}) error
}

// NotificationService defines the main notification service interface
type NotificationService interface {
	// Send a simple SMS notification
	SendSMS(ctx context.Context, phoneNumber, message string) error

	// Send a transaction notification to a phone number
	SendTransactionNotification(ctx context.Context, req TransactionNotificationRequest) error

	// Send notification using template
	SendTemplatedNotification(ctx context.Context, req TemplatedNotificationRequest) error

	// Get notification status
	GetNotificationStatus(ctx context.Context, notificationID uuid.UUID) (*entity.Notification, error)
}

// TransactionNotificationRequest represents a request to send transaction notification
type TransactionNotificationRequest struct {
	TransactionID   string                             `json:"transaction_id"`
	Type            entity.NotificationType            `json:"type"`
	Recipients      []NotificationRecipient            `json:"recipients"`
	TransactionData entity.TransactionNotificationData `json:"transaction_data"`
}

// NotificationRecipient represents a notification recipient
type NotificationRecipient struct {
	Type       entity.NotificationType `json:"type"`       // SMS, EMAIL, IN_APP
	Identifier string                  `json:"identifier"` // phone number, email, or user ID
	Name       string                  `json:"name"`       // recipient name
	Role       string                  `json:"role"`       // payer, merchant, tipee
}

// TemplatedNotificationRequest represents a request to send templated notification
type TemplatedNotificationRequest struct {
	Template  entity.NotificationTemplate `json:"template"`
	Type      entity.NotificationType     `json:"type"`
	Recipient string                      `json:"recipient"`
	Subject   string                      `json:"subject,omitempty"`
	Data      map[string]interface{}      `json:"data,omitempty"`
}

// Error represents a notification error
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

// Error constants
const (
	ErrInvalidPhoneNumber   = "INVALID_PHONE_NUMBER"
	ErrInvalidTemplate      = "INVALID_TEMPLATE"
	ErrProviderFailure      = "PROVIDER_FAILURE"
	ErrNotificationNotFound = "NOTIFICATION_NOT_FOUND"
)
