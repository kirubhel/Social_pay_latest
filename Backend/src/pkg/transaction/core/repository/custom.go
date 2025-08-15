package repository

import (
	"context"
	"log"
	"regexp"
	"strings"

	repository "github.com/socialpay/socialpay/src/pkg/transaction/core/repository/generated"
)

var withParameterQuery = `
SELECT 
    id,
    phone_number,
    user_id,
    merchant_id,
    type,
    medium,
    reference,
    comment,
    reference_number,
    description,
    verified,
    status,
    test,
    has_challenge,
	webhook_received,
    ttl,
    created_at,
    updated_at,
    confirm_timestamp,
    base_amount,
    fee_amount,
    admin_net,
    vat_amount,
    merchant_net,
	customer_net,
    total_amount,
    currency,
    details,
    token,
    callback_url,
    success_url,
    failed_url,
	provider_tx_id
FROM public.transactions `

func (r *TransactionRepositoryImpl) GetTransactionWithParameter(clause string, args []interface{}) ([]repository.Transaction, error) {

	Query := withParameterQuery + " " + clause
	rows, err := r.q.QueryContext(context.Background(), Query, args...)
	if err != nil {
		log.Println("Err::query_err::", Query)
		return nil, err
	}

	defer rows.Close()

	var txs []repository.Transaction

	for rows.Next() {

		var i repository.Transaction
		if err := rows.Scan(
			&i.ID,               // 0  UUID
			&i.PhoneNumber,      // 1  string or sql.NullString
			&i.UserID,           // 2  UUID
			&i.MerchantID,       // 3  UUID (nullable, so maybe *uuid.UUID or sql.NullString)
			&i.Type,             // 4  string
			&i.Medium,           // 5  string
			&i.Reference,        // 6  string or sql.NullString
			&i.Comment,          // 7  string or sql.NullString
			&i.ReferenceNumber,  // 8  string or sql.NullString
			&i.Description,      // 9  string or sql.NullString
			&i.Verified,         // 10 bool
			&i.Status,           // 11 string or custom type
			&i.Test,             // 12 bool
			&i.HasChallenge,     // 13 bool
			&i.WebhookReceived,  // 14 bool
			&i.Ttl,              // 15 int64 (BIGINT)
			&i.CreatedAt,        // 16 time.Time
			&i.UpdatedAt,        // 17 time.Time
			&i.ConfirmTimestamp, // 18 time.Time or *time.Time (nullable)
			&i.BaseAmount,       // 19 float64 or decimal type
			&i.FeeAmount,        // 20 float64 or decimal type (nullable)
			&i.AdminNet,         // 21 float64 or decimal type (nullable)
			&i.VatAmount,        // 22 float64 or decimal type (nullable)
			&i.MerchantNet,      // 23 float64 or decimal type (nullable)
			&i.CustomerNet,      // 24 float64 or decimal type (nullable)
			&i.TotalAmount,      // 24 float64 or decimal type (nullable)
			&i.Currency,         // 25 string
			&i.Details,          // 26 []byte / json.RawMessage / custom struct
			&i.Token,            // 27 string or sql.NullString
			&i.CallbackUrl,      // 28 string or sql.NullString
			&i.SuccessUrl,       // 29 string or sql.NullString
			&i.FailedUrl,        // 30 string or sql.NullString
			&i.ProviderTxID,     // 31 string or sql.NullString
		); err != nil {
			return nil, err
		}

		txs = append(txs, i)
	}

	return txs, nil

}

func (q *TransactionRepositoryImpl) CountTransactionWithParameter(ctx context.Context, clause string,
	args []interface{}) (int, error) {

	wclause := StripOrderBy(clause)

	log.Println("clause::", wclause)

	// Build the full count query
	countQuery := "SELECT COUNT(*) FROM public.transactions " + wclause

	var total int
	err := q.q.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		log.Println("Err::count_query_err::", countQuery)

		return 0, err
	}

	log.Println("Count value::", total)

	return total, nil
}

func (q *TransactionRepositoryImpl) CountWithClause(ctx context.Context,
	baseTable string, clause string, args ...interface{}) (int, error) {
	query := "SELECT COUNT(*) FROM " + baseTable + " " + clause

	var total int
	err := q.q.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		log.Println("count query failed:", query, "err:", err)
		return 0, err
	}

	return total, nil
}

// StripOrderBy removes ORDER BY ... before semicolon, keeping the semicolon
func StripOrderBy(clause string) string {
	hasSemicolon := strings.HasSuffix(strings.TrimSpace(clause), ";")

	// Remove the semicolon temporarily
	clause = strings.TrimSuffix(clause, ";")

	// Regex: remove ORDER BY ... at the end (case-insensitive)
	re := regexp.MustCompile(`(?i)\s+ORDER\s+BY\s+[\w\s.,_"()]+$`)
	clause = re.ReplaceAllString(clause, "")

	// Restore semicolon if it was there
	if hasSemicolon {
		clause += ";"
	}
	return clause
}
