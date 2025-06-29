package usecase

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/notifications/core/entity"
)

// NotificationServiceImpl implements the NotificationService interface
type NotificationServiceImpl struct {
	smsProvider   SMSProvider
	emailProvider EmailProvider
	inAppProvider InAppProvider
	log           *log.Logger
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(smsProvider SMSProvider, log *log.Logger) NotificationService {
	return &NotificationServiceImpl{
		smsProvider: smsProvider,
		log:         log,
	}
}

// SendSMS sends a simple SMS notification
func (ns *NotificationServiceImpl) SendSMS(ctx context.Context, phoneNumber, message string) error {
	ns.log.Printf("[SendSMS] Sending SMS to: %s", phoneNumber)

	// Validate and normalize phone number
	normalizedPhone, err := normalizeEthiopianPhoneNumber(phoneNumber)
	ns.log.Printf("[SendSMS] Normalized phone number: %v", normalizedPhone)
	if err != nil {
		ns.log.Printf("[SendSMS] Invalid phone number: %v", err)
		return Error{
			Type:    ErrInvalidPhoneNumber,
			Message: err.Error(),
		}
	}

	// Send SMS using provider
	if err := ns.smsProvider.SendSMS(normalizedPhone, message); err != nil {
		ns.log.Printf("[SendSMS] Failed to send SMS: normalizedPhone %v, message %v, err %v", normalizedPhone, message, err)
		return Error{
			Type:    ErrProviderFailure,
			Message: fmt.Sprintf("Failed to send SMS: %v", err),
		}
	}

	ns.log.Printf("[SendSMS] SMS sent successfully to: %s", normalizedPhone)
	return nil
}

// SendTransactionNotification sends transaction notifications to multiple recipients
func (ns *NotificationServiceImpl) SendTransactionNotification(ctx context.Context, req TransactionNotificationRequest) error {
	ns.log.Printf("[SendTransactionNotification] Sending transaction notification for: %s", req.TransactionID)

	for _, recipient := range req.Recipients {
		switch recipient.Type {
		case entity.TypeSMS:
			message := ns.buildTransactionSMSMessage(recipient, req.TransactionData)
			if err := ns.SendSMS(ctx, recipient.Identifier, message); err != nil {
				ns.log.Printf("[SendTransactionNotification] Failed to send SMS to %s (%s): %v",
					recipient.Name, recipient.Identifier, err)
				// Continue with other recipients even if one fails
				continue
			}
			ns.log.Printf("[SendTransactionNotification] SMS sent to %s (%s)",
				recipient.Name, recipient.Role)

		case entity.TypeEmail:
			// TODO: Implement email notifications
			ns.log.Printf("[SendTransactionNotification] Email notifications not implemented yet")

		case entity.TypeInApp:
			// TODO: Implement in-app notifications
			ns.log.Printf("[SendTransactionNotification] In-app notifications not implemented yet")
		}
	}

	return nil
}

// SendTemplatedNotification sends a notification using a predefined template
func (ns *NotificationServiceImpl) SendTemplatedNotification(ctx context.Context, req TemplatedNotificationRequest) error {
	ns.log.Printf("[SendTemplatedNotification] Sending templated notification: %s", req.Template)

	message, err := ns.buildTemplatedMessage(req.Template, req.Data)
	if err != nil {
		return Error{
			Type:    ErrInvalidTemplate,
			Message: err.Error(),
		}
	}

	switch req.Type {
	case entity.TypeSMS:
		return ns.SendSMS(ctx, req.Recipient, message)
	case entity.TypeEmail:
		// TODO: Implement email
		return Error{
			Type:    ErrProviderFailure,
			Message: "Email notifications not implemented yet",
		}
	case entity.TypeInApp:
		// TODO: Implement in-app
		return Error{
			Type:    ErrProviderFailure,
			Message: "In-app notifications not implemented yet",
		}
	default:
		return Error{
			Type:    ErrInvalidTemplate,
			Message: "Invalid notification type",
		}
	}
}

// GetNotificationStatus retrieves the status of a notification (placeholder for future implementation)
func (ns *NotificationServiceImpl) GetNotificationStatus(ctx context.Context, notificationID uuid.UUID) (*entity.Notification, error) {
	// TODO: Implement notification status tracking with database
	return nil, Error{
		Type:    ErrNotificationNotFound,
		Message: "Notification status tracking not implemented yet",
	}
}

// buildTransactionSMSMessage creates an SMS message for transaction notifications
func (ns *NotificationServiceImpl) buildTransactionSMSMessage(recipient NotificationRecipient, data entity.TransactionNotificationData) string {
	// Generate current timestamp in East Africa Time (UTC+3) - equivalent to Addis Ababa
	// Using fixed offset instead of timezone name for better compatibility
	loc := time.FixedZone("EAT", 3*60*60) // UTC+3 (3 hours * 60 minutes * 60 seconds)
	now := time.Now().In(loc)
	dateStr := now.Format("02 Jan 2006")
	timeStr := now.Format("3:04 PM")

	if data.MerchantName == "" {
		data.MerchantName = "Merchant"
	}

	var message string

	switch strings.ToLower(recipient.Role) {
	case "payer", "customer":
		if data.Status == "SUCCESS" {
			message = fmt.Sprintf(
				`Dear %s,
Your payment of %.2f %s to %s has been successfully processed.
Reference: %s
Date: %s at %s

📞6562
አሸናፊዎች በላኪፔይ ይከፍላሉ!!`,
				recipient.Name,
				data.Amount,
				data.Currency,
				data.MerchantName,
				data.Reference,
				dateStr,
				timeStr,
			)
		} else {
			message = fmt.Sprintf(
				`Dear %s,
Your payment of %.2f %s to %s has failed.
Reference: %s
Date: %s at %s

Please try again or contact support.`,
				recipient.Name,
				data.Amount,
				data.Currency,
				data.MerchantName,
				data.Reference,
				dateStr,
				timeStr,
			)
		}

	case "merchant":
		if data.Status == "SUCCESS" {
			message = fmt.Sprintf(
				`Dear %s,
You received %.2f %s from %s.
Reference: %s
Date: %s at %s

Social Pay - Your trusted payment partner!`,
				recipient.Name,
				data.Amount,
				data.Currency,
				data.CustomerName,
				data.Reference,
				dateStr,
				timeStr,
			)
		} else {
			message = fmt.Sprintf(
				`Dear %s,
Payment of %.2f %s from %s failed.
Reference: %s
Date: %s at %s

Transaction was not completed.`,
				recipient.Name,
				data.Amount,
				data.Currency,
				data.CustomerName,
				data.Reference,
				dateStr,
				timeStr,
			)
		}

	case "tipee":
		if data.Status == "SUCCESS" {
			message = fmt.Sprintf(
				`Dear %s,
You received a tip of %.2f %s.
Reference: %s
Date: %s at %s

Thank you for your service!`,
				recipient.Name,
				data.Amount,
				data.Currency,
				data.Reference,
				dateStr,
				timeStr,
			)
		}

	default:
		message = fmt.Sprintf(
			`Transaction %s: %.2f %s
Status: %s
Reference: %s
Date: %s at %s`,
			data.TransactionType,
			data.Amount,
			data.Currency,
			data.Status,
			data.Reference,
			dateStr,
			timeStr,
		)
	}

	return message
}

// buildTemplatedMessage creates a message from a template
func (ns *NotificationServiceImpl) buildTemplatedMessage(template entity.NotificationTemplate, data map[string]interface{}) (string, error) {
	switch template {
	case entity.TemplateOTP:
		if otp, ok := data["otp"].(string); ok {
			return fmt.Sprintf("Your Social Pay verification code is %s. Do not share this code with anyone.", otp), nil
		}
		return "", fmt.Errorf("OTP template requires 'otp' field")

	case entity.TemplateGeneral:
		if message, ok := data["message"].(string); ok {
			return message, nil
		}
		return "", fmt.Errorf("General template requires 'message' field")

	default:
		return "", fmt.Errorf("Template %s not implemented", template)
	}
}

// normalizeEthiopianPhoneNumber accepts various formats and returns 251xxxxxxxxx
func normalizeEthiopianPhoneNumber(phone string) (string, error) {
	// Remove all non-digit characters except +
	cleaned := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Remove + if present
	cleaned = strings.TrimPrefix(cleaned, "+")

	// Handle different formats
	switch {
	case len(cleaned) == 9 && cleaned[0] == '9': // 9xxxxxxxx
		return "+251" + cleaned, nil
	case len(cleaned) == 10 && cleaned[:2] == "09": // 09xxxxxxxx
		return "+251" + cleaned[1:], nil
	case len(cleaned) == 12 && strings.HasPrefix(cleaned, "251"): // 2519xxxxxxxx
		return "+" + cleaned, nil
	default:
		return "", fmt.Errorf("invalid Ethiopian phone number format. Accepted formats: 9xxxxxxxx, 09xxxxxxxx, 2519xxxxxxxx, +2519xxxxxxxx")
	}
}
