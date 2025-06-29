package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/notifications/core/entity"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	v2MerchantRepo "github.com/socialpay/socialpay/src/pkg/v2_merchant/core/repository"
)

// TransactionNotifier handles sending notifications for transaction events
type TransactionNotifier struct {
	notificationService NotificationService
	merchantRepo        v2MerchantRepo.Repository
	log                 logging.Logger
}

// NewTransactionNotifier creates a new transaction notifier
func NewTransactionNotifier(
	notificationService NotificationService,
	merchantRepo v2MerchantRepo.Repository,
) *TransactionNotifier {
	return &TransactionNotifier{
		notificationService: notificationService,
		merchantRepo:        merchantRepo,
		log:                 logging.NewStdLogger("[notifications]"),
	}
}

// NotifyTransactionStatus sends notifications to all parties involved in a transaction
func (tn *TransactionNotifier) NotifyTransactionStatus(ctx context.Context, transaction *txEntity.Transaction, status string) error {
	if transaction == nil {
		return fmt.Errorf("transaction is required")
	}

	tn.log.Info("[TransactionNotifier] Processing transaction notification", map[string]interface{}{
		"transactionID": transaction.Id,
	})

	var recipients []NotificationRecipient
	transactionData := entity.TransactionNotificationData{
		TransactionID:   transaction.Id.String(),
		Amount:          transaction.Amount,
		Currency:        transaction.Currency,
		Status:          status,
		Reference:       transaction.Reference,
		Timestamp:       time.Now().Format("2006-01-02 15:04:05"),
		TransactionType: string(transaction.Type),
	}

	// For customers/payers: only use phone number from transaction without fetching user data
	if transaction.PhoneNumber != "" {
		// We don't have customer name, so we'll use "Customer" as a generic name
		transactionData.CustomerName = "Customer"

		recipients = append(recipients, NotificationRecipient{
			Type:       entity.TypeSMS,
			Identifier: transaction.PhoneNumber,
			Name:       "Customer", // Generic name since we don't fetch customer data
			Role:       "payer",
		})
		tn.log.Info("[TransactionNotifier] Added payer", map[string]interface{}{
			"phoneNumber": transaction.PhoneNumber,
		})
	}

	// For merchants: fetch merchant information if transaction has merchant ID
	if transaction.MerchantId != uuid.Nil {
		merchant, err := tn.merchantRepo.GetMerchantDetails(ctx, transaction.MerchantId)
		tn.log.Info("[TransactionNotifier] Merchant details", map[string]interface{}{
			"merchant": merchant,
		})
		if err != nil {
			tn.log.Error("[TransactionNotifier] Failed to get merchant details", map[string]interface{}{
				"error": err,
			})
			// Continue without merchant notification rather than failing
		} else if merchant != nil {
			// Access fields from the nested Merchant struct
			merchantName := ""
			if merchant.Merchant.TradingName != nil && *merchant.Merchant.TradingName != "" {
				merchantName = *merchant.Merchant.TradingName
			} else {
				merchantName = merchant.Merchant.LegalName
			}
			transactionData.MerchantName = merchantName

			// Get merchant phone number from contacts
			var merchantPhone string
			for _, contact := range merchant.Contacts {
				if contact.ContactType == "primary" || contact.ContactType == "business" {
					merchantPhone = contact.PhoneNumber
					break
				}
			}
			// If no primary contact, use the first available contact
			if merchantPhone == "" && len(merchant.Contacts) > 0 {
				merchantPhone = merchant.Contacts[0].PhoneNumber
			}

			if merchantPhone != "" {
				recipients = append(recipients, NotificationRecipient{
					Type:       entity.TypeSMS,
					Identifier: merchantPhone,
					Name:       merchantName,
					Role:       "merchant",
				})
				tn.log.Info("[TransactionNotifier] Added merchant", map[string]interface{}{
					"merchantName":  merchantName,
					"merchantPhone": merchantPhone,
				})
			} else {
				tn.log.Info("[TransactionNotifier] Merchant has no phone number", map[string]interface{}{
					"merchantName": merchantName,
				})
			}
		}
	}

	// For tipees: handle if transaction has tipee information
	if transaction.TipeePhone != nil && *transaction.TipeePhone != "" {
		recipients = append(recipients, NotificationRecipient{
			Type:       entity.TypeSMS,
			Identifier: *transaction.TipeePhone,
			Name:       "Tipee", // We don't have tipee name in transaction
			Role:       "tipee",
		})
		tn.log.Info("[TransactionNotifier] Added tipee", map[string]interface{}{
			"tipeePhone": *transaction.TipeePhone,
		})
	}

	if len(recipients) == 0 {
		tn.log.Info("[TransactionNotifier] No recipients found for transaction", map[string]interface{}{
			"transactionID": transaction.Id,
		})
		return nil
	}

	// Send notifications to all recipients
	notificationReq := TransactionNotificationRequest{
		TransactionID:   transaction.Id.String(),
		Type:            entity.TypeSMS,
		Recipients:      recipients,
		TransactionData: transactionData,
	}

	err := tn.notificationService.SendTransactionNotification(ctx, notificationReq)
	if err != nil {
		tn.log.Error("[TransactionNotifier] Failed to send transaction notifications", map[string]interface{}{
			"error": err,
		})
		return fmt.Errorf("failed to send transaction notifications: %w", err)
	}

	tn.log.Info("[TransactionNotifier] Successfully sent notifications to %d recipients", map[string]interface{}{
		"recipients": len(recipients),
	})
	return nil
}

// NotifyTransactionStatusByPhone sends a simple notification to a specific phone number
func (tn *TransactionNotifier) NotifyTransactionStatusByPhone(ctx context.Context, phoneNumber, customerName, merchantName string, amount float64, currency, status, reference string) error {
	tn.log.Info("[TransactionNotifier] Sending notification to phone", map[string]interface{}{
		"phoneNumber": phoneNumber,
	})

	transactionData := entity.TransactionNotificationData{
		CustomerName: customerName,
		MerchantName: merchantName,
		Amount:       amount,
		Currency:     currency,
		Status:       status,
		Reference:    reference,
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
	}

	recipient := NotificationRecipient{
		Type:       entity.TypeSMS,
		Identifier: phoneNumber,
		Name:       customerName,
		Role:       "payer",
	}

	req := TransactionNotificationRequest{
		Type:            entity.TypeSMS,
		Recipients:      []NotificationRecipient{recipient},
		TransactionData: transactionData,
	}

	return tn.notificationService.SendTransactionNotification(ctx, req)
}
