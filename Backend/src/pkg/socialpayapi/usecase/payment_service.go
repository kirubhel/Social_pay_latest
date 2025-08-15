package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"
)

// PaymentProcessor handles payment processing using different processors
type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error)
	ProcessWithdrawal(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error)
	QueryTransactionStatus(ctx context.Context, medium txEntity.TransactionMedium, transactionID string) (*payment.TransactionStatusQueryResponse, error)
}

type paymentService struct {
	processors map[txEntity.TransactionMedium]payment.Processor
	log        logging.Logger
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(processors ...payment.Processor) PaymentProcessor {
	processorMap := make(map[txEntity.TransactionMedium]payment.Processor)
	for _, p := range processors {
		processorMap[p.GetType()] = p
	}
	return &paymentService{
		processors: processorMap,
		log:        logging.NewStdLogger("[SOCIALPAY-API] [PAYMENT-SERVICE]"),
	}
}

func (s *paymentService) ProcessPayment(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {

	processor, ok := s.processors[req.Medium]
	if !ok {
		s.log.Error("No payment processor available", map[string]interface{}{
			"processor_type": req.Medium,
		})
		return nil, fmt.Errorf("no payment processor available")
	}

	// Round amount to atmost 2 decimal places
	req.Amount = math.Round(req.Amount*100) / 100

	return processor.InitiatePayment(ctx, apikey, req)
}

func (s *paymentService) QueryTransactionStatus(ctx context.Context, medium txEntity.TransactionMedium, transactionID string) (*payment.TransactionStatusQueryResponse, error) {
	processor, ok := s.processors[medium]
	if !ok {
		return nil, fmt.Errorf("no payment processor available")
	}

	return processor.QueryTransactionStatus(ctx, transactionID)
}

func (s *paymentService) ProcessWithdrawal(ctx context.Context, apikey string, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	processor, ok := s.processors[req.Medium]
	if !ok {
		return nil, fmt.Errorf("no payment processor available")
	}

	// Round amount to atmost 2 decimal places
	req.Amount = math.Round(req.Amount*100) / 100

	return processor.InitiateWithdrawal(ctx, apikey, req)
}
