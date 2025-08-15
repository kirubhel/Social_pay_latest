package usecase

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/google/uuid"

	commission_usecase "github.com/socialpay/socialpay/src/pkg/commission/usecase"
	socialPayEntity "github.com/socialpay/socialpay/src/pkg/socialpayapi/core/entity"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
	txRepo "github.com/socialpay/socialpay/src/pkg/transaction/core/repository"
	transaction_usecase "github.com/socialpay/socialpay/src/pkg/transaction/usecase"
	merchantEntity "github.com/socialpay/socialpay/src/pkg/v2_merchant/core/entity"
	v2MerchantUsecase "github.com/socialpay/socialpay/src/pkg/v2_merchant/usecase"
	walletEntity "github.com/socialpay/socialpay/src/pkg/wallet/core/entity"
	walletUsecase "github.com/socialpay/socialpay/src/pkg/wallet/usecase"
)

// PaymentUseCase defines the interface for payment operations
type PaymentUseCase interface {
	// DirectPay handles direct payment requests
	ProcessDirectPayment(ctx context.Context, apikey string, userID uuid.UUID, merchantID uuid.UUID, req *socialPayEntity.DirectPaymentRequest) (*socialPayEntity.PaymentResponse, error)

	// GetTransaction retrieves transaction details by ID
	GetTransactionWithMerchant(ctx context.Context, id uuid.UUID) (*txEntity.Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*txEntity.Transaction, error)

	// QueryTransactionStatus queries transaction status from provider
	QueryTransactionStatus(ctx context.Context, medium txEntity.TransactionMedium, transactionID string) (*payment.TransactionStatusQueryResponse, error)

	// RequestWithdrawal handles withdrawal requests
	RequestWithdrawal(ctx context.Context, apiKey string, userID uuid.UUID, merchantID uuid.UUID, req *socialPayEntity.WithdrawalRequest) (*socialPayEntity.PaymentResponse, error)

	// GetWalletBalance retrieves the wallet balance for a merchant
	GetWalletBalance(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID) (*walletEntity.MerchantWallet, error)

	// Hosted Checkout methods
	CreateHostedCheckout(ctx context.Context, apikey string, userID uuid.UUID, merchantID uuid.UUID, req *socialPayEntity.HostedCheckoutRequest) (*socialPayEntity.PaymentResponse, error)
	UpdateHostedCheckout(ctx context.Context, apikey string, userID uuid.UUID, merchantID uuid.UUID, id uuid.UUID, req *socialPayEntity.UpdateHostedCheckoutRequest) (*socialPayEntity.PaymentResponse, error)
	GetHostedCheckout(ctx context.Context, id uuid.UUID) (*socialPayEntity.HostedCheckoutResponseDTO, error)
	GetHostedCheckoutWithMerchant(ctx context.Context, id uuid.UUID) (*socialPayEntity.HostedCheckoutWithMerchantResponseDTO, error)
	ProcessCheckoutPayment(ctx context.Context, req *socialPayEntity.CheckoutPaymentRequest) (*socialPayEntity.PaymentResponse, error)
}

type paymentUseCase struct {
	transactionRepo            txRepo.TransactionRepository
	hostedPaymentRepo          txRepo.HostedPaymentRepository
	transactionUseCase         transaction_usecase.TransactionUseCase
	walletUseCase              walletUsecase.MerchantWalletUsecase
	merchantUseCase            v2MerchantUsecase.MerchantUseCase
	paymentService             PaymentProcessor
	transactionCreationService *TransactionCreationService
	log                        logging.Logger
}

func (uc *paymentUseCase) ProcessDirectPayment(ctx context.Context, apikey string, userID uuid.UUID, merchantID uuid.UUID, req *socialPayEntity.DirectPaymentRequest) (*socialPayEntity.PaymentResponse, error) {

	uc.log.Info("Starting direct payment processing", map[string]interface{}{
		"user_id":      userID,
		"merchant_id":  merchantID,
		"medium":       req.Medium,
		"phone_number": req.PhoneNumber,
	})

	// validate reference id
	if err := uc.transactionUseCase.ValidateReferenceId(ctx, merchantID, req.Reference); err != nil {
		return nil, err
	}

	// Use unified transaction creation service
	txCreationReq := TransactionCreationRequest{
		UserID:          userID,
		MerchantID:      merchantID,
		BaseAmount:      req.Amount,
		Description:     req.Description,
		Medium:          req.Medium,
		Type:            txEntity.DEPOSIT,
		Status:          txEntity.INITIATED,
		PaymentType:     "direct",
		MerchantPaysFee: req.MerchantPaysFee,
		CallbackURL:     req.CallbackURL,
		SuccessURL:      req.Redirects.Success,
		FailedURL:       req.Redirects.Failed,
	}

	txCreationResp, err := uc.transactionCreationService.CreateTransaction(ctx, txCreationReq)
	if err != nil {
		uc.log.Error("Failed to create transaction using unified service", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	tx := txCreationResp.Transaction
	tx.PhoneNumber = req.PhoneNumber
	tx.Currency = req.Currency
	tx.Reference = req.Reference

	uc.log.Info("Created transaction using unified service", map[string]interface{}{
		"transaction_id":    tx.Id,
		"original_amount":   tx.BaseAmount,
		"total_amount":      tx.TotalAmount,
		"commission_rate":   txCreationResp.CommissionRate,
		"merchant_pays_fee": req.MerchantPaysFee,
		"tip_amount":        txCreationResp.TipAmount,
	})

	// Store transaction
	if err := uc.transactionRepo.Create(ctx, tx); err != nil {
		uc.log.Error("Failed to create transaction", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	uc.log.Info("Successfully stored transaction", map[string]interface{}{
		"transaction_id": tx.Id,
	})

	// Process payment using payment service
	paymentReq := &payment.PaymentRequest{
		TransactionID: tx.Id,
		Medium:        req.Medium,
		Amount:        tx.CustomerNet,
		Currency:      tx.Currency,
		PhoneNumber:   tx.PhoneNumber,
		Reference:     tx.Reference,
		Description:   tx.Description,
		CallbackURL:   tx.CallbackURL,
		SuccessURL:    tx.SuccessURL,
		FailedURL:     tx.FailedURL,
	}

	uc.log.Info("Initiating payment processing", map[string]interface{}{
		"transaction_id": tx.Id,
	})
	paymentResp, err := uc.paymentService.ProcessPayment(ctx, apikey, paymentReq)
	if err != nil {
		uc.log.Error("Payment processing failed", map[string]interface{}{
			"error": err.Error(),
		})
		tx.Status = txEntity.TransactionStatus(txEntity.FAILED)
		tx.Comment = err.Error()
		_ = uc.transactionRepo.Update(ctx, tx)
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}
	uc.log.Info("Payment processing completed", map[string]interface{}{
		"status": paymentResp.Status,
	})

	// Update transaction status
	tx.Status = txEntity.TransactionStatus(paymentResp.Status)
	if err := uc.transactionRepo.Update(ctx, tx); err != nil {
		uc.log.Error("Failed to update transaction status", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}
	// loging provider_tx_id
	uc.log.Info("provider_tx_id", map[string]interface{}{
		"provider_tx_id": paymentResp.ProcessorRef,
	})

	param := map[string]interface{}{
		"provider_tx_id": paymentResp.ProcessorRef,
	}
	// Setting provider_id
	if err := uc.transactionRepo.UpdateTransactionWithProviderData(ctx, tx.Id, param); err != nil {
		// loging
		uc.log.Error("error while setting provider_tx_id", map[string]interface{}{
			"err": err,
		})
	}

	uc.log.Info("Successfully updated transaction status", map[string]interface{}{
		"status": tx.Status,
	})

	return &socialPayEntity.PaymentResponse{
		Success:              paymentResp.Success,
		Status:               string(paymentResp.Status),
		Message:              paymentResp.Message,
		Reference:            tx.Reference,
		PaymentURL:           paymentResp.PaymentURL,
		SocialPayTransactionID: tx.Id.String(),
	}, nil
}

func (uc *paymentUseCase) GetTransactionWithMerchant(ctx context.Context, id uuid.UUID) (*txEntity.Transaction, error) {
	tx, err := uc.transactionRepo.GetByIDWithMerchant(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return tx, nil
}

func (uc *paymentUseCase) GetTransaction(ctx context.Context, id uuid.UUID) (*txEntity.Transaction, error) {
	tx, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return tx, nil
}

// RequestWithdrawal handles withdrawal requests
func (uc *paymentUseCase) RequestWithdrawal(ctx context.Context, apikey string, userID uuid.UUID, merchantID uuid.UUID, req *socialPayEntity.WithdrawalRequest) (*socialPayEntity.PaymentResponse, error) {
	uc.log.Info("[Withdrawal] Starting withdrawal request", map[string]interface{}{
		"user_id":     userID,
		"merchant_id": merchantID,
		"amount":      req.Amount,
		"currency":    req.Currency,
	})

	if err := uc.transactionUseCase.ValidateReferenceId(ctx, merchantID, req.Reference); err != nil {
		return nil, err
	}

	// Use unified transaction creation service for withdrawal
	txCreationReq := TransactionCreationRequest{
		UserID:          userID,
		MerchantID:      merchantID,
		BaseAmount:      req.Amount,
		Description:     "Withdrawal for phone number: " + req.PhoneNumber, // Withdrawal doesn't have description in request
		Medium:          req.Medium,
		Type:            txEntity.WITHDRAWAL,
		PaymentType:     "withdrawal",
		MerchantPaysFee: req.MerchantPaysFee,
		CallbackURL:     req.CallbackURL,
	}

	txCreationResp, txErr := uc.transactionCreationService.CreateTransaction(ctx, txCreationReq)
	if txErr != nil {
		uc.log.Error("Failed to create withdrawal transaction using unified service", map[string]interface{}{
			"error": txErr.Error(),
		})
		return nil, fmt.Errorf("failed to create withdrawal transaction: %w", txErr)
	}

	tx := txCreationResp.Transaction
	tx.PhoneNumber = req.PhoneNumber
	tx.Currency = req.Currency
	tx.Reference = req.Reference

	uc.log.Info("[Withdrawal] Created withdrawal transaction", map[string]interface{}{
		"transaction_id":  tx.Id,
		"original amount": tx.BaseAmount,
		"currency":        tx.Currency,
	})

	uc.log.Info("[Withdrawal] Calculated fees", map[string]interface{}{
		"original_amount": tx.BaseAmount,
		"fee_amount":      tx.FeeAmount,
		"vat_amount":      tx.VatAmount,
		"admin_net":       tx.AdminNet,
		"merchant_net":    tx.MerchantNet,
		"total_amount":    tx.TotalAmount,
		"customer_net":    tx.CustomerNet,
	})

	// Get the wallet and lock the withdrawal amount
	// This uses row-level locking to prevent race conditions
	err := uc.walletUseCase.LockWithdrawalAmount(ctx, merchantID, tx.MerchantNet)
	uc.log.Info("[Withdrawal] Locked withdrawal amount", map[string]interface{}{
		"merchant_id": merchantID,
		"amount":      tx.MerchantNet,
	})
	if err != nil {
		uc.log.Error("[Withdrawal] Failed to lock withdrawal amount", map[string]interface{}{
			"error":       err.Error(),
			"merchant_id": merchantID,
			"amount":      tx.TotalAmount,
		})
		return nil, err
	}

	// Save the transaction
	if err := uc.transactionRepo.Create(ctx, tx); err != nil {
		uc.log.Error("[Withdrawal] Failed to create withdrawal transaction", map[string]interface{}{
			"error": err.Error(),
		})

		// Try to unlock the amount since the transaction creation failed
		unlockErr := uc.walletUseCase.ProcessTransactionStatus(ctx, merchantID, tx.MerchantNet, 0, false, true)
		if unlockErr != nil {
			uc.log.Error("[Withdrawal] Failed to unlock withdrawal amount after transaction creation failure", map[string]interface{}{
				"error":        unlockErr.Error(),
				"merchant_id":  merchantID,
				"total_amount": tx.TotalAmount,
				"amount":       tx.BaseAmount,
			})
		}

		return nil, fmt.Errorf("failed to create withdrawal transaction: %w", err)
	}
	uc.log.Info("[Withdrawal] Successfully stored withdrawal transaction", map[string]interface{}{
		"transaction_id": tx.Id,
	})

	// Process withdrawal
	uc.log.Info("[Withdrawal] Initiating withdrawal processing", map[string]interface{}{
		"transaction_id": tx.Id,
	})

	// Process payment using payment service
	paymentReq := &payment.PaymentRequest{
		TransactionID: tx.Id,
		Medium:        req.Medium,
		Amount:        tx.CustomerNet,
		Currency:      tx.Currency,
		PhoneNumber:   tx.PhoneNumber,
		Reference:     tx.Reference,
		Description:   tx.Description,
		CallbackURL:   tx.CallbackURL,
		SuccessURL:    tx.SuccessURL,
		FailedURL:     tx.FailedURL,
	}

	withdrawalResponse, err := uc.paymentService.ProcessWithdrawal(ctx, apikey, paymentReq)

	if err != nil {
		uc.log.Error("[Withdrawal] Withdrawal processing failed", map[string]interface{}{
			"error": err.Error(),
		})
		tx.Status = txEntity.FAILED
		tx.Comment = err.Error()
		_ = uc.transactionRepo.Update(ctx, tx)

		// Unlock the amount since the withdrawal failed
		unlockErr := uc.walletUseCase.ProcessTransactionStatus(ctx, merchantID, tx.MerchantNet, 0, false, true)
		if unlockErr != nil {
			uc.log.Error("[Withdrawal] Failed to unlock withdrawal amount after processing failure", map[string]interface{}{
				"error":       unlockErr.Error(),
				"merchant_id": merchantID,
				"amount":      tx.TotalAmount,
			})
		}

		return nil, fmt.Errorf("failed to process withdrawal: %w", err)
	}
	uc.log.Info("[Withdrawal] Withdrawal processing completed", map[string]interface{}{
		"transaction_id": tx.Id,
	})

	tx.Status = txEntity.TransactionStatus(withdrawalResponse.Status)
	if err := uc.transactionRepo.Update(ctx, tx); err != nil {
		uc.log.Error("Failed to update withdrawal transaction status", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}
	uc.log.Info("Successfully updated withdrawal transaction status", map[string]interface{}{
		"status": tx.Status,
	})

	return &socialPayEntity.PaymentResponse{
		Success:              withdrawalResponse.Success,
		Status:               string(withdrawalResponse.Status),
		Message:              withdrawalResponse.Message,
		Reference:            tx.Reference,
		PaymentURL:           withdrawalResponse.PaymentURL,
		SocialPayTransactionID: tx.Id.String(),
	}, nil
}

func (uc *paymentUseCase) GetWalletBalance(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID) (*walletEntity.MerchantWallet, error) {
	uc.log.Info("Getting wallet balance", map[string]interface{}{
		"user_id":     userID,
		"merchant_id": merchantID,
	})

	wallet, err := uc.walletUseCase.GetMerchantWallet(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant wallet: %w", err)
	}

	// Convert from wallet entity to webhook entity
	return &walletEntity.MerchantWallet{
		ID:           wallet.ID,
		UserID:       wallet.UserID,
		Amount:       wallet.Amount,
		LockedAmount: wallet.LockedAmount,
		Currency:     walletEntity.Currency(wallet.Currency),
		CreatedAt:    wallet.CreatedAt,
		UpdatedAt:    wallet.UpdatedAt,
	}, nil
}

// CreateHostedCheckout creates a hosted checkout session
func (uc *paymentUseCase) CreateHostedCheckout(ctx context.Context, apikey string, userID uuid.UUID, merchantID uuid.UUID, req *socialPayEntity.HostedCheckoutRequest) (*socialPayEntity.PaymentResponse, error) {
	uc.log.Info("Starting hosted checkout creation", map[string]interface{}{
		"user_id":     userID,
		"merchant_id": merchantID,
		"amount":      req.Amount,
		"currency":    req.Currency,
	})

	// Validate reference ID
	if err := uc.transactionUseCase.ValidateReferenceId(ctx, merchantID, req.Reference); err != nil {
		return nil, err
	}

	// hosted checkout validate reference id
	if err := uc.hostedPaymentRepo.ValidateReferenceId(ctx, merchantID, req.Reference); err != nil {
		return nil, err
	}

	// Set default expiry time if not provided (24 hours from now)
	var expiresAt time.Time
	if req.ExpiresAt != nil {
		expiresAt = req.ExpiresAt.UTC()
	}

	// Create hosted payment
	hostedPayment := &txEntity.HostedPayment{
		ID:               uuid.New(),
		UserID:           userID,
		MerchantID:       merchantID,
		Amount:           req.Amount,
		Currency:         req.Currency,
		Description:      req.Description,
		Reference:        req.Reference,
		SupportedMediums: req.SupportedMediums,
		PhoneNumber:      req.PhoneNumber,
		SuccessURL:       req.Redirects.Success,
		FailedURL:        req.Redirects.Failed,
		CallbackURL:      req.CallbackURL,
		Status:           txEntity.HostedPaymentPending,
		MerchantPaysFee:  req.MerchantPaysFee,
		ExpiresAt:        expiresAt,
		AcceptTip:        req.AcceptTip,
	}

	// Store hosted payment
	if err := uc.hostedPaymentRepo.Create(ctx, hostedPayment); err != nil {
		uc.log.Error("Failed to create hosted payment", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create hosted payment: %w", err)
	}

	// Get checkout URL from environment
	checkoutURL := os.Getenv("APP_CHECKOUT_URL")
	if checkoutURL == "" {
		checkoutURL = "http://localhost:3000" // fallback
	}
	paymentURL := fmt.Sprintf("%s/checkout/%s", checkoutURL, hostedPayment.ID.String())

	uc.log.Info("Successfully created hosted checkout", map[string]interface{}{
		"hosted_payment_id": hostedPayment.ID,
		"payment_url":       paymentURL,
		"expires_at":        hostedPayment.ExpiresAt,
	})

	return &socialPayEntity.PaymentResponse{
		Success:              true,
		Status:               string(txEntity.HostedPaymentPending),
		Message:              "Hosted checkout created successfully",
		PaymentURL:           paymentURL,
		Reference:            req.Reference,
		SocialPayTransactionID: hostedPayment.ID.String(),
	}, nil
}

// UpdateHostedCheckout updates an existing hosted checkout session
func (uc *paymentUseCase) UpdateHostedCheckout(ctx context.Context, apikey string, userID uuid.UUID, merchantID uuid.UUID, id uuid.UUID, req *socialPayEntity.UpdateHostedCheckoutRequest) (*socialPayEntity.PaymentResponse, error) {
	uc.log.Info("Starting hosted checkout update", map[string]interface{}{
		"user_id":           userID,
		"merchant_id":       merchantID,
		"hosted_payment_id": id,
	})

	// Get existing hosted payment
	existingPayment, err := uc.hostedPaymentRepo.GetByID(ctx, id)
	if err != nil {
		uc.log.Error("Failed to get hosted payment for update", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("hosted payment not found: %w", err)
	}

	// Verify ownership
	if existingPayment.UserID != userID || existingPayment.MerchantID != merchantID {
		return nil, fmt.Errorf("access denied: hosted payment does not belong to this merchant")
	}

	// Check if hosted payment is still pending and not expired
	if existingPayment.Status != txEntity.HostedPaymentPending {
		return nil, fmt.Errorf("cannot update hosted payment: status is %s", existingPayment.Status)
	}

	if time.Now().UTC().After(existingPayment.ExpiresAt) {
		return nil, fmt.Errorf("cannot update hosted payment: it has expired")
	}

	// Update fields if provided
	updated := false

	if req.Amount != nil {
		existingPayment.Amount = *req.Amount
		updated = true
	}

	if req.Currency != nil {
		existingPayment.Currency = *req.Currency
		updated = true
	}

	if req.Description != nil {
		existingPayment.Description = *req.Description
		updated = true
	}

	if len(req.SupportedMediums) > 0 {
		existingPayment.SupportedMediums = req.SupportedMediums
		updated = true
	}

	if req.MerchantPaysFee != nil {
		existingPayment.MerchantPaysFee = *req.MerchantPaysFee
		updated = true
	}

	if req.PhoneNumber != nil {
		existingPayment.PhoneNumber = *req.PhoneNumber
		updated = true
	}

	if req.Redirects != nil {
		existingPayment.SuccessURL = req.Redirects.Success
		existingPayment.FailedURL = req.Redirects.Failed
		updated = true
	}

	if req.CallbackURL != nil {
		existingPayment.CallbackURL = *req.CallbackURL
		updated = true
	}

	if req.ExpiresAt != nil {
		existingPayment.ExpiresAt = req.ExpiresAt.UTC()
		updated = true
	}

	if !updated {
		return nil, fmt.Errorf("no fields provided for update")
	}

	// Update timestamp
	existingPayment.UpdatedAt = time.Now().UTC()

	// Store updated hosted payment
	if err := uc.hostedPaymentRepo.Update(ctx, existingPayment); err != nil {
		uc.log.Error("Failed to update hosted payment", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to update hosted payment: %w", err)
	}

	// Get checkout URL from environment
	checkoutURL := os.Getenv("APP_CHECKOUT_URL")
	if checkoutURL == "" {
		checkoutURL = "http://localhost:3000" // fallback
	}
	paymentURL := fmt.Sprintf("%s/checkout/%s", checkoutURL, existingPayment.ID.String())

	uc.log.Info("Successfully updated hosted checkout", map[string]interface{}{
		"hosted_payment_id": existingPayment.ID,
		"payment_url":       paymentURL,
		"expires_at":        existingPayment.ExpiresAt,
	})

	return &socialPayEntity.PaymentResponse{
		Success:              true,
		Status:               string(existingPayment.Status),
		Message:              "Hosted checkout updated successfully",
		PaymentURL:           paymentURL,
		Reference:            existingPayment.Reference,
		SocialPayTransactionID: existingPayment.ID.String(),
	}, nil
}

// GetHostedCheckout retrieves hosted checkout details by ID
func (uc *paymentUseCase) GetHostedCheckout(ctx context.Context, id uuid.UUID) (*socialPayEntity.HostedCheckoutResponseDTO, error) {
	uc.log.Info("Getting hosted checkout details", map[string]interface{}{
		"hosted_checkout_id": id,
	})

	// Get hosted payment
	hostedPayment, err := uc.hostedPaymentRepo.GetByID(ctx, id)
	if err != nil {
		uc.log.Error("Failed to get hosted payment", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("hosted payment not found: %w", err)
	}

	// Check if hosted payment has expired
	if time.Now().After(hostedPayment.ExpiresAt) {
		uc.log.Info("Hosted payment has expired", map[string]interface{}{
			"hosted_checkout_id": id,
			"expires_at":         hostedPayment.ExpiresAt,
		})
		return nil, fmt.Errorf("hosted payment has expired")
	}

	// Check if hosted payment is still pending
	if hostedPayment.Status != txEntity.HostedPaymentPending {
		uc.log.Info("Hosted payment is no longer pending", map[string]interface{}{
			"hosted_checkout_id": id,
			"status":             hostedPayment.Status,
		})
		return nil, fmt.Errorf("hosted payment is no longer available")
	}

	// Convert to response DTO
	response := &socialPayEntity.HostedCheckoutResponseDTO{
		ID:               hostedPayment.ID,
		Amount:           hostedPayment.Amount,
		Currency:         hostedPayment.Currency,
		Description:      hostedPayment.Description,
		Reference:        hostedPayment.Reference,
		SupportedMediums: hostedPayment.SupportedMediums,
		MerchantPaysFee:  hostedPayment.MerchantPaysFee,
		PhoneNumber:      hostedPayment.PhoneNumber,
		SuccessURL:       hostedPayment.SuccessURL,
		FailedURL:        hostedPayment.FailedURL,
		Status:           string(hostedPayment.Status),
		CreatedAt:        hostedPayment.CreatedAt,
		ExpiresAt:        hostedPayment.ExpiresAt,
	}

	uc.log.Info("Successfully retrieved hosted checkout details", map[string]interface{}{
		"hosted_checkout_id": id,
		"amount":             hostedPayment.Amount,
		"currency":           hostedPayment.Currency,
	})

	return response, nil
}

// GetHostedCheckoutWithMerchant retrieves hosted checkout details with merchant information
func (uc *paymentUseCase) GetHostedCheckoutWithMerchant(ctx context.Context, id uuid.UUID) (*socialPayEntity.HostedCheckoutWithMerchantResponseDTO, error) {
	uc.log.Info("Getting hosted checkout details with merchant information", map[string]interface{}{
		"hosted_checkout_id": id,
	})

	// Get hosted payment
	hostedPayment, err := uc.hostedPaymentRepo.GetByID(ctx, id)
	if err != nil {
		uc.log.Error("Failed to get hosted payment", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("hosted payment not found: %w", err)
	}

	// Check if hosted payment has expired
	if time.Now().After(hostedPayment.ExpiresAt) {
		uc.log.Info("Hosted payment has expired", map[string]interface{}{
			"hosted_checkout_id": id,
			"expires_at":         hostedPayment.ExpiresAt,
		})
		return nil, fmt.Errorf("hosted payment has expired")
	}

	// Check if hosted payment is still pending
	if hostedPayment.Status != txEntity.HostedPaymentPending {
		uc.log.Info("Hosted payment is no longer pending", map[string]interface{}{
			"hosted_checkout_id": id,
			"status":             hostedPayment.Status,
		})
		return nil, fmt.Errorf("hosted payment is no longer available")
	}

	// Get merchant information
	merchant, err := uc.merchantUseCase.GetMerchant(ctx, hostedPayment.MerchantID)
	if err != nil {
		uc.log.Error("Failed to get merchant information", map[string]interface{}{
			"error":       err.Error(),
			"merchant_id": hostedPayment.MerchantID,
		})
		// Continue without merchant info rather than failing the whole request
		merchant = nil
	}

	// Convert to response DTO
	response := &socialPayEntity.HostedCheckoutWithMerchantResponseDTO{
		ID:               hostedPayment.ID,
		Amount:           hostedPayment.Amount,
		Currency:         hostedPayment.Currency,
		Description:      hostedPayment.Description,
		Reference:        hostedPayment.Reference,
		SupportedMediums: hostedPayment.SupportedMediums,
		PhoneNumber:      hostedPayment.PhoneNumber,
		SuccessURL:       hostedPayment.SuccessURL,
		FailedURL:        hostedPayment.FailedURL,
		Status:           string(hostedPayment.Status),
		CreatedAt:        hostedPayment.CreatedAt,
		ExpiresAt:        hostedPayment.ExpiresAt,
		MerchantPaysFee:  hostedPayment.MerchantPaysFee,
		AcceptTip:        hostedPayment.AcceptTip,
		Merchant: &merchantEntity.Merchant{
			ID:           merchant.ID,
			LegalName:    merchant.LegalName,
			TradingName:  merchant.TradingName,
			BusinessType: merchant.BusinessType,
			WebsiteURL:   merchant.WebsiteURL,
			Status:       merchant.Status,
			CreatedAt:    merchant.CreatedAt,
			UpdatedAt:    merchant.UpdatedAt,
		},
	}

	uc.log.Info("Successfully retrieved hosted checkout details with merchant information", map[string]interface{}{
		"hosted_checkout_id": id,
		"amount":             hostedPayment.Amount,
		"currency":           hostedPayment.Currency,
		"merchant_id":        hostedPayment.MerchantID,
	})

	return response, nil
}

// ProcessCheckoutPayment processes payment from hosted checkout page
func (uc *paymentUseCase) ProcessCheckoutPayment(ctx context.Context, req *socialPayEntity.CheckoutPaymentRequest) (*socialPayEntity.PaymentResponse, error) {
	uc.log.Info("Starting checkout payment processing", map[string]interface{}{
		"hosted_checkout_id": req.HostedCheckoutID,
		"medium":             req.Medium,
		"phone_number":       req.PhoneNumber,
	})

	// Get hosted payment
	hostedPayment, err := uc.hostedPaymentRepo.GetByID(ctx, req.HostedCheckoutID)
	if err != nil {
		uc.log.Error("Failed to get hosted payment", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("hosted payment not found: %w", err)
	}

	// Check if hosted payment is still valid
	if hostedPayment.Status != txEntity.HostedPaymentPending {
		uc.log.Error("Hosted payment is no longer pending", map[string]interface{}{
			"hosted_checkout_id": req.HostedCheckoutID,
			"status":             hostedPayment.Status,
		})
		return nil, fmt.Errorf("hosted payment is no longer pending")
	}

	if time.Now().UTC().After(hostedPayment.ExpiresAt) {
		uc.log.Error("Hosted payment has expired", map[string]interface{}{
			"hosted_checkout_id": req.HostedCheckoutID,
			"expires_at":         hostedPayment.ExpiresAt,
			"current_time":       time.Now().UTC(),
		})
		return nil, fmt.Errorf("hosted payment has expired")
	}

	// Validate that the selected medium is supported
	mediumSupported := false
	for _, supportedMedium := range hostedPayment.SupportedMediums {
		if supportedMedium == req.Medium {
			mediumSupported = true
			break
		}
	}
	if !mediumSupported {
		return nil, fmt.Errorf("selected payment medium is not supported")
	}

	// Use unified transaction creation service for checkout
	txCreationReq := TransactionCreationRequest{
		UserID:          hostedPayment.UserID,
		MerchantID:      hostedPayment.MerchantID,
		BaseAmount:      hostedPayment.Amount,
		Description:     hostedPayment.Description,
		Medium:          req.Medium,
		Type:            txEntity.DEPOSIT,
		PaymentType:     "checkout",
		MerchantPaysFee: hostedPayment.MerchantPaysFee,
		CallbackURL:     hostedPayment.CallbackURL,
		SuccessURL:      hostedPayment.SuccessURL,
		FailedURL:       hostedPayment.FailedURL,
		TipAmount:       req.TipAmount,
		TipeePhone:      req.TipeePhone,
		TipMedium:       req.TipMedium,
	}

	txCreationResp, txErr := uc.transactionCreationService.CreateTransaction(ctx, txCreationReq)
	if txErr != nil {
		uc.log.Error("Failed to create checkout transaction using unified service", map[string]interface{}{
			"error": txErr.Error(),
		})
		return nil, fmt.Errorf("failed to create checkout transaction: %w", txErr)
	}

	tx := txCreationResp.Transaction
	tx.PhoneNumber = req.PhoneNumber
	tx.Currency = hostedPayment.Currency
	tx.Reference = hostedPayment.Reference
	tx.Status = txEntity.INITIATED
	tx.HasTip = req.TipAmount != nil && *req.TipAmount > 0

	uc.log.Info("Created checkout transaction using unified service", map[string]interface{}{
		"transaction_id":    tx.Id,
		"original_amount":   tx.BaseAmount,
		"total_amount":      tx.TotalAmount,
		"commission_rate":   txCreationResp.CommissionRate,
		"merchant_pays_fee": hostedPayment.MerchantPaysFee,
	})

	// Store transaction
	if err := uc.transactionRepo.Create(ctx, tx); err != nil {
		uc.log.Error("Failed to create transaction", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Process payment using payment service
	paymentReq := &payment.PaymentRequest{
		TransactionID: tx.Id,
		Medium:        req.Medium,
		Amount:        tx.CustomerNet,
		Currency:      tx.Currency,
		PhoneNumber:   tx.PhoneNumber,
		Reference:     tx.Reference,
		Description:   tx.Description,
		CallbackURL:   tx.CallbackURL,
		SuccessURL:    tx.SuccessURL,
		FailedURL:     tx.FailedURL,
	}

	paymentResp, err := uc.paymentService.ProcessPayment(ctx, hostedPayment.MerchantID.String(), paymentReq)
	if err != nil {
		uc.log.Error("Payment processing failed", map[string]interface{}{
			"error": err.Error(),
		})
		tx.Status = txEntity.TransactionStatus(txEntity.FAILED)
		tx.Comment = err.Error()
		_ = uc.transactionRepo.Update(ctx, tx)
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Update transaction status
	tx.Status = txEntity.TransactionStatus(paymentResp.Status)
	if err := uc.transactionRepo.Update(ctx, tx); err != nil {
		uc.log.Error("Failed to update transaction status", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// loging provider_tx_id
	uc.log.Info("provider_tx_id", map[string]interface{}{
		"provider_tx_id": paymentResp.ProcessorRef,
	})

	param := map[string]interface{}{
		"provider_tx_id": paymentResp.ProcessorRef,
	}
	// Setting provider_id
	if err := uc.transactionRepo.UpdateTransactionWithProviderData(ctx, tx.Id, param); err != nil {
		// loging
		uc.log.Error("error while setting provider_tx_id", map[string]interface{}{
			"err": err,
		})
	}

	// Update hosted payment with transaction details
	if err := uc.hostedPaymentRepo.UpdateWithTransaction(ctx, hostedPayment.ID,
		tx.Id, string(req.Medium), req.PhoneNumber, txEntity.HostedPaymentCompleted); err != nil {
		uc.log.Error("Failed to update hosted payment", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to update hosted payment: %w", err)
	}

	uc.log.Info("Successfully processed checkout payment", map[string]interface{}{
		"transaction_id": tx.Id,
		"status":         tx.Status,
	})

	return &socialPayEntity.PaymentResponse{
		Success:              paymentResp.Success,
		Status:               string(paymentResp.Status),
		Message:              paymentResp.Message,
		PaymentURL:           paymentResp.PaymentURL,
		Reference:            tx.Reference,
		SocialPayTransactionID: tx.Id.String(),
	}, nil
}

func (uc *paymentUseCase) QueryTransactionStatus(ctx context.Context, medium txEntity.TransactionMedium, transactionID string) (*payment.TransactionStatusQueryResponse, error) {
	return uc.paymentService.QueryTransactionStatus(ctx, medium, transactionID)
}

type UseCaseConfig struct {
	TransactionRepo    txRepo.TransactionRepository
	HostedPaymentRepo  txRepo.HostedPaymentRepository
	TransactionUseCase transaction_usecase.TransactionUseCase
	PaymentService     PaymentProcessor
	WalletUseCase      walletUsecase.MerchantWalletUsecase
	MerchantUseCase    v2MerchantUsecase.MerchantUseCase
	CommissionUseCase  commission_usecase.CommissionUseCase
}

func NewPaymentUseCase(config UseCaseConfig) PaymentUseCase {
	logger := logging.NewStdLogger("[socialpay-api]")
	transactionCreationService := NewTransactionCreationService(config.CommissionUseCase, logger)

	return &paymentUseCase{
		transactionRepo:            config.TransactionRepo,
		hostedPaymentRepo:          config.HostedPaymentRepo,
		transactionUseCase:         config.TransactionUseCase,
		walletUseCase:              config.WalletUseCase,
		paymentService:             config.PaymentService,
		merchantUseCase:            config.MerchantUseCase,
		transactionCreationService: transactionCreationService,
		log:                        logger,
	}
}

func RoundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
