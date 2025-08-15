package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/wallet/adapter/gateway/repository"
	"github.com/socialpay/socialpay/src/pkg/wallet/core/entity"
)

type MerchantWalletUsecase struct {
	walletRepository repository.WalletRepository
	logger           logging.Logger
}

func NewMerchantWalletUsecase(walletRepository repository.WalletRepository, logger logging.Logger) MerchantWalletUsecase {
	return MerchantWalletUsecase{
		walletRepository: walletRepository,
		logger:           logger,
	}
}

func (u *MerchantWalletUsecase) UpdateMerchantWallet(ctx context.Context, walletID uuid.UUID, amount float64, lockedAmount float64) error {
	return u.walletRepository.UpdateMerchantWallet(ctx, walletID, amount, lockedAmount)
}

func (u *MerchantWalletUsecase) CreateMerchantWallet(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID, amount float64, lockedAmount float64, currency string) error {
	return u.walletRepository.CreateMerchantWallet(ctx, userID, merchantID, amount, lockedAmount, currency)
}

func (u *MerchantWalletUsecase) GetMerchantWallet(ctx context.Context, merchantID uuid.UUID) (*entity.MerchantWallet, error) {
	return u.walletRepository.GetMerchantWalletByMerchantID(ctx, merchantID)
}

// LockWithdrawalAmount locks the specified amount in the merchant wallet for withdrawal
// This prevents the amount from being used for other withdrawals while the transaction is processing
func (u *MerchantWalletUsecase) LockWithdrawalAmount(ctx context.Context, merchantID uuid.UUID, amount float64) error {
	// Create a timeout context to prevent indefinite hangs
	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Use the atomic locking operation
	err := u.walletRepository.LockWithdrawalAmountAtomic(txCtx, merchantID, amount)
	if err != nil {
		return fmt.Errorf("failed to lock withdrawal amount: %w", err)
	}

	u.logger.Debug("Withdrawal amount locked successfully", map[string]interface{}{
		"merchantID": merchantID,
		"amount":     amount,
	})

	return nil
}

// ProcessTransactionStatus handles the final transaction status (success or failure) for both deposits and withdrawals
// This is a high-performance implementation using atomic SQL operations
// For withdrawals:
//   - Success: unlocks the locked amount (amount already deducted during lock phase)
//   - Failure: returns locked amount to available balance and unlocks it
//
// For deposits:
//   - Success: increases the available balance by the specified amount
//   - Failure: no action needed (no prior locking was done)
//
// Admin wallet always gets commission when transaction is successful
func (u *MerchantWalletUsecase) ProcessTransactionStatus(ctx context.Context, merchantID uuid.UUID, merchantAmount float64, adminAmount float64, isSuccess bool, isWithdrawal bool) error {
	u.logger.Info("Processing transaction status", map[string]interface{}{
		"merchantID":   merchantID,
		"amount":       merchantAmount,
		"adminAmount":  adminAmount,
		"isSuccess":    isSuccess,
		"isWithdrawal": isWithdrawal,
	})

	// Single atomic SQL operation - no Go logic, no fetching, no complex transactions
	if isWithdrawal {
		if isSuccess {
			// Withdrawal success: unlock amount (don't change available balance) + admin commission
			err := u.walletRepository.ProcessWithdrawalSuccess(ctx, merchantID, merchantAmount, adminAmount)
			if err != nil {
				u.logger.Error("Failed to process withdrawal success", map[string]interface{}{
					"error":      err,
					"merchantID": merchantID,
					"amount":     merchantAmount,
				})
				return fmt.Errorf("failed to process withdrawal success: %w", err)
			}
		} else {
			// Withdrawal failure: return locked amount to available balance
			err := u.walletRepository.ProcessWithdrawalFailure(ctx, merchantID, merchantAmount)
			if err != nil {
				u.logger.Error("Failed to process withdrawal failure", map[string]interface{}{
					"error":      err,
					"merchantID": merchantID,
					"amount":     merchantAmount,
				})
				return fmt.Errorf("failed to process withdrawal failure: %w", err)
			}
		}
	} else {
		if isSuccess {
			// Deposit success: add to available balance + admin commission
			err := u.walletRepository.ProcessDepositSuccess(ctx, merchantID, merchantAmount, adminAmount)
			if err != nil {
				u.logger.Error("Failed to process deposit success", map[string]interface{}{
					"error":      err,
					"merchantID": merchantID,
					"amount":     merchantAmount,
				})
				return fmt.Errorf("failed to process deposit success: %w", err)
			}
		}
		// Deposit failure: no-op (no funds were locked)
	}

	u.logger.Info("Transaction status processed successfully", map[string]interface{}{
		"merchantID":   merchantID,
		"amount":       merchantAmount,
		"isSuccess":    isSuccess,
		"isWithdrawal": isWithdrawal,
	})

	return nil
}
