package dto

import (
	"time"

	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

type WebhookRequest struct {
	Type             txEntity.TransactionType `json:"type"`
	TransactionID    string                   `json:"transactionId" binding:"required"`
	Status           string                   `json:"status" binding:"required"`
	Message          string                   `json:"message" binding:"required"`
	ProviderTxID     string                   `json:"providerTxId"`
	ProviderData     string                   `json:"providerData"`
	Timestamp        time.Time                `json:"timestamp" binding:"required"`
	CallbackURL      string                   `json:"callbackUrl" binding:"required"`
	MerchantID       string                   `json:"merchantId" binding:"required"`
	UserID           string                   `json:"userId" binding:"required"`
	IsHostedCheckout bool                     `json:"isHostedCheckout"`
}

type WebhookEventMerchant struct {
	Event        txEntity.TransactionType `json:"event"`
	ReferenceId  string                   `json:"referenceId"`
	SocialPayTxnID string                   `json:"socialpayTxnId"`
	Status       string                   `json:"status"`
	Amount       string                   `json:"amount"`
	CallbackURL  string                   `json:"callbackUrl"`
	Message      string                   `json:"message"`
	ProviderTxID string                   `json:"providerTxId"`
	Timestamp    time.Time                `json:"timestamp"`
	MerchantID   string                   `json:"merchantId"`
	UserID       string                   `json:"userId"`
}

type WebhookMessage struct {
	Type             txEntity.TransactionType `json:"type"`
	TransactionID    string                   `json:"transaction_id"`
	MerchantID       string                   `json:"merchant_id"`
	TotalAmount      float64                  `json:"total_amount"`
	Status           string                   `json:"status"`
	Message          string                   `json:"message"`
	ProviderTxID     string                   `json:"provider_txid"`
	ProviderData     string                   `json:"provider_data"`
	Timestamp        time.Time                `json:"timestamp"`
	UserID           string                   `json:"user_id"`
	IsHostedCheckout bool                     `json:"isHostedCheckout"`
}
