package usecase

import (
	"context"
	"fmt"

	"github.com/socialpay/socialpay/src/pkg/shared/logging"
	"github.com/robfig/cron/v3"
)

type CronService struct {
	cron                     *cron.Cron
	transactionStatusChecker *TransactionStatusChecker
	log                      logging.Logger
	ctx                      context.Context
}

func NewCronService(
	transactionStatusChecker *TransactionStatusChecker,
	ctx context.Context,
) *CronService {
	// Create cron with seconds support
	cronInstance := cron.New(cron.WithSeconds())

	return &CronService{
		cron:                     cronInstance,
		transactionStatusChecker: transactionStatusChecker,
		log:                      logging.NewStdLogger("[CRON-SERVICE]"),
		ctx:                      ctx,
	}
}

func (cs *CronService) Start() error {
	cs.log.Info("Starting cron service", map[string]interface{}{})

	// Add transaction status checker job - runs every 30 minutes
	// Cron format: "0 */30 * * * *" = every 30 minutes at 0 seconds
	_, err := cs.cron.AddFunc("0 */30 * * * *", func() {
		cs.log.Info("Running scheduled transaction status check", map[string]interface{}{})

		if err := cs.transactionStatusChecker.CheckPendingCBETransactions(cs.ctx); err != nil {
			cs.log.Error("Transaction status check failed", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			cs.log.Info("Transaction status check completed successfully", map[string]interface{}{})
		}
	})

	if err != nil {
		cs.log.Error("Failed to add transaction status checker job", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to add transaction status checker job: %w", err)
	}

	// Add more cron jobs here in the future
	// Example:
	// _, err = cs.cron.AddFunc("@daily", func() {
	//     // Daily cleanup job
	// })

	// Start the cron scheduler
	cs.cron.Start()
	cs.log.Info("Cron service started successfully", map[string]interface{}{
		"total_jobs": len(cs.cron.Entries()),
	})

	// Run initial transaction status check on startup
	cs.log.Info("Running initial transaction status check on startup", map[string]interface{}{})
	go func() {
		if err := cs.transactionStatusChecker.CheckPendingCBETransactions(cs.ctx); err != nil {
			cs.log.Error("Initial transaction status check failed", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			cs.log.Info("Initial transaction status check completed successfully", map[string]interface{}{})
		}
	}()

	return nil
}

func (cs *CronService) Stop() {
	cs.log.Info("Stopping cron service", map[string]interface{}{})
	cs.cron.Stop()
	cs.log.Info("Cron service stopped", map[string]interface{}{})
}

func (cs *CronService) GetJobCount() int {
	return len(cs.cron.Entries())
}

// GetNextRunTimes returns the next run times for all scheduled jobs
func (cs *CronService) GetNextRunTimes() []map[string]interface{} {
	entries := cs.cron.Entries()
	var runTimes []map[string]interface{}

	for i, entry := range entries {
		runTimes = append(runTimes, map[string]interface{}{
			"job_id":       i,
			"next_run":     entry.Next,
			"previous_run": entry.Prev,
		})
	}

	return runTimes
}
