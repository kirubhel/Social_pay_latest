package notifications

import (
	"log"

	"github.com/socialpay/socialpay/src/pkg/notifications/adapter/gateway/sms"
	"github.com/socialpay/socialpay/src/pkg/notifications/usecase"
	v2MerchantRepo "github.com/socialpay/socialpay/src/pkg/v2_merchant/core/repository"
)

// NewNotificationService creates a new notification service with AfroSMS provider
func NewNotificationService(log *log.Logger) usecase.NotificationService {
	// Initialize SMS provider
	smsProvider := sms.New(log)

	// Create and return notification service
	return usecase.NewNotificationService(smsProvider, log)
}

// NewTransactionNotifier creates a new transaction notifier with merchant repository
func NewTransactionNotifier(merchantRepo v2MerchantRepo.Repository, log *log.Logger) *usecase.TransactionNotifier {
	// Create notification service
	notificationService := NewNotificationService(log)

	// Create and return transaction notifier
	return usecase.NewTransactionNotifier(notificationService, merchantRepo)
}
