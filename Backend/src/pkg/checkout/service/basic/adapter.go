package basic

import (
	"log"

	account "github.com/socialpay/socialpay/src/pkg/account/usecase"
	"github.com/socialpay/socialpay/src/pkg/checkout/service"
)

type BasicCheckoutService struct {
	log     *log.Logger
	repo    service.CheckoutRepo
	account account.Interactor
}

func New(log *log.Logger, repo service.CheckoutRepo, account account.Interactor) service.CheckoutInteractor {
	return BasicCheckoutService{log, repo, account}
}
