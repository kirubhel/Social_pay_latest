package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/webhook/adapter/gateway/repository"
	"github.com/socialpay/socialpay/src/pkg/webhook/core/entity"
)

type WalletUsecase struct {
	walletRepository repository.WalletRepository
}

func NewWalletUsecase(walletRepository repository.WalletRepository) WalletUsecase {
	return WalletUsecase{walletRepository: walletRepository}
}

func (u *WalletUsecase) UpdateMerchantWallet(ctx context.Context, walletID uuid.UUID, amount float64, lockedAmount float64) error {
	return u.walletRepository.UpdateMerchantWallet(ctx, walletID, amount, lockedAmount)
}

func (u *WalletUsecase) CreateMerchantWallet(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID, amount float64, lockedAmount float64, currency string) error {
	return u.walletRepository.CreateMerchantWallet(ctx, userID, merchantID, amount, lockedAmount, currency)
}

func (u *WalletUsecase) GetMerchantWallet(ctx context.Context, merchantID uuid.UUID) (*entity.Wallet, error) {
	return u.walletRepository.GetMerchantWallet(ctx, merchantID)
}

// LockWithdrawalAmount locks the specified amount in the merchant wallet for withdrawal
// This prevents the amount from being used for other withdrawals while the transaction is processing
func (u *WalletUsecase) LockWithdrawalAmount(ctx context.Context, merchantID uuid.UUID, amount float64) error {
	// Start a database transaction
	tx, err := u.walletRepository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure we either commit or rollback
	defer func() {
		if err != nil {
			if rbErr := u.walletRepository.RollbackTx(tx); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Get the wallet with exclusive lock
	wallet, err := u.walletRepository.GetMerchantWalletForUpdate(ctx, tx, merchantID)
	if err != nil {
		return fmt.Errorf("failed to get merchant wallet: %w", err)
	}

	// Check if wallet has sufficient funds
	if wallet.Amount < amount {
		if rbErr := u.walletRepository.RollbackTx(tx); rbErr != nil {
			log.Printf("Error rolling back transaction: %v", rbErr)
		}
		return fmt.Errorf("insufficient funds: available balance is %.2f %s", wallet.Amount, wallet.Currency)
	}

	// Lock the amount (reduce available balance, increase locked amount)
	wallet.Amount -= amount
	wallet.LockedAmount += amount

	// Update the wallet within the transaction
	err = u.walletRepository.UpdateMerchantWalletWithTx(ctx, tx, wallet.ID, wallet.Amount, wallet.LockedAmount)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	// Commit the transaction
	if err = u.walletRepository.CommitTx(tx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ProcessTransactionStatus handles the final transaction status (success or failure) for both deposits and withdrawals
// For withdrawals:
//
//	Success: keeps the amount deducted from available balance and removes it from locked amount
//	Failure: returns the amount to available balance and removes it from locked amount
//
// For deposits:
//
//	Success: increases the available balance by the specified amount
//	Failure: no changes to the wallet (assuming no prior locking was done)
func (u *WalletUsecase) ProcessTransactionStatus(ctx context.Context, merchantID uuid.UUID, amount float64, isSuccess bool, isWithdrawal bool) error {
	// Start a database transaction
	tx, err := u.walletRepository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure we either commit or rollback
	defer func() {
		if err != nil {
			if rbErr := u.walletRepository.RollbackTx(tx); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Get the wallet with exclusive lock
	wallet, err := u.walletRepository.GetMerchantWalletForUpdate(ctx, tx, merchantID)
	if err != nil {
		return fmt.Errorf("failed to get merchant wallet: %w", err)
	}

	// Update wallet based on transaction type and status
	if isWithdrawal {
		// Withdrawal logic
		if isSuccess {
			// Success: amount stays deducted from available balance, but remove from locked
			wallet.LockedAmount -= amount
		} else {
			// Failure: return amount to available balance and remove from locked
			wallet.Amount += amount
			wallet.LockedAmount -= amount
		}
	} else {
		// Deposit logic
		if isSuccess {
			// Success: increase available balance
			wallet.Amount += amount
		}
		// For deposit failure, no action needed as funds haven't been locked
	}

	// Update the wallet within the transaction
	err = u.walletRepository.UpdateMerchantWalletWithTx(ctx, tx, wallet.ID, wallet.Amount, wallet.LockedAmount)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	// Commit the transaction
	if err = u.walletRepository.CommitTx(tx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ProcessWithdrawalStatus is kept for backward compatibility
// Deprecated: Use ProcessTransactionStatus instead
func (u *WalletUsecase) ProcessWithdrawalStatus(ctx context.Context, merchantID uuid.UUID, amount float64, isSuccess bool) error {
	return u.ProcessTransactionStatus(ctx, merchantID, amount, isSuccess, true)
}
