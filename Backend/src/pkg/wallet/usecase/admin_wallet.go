package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/wallet/adapter/gateway/repository"
	"github.com/socialpay/socialpay/src/pkg/wallet/core/entity"
)

type AdminWalletUsecase struct {
	merchantWalletRepository repository.WalletRepository
}

func NewAdminWalletUsecase(merchantWalletRepository repository.WalletRepository) AdminWalletUsecase {
	return AdminWalletUsecase{merchantWalletRepository: merchantWalletRepository}
}

// GetAdminWallet retrieves the single admin wallet for the given user
// Note: There is only one admin wallet in the system
func (u *AdminWalletUsecase) GetAdminWallet(ctx context.Context) (*entity.MerchantWallet, error) {
	return u.merchantWalletRepository.GetAdminWallet(ctx)
}

// GetTotalAdminWalletAmount gets the total amount across all admin wallets
// Note: Since there's only one admin wallet, this returns the balance of that single wallet
func (u *AdminWalletUsecase) GetTotalAdminWalletAmount(ctx context.Context) (map[string]float64, error) {
	return u.merchantWalletRepository.GetTotalAdminWalletAmount(ctx)
}

// CheckWalletBalanceHealth verifies if wallet balances match transaction history
// This includes checking the single admin wallet and all merchant wallets
func (u *AdminWalletUsecase) CheckWalletBalanceHealth(ctx context.Context, userID uuid.UUID) (*entity.WalletHealthCheck, error) {
	return u.merchantWalletRepository.CheckWalletBalanceHealth(ctx)
}

// TODO: flag as not ready for use
// needs to be coupled with the merchant processing - on the same function
func (u *AdminWalletUsecase) ProcessTransactionStatus(ctx context.Context, userID uuid.UUID, amount float64, isSuccess bool, isCredit bool) error {
	return fmt.Errorf("not implemented yet")
}
