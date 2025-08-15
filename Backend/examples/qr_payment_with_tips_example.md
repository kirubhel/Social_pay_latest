# QR Payment with Tips Integration Example

This example demonstrates how QR payments with tip functionality integrate with the withdrawal API in the SocialPay system.

## Overview

The SocialPay system now supports:
1. **QR Payments with Tips** - Customers can add tips when paying through QR codes
2. **Automatic Tip Processing** - Tips are automatically processed as withdrawals to service staff
3. **Withdrawal API** - Complete withdrawal functionality for merchants and tip recipients

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   QR Payment    │    │ Tip Processing   │    │ Withdrawal API  │
│                 │───▶│    Service       │───▶│                 │
│ • Customer pays │    │ • Detects tips   │    │ • Processes     │
│ • Adds tip      │    │ • Creates        │    │   withdrawals   │
│ • Transaction   │    │   withdrawals    │    │ • Updates       │
│   recorded      │    │ • Links to orig  │    │   status        │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Database Schema

The system uses enhanced transaction tables with QR and tip context:

```sql
-- Transaction table with QR payment and tip support
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    -- ... existing fields ...
    
    -- QR Payment Context
    transaction_source transaction_source DEFAULT 'DIRECT',
    qr_link_id UUID NULL,
    hosted_checkout_id UUID NULL,
    qr_tag VARCHAR(30) NULL,
    
    -- Tip Information
    has_tip BOOLEAN DEFAULT FALSE,
    tip_amount DECIMAL(20,2) NULL,
    tipee_phone VARCHAR(20) NULL,
    tip_medium VARCHAR(20) NULL,
    tip_transaction_id UUID NULL,
    tip_processed BOOLEAN DEFAULT FALSE
);
```

## Example Flow

### 1. Customer Scans QR Code and Pays with Tip

```http
POST /qr/merchant/payment
Content-Type: application/json
X-API-Key: your-api-key

{
    "qr_link_id": "123e4567-e89b-12d3-a456-426614174000",
    "amount": 100.00,
    "currency": "ETB",
    "medium": "TELEBIRR",
    "phone_number": "+251911234567",
    "has_tip": true,
    "tip_amount": 15.00,
    "tipee_phone": "+251922345678",
    "tip_medium": "TELEBIRR"
}
```

**Response:**
```json
{
    "success": true,
    "status": "PENDING",
    "message": "Payment initiated successfully",
    "payment_url": "https://telebirr.pay/checkout/abc123",
    "reference_id": "QR-PAY-123456",
    "socialpay_transaction_id": "456e7890-e89b-12d3-a456-426614174001"
}
```

### 2. Payment Completes - Transaction Record

The system creates a transaction record with tip information:

```json
{
    "id": "456e7890-e89b-12d3-a456-426614174001",
    "merchant_id": "789e0123-e89b-12d3-a456-426614174002",
    "amount": 100.00,
    "status": "SUCCESS",
    "transaction_source": "QR_PAYMENT",
    "qr_link_id": "123e4567-e89b-12d3-a456-426614174000",
    "qr_tag": "QR_RESTAURANT_PAYMENT",
    "has_tip": true,
    "tip_amount": 15.00,
    "tipee_phone": "+251922345678",
    "tip_medium": "TELEBIRR",
    "tip_processed": false
}
```

### 3. Automatic Tip Processing

The tip processing service automatically detects successful payments with tips:

```go
// Tip Processing Service detects pending tips
func (s *tipProcessingService) ProcessPendingTips(ctx context.Context) error {
    // Get all transactions with pending tips
    transactions, err := s.transactionRepo.GetTransactionsWithPendingTips(ctx)
    
    for _, tx := range transactions {
        // Create withdrawal transaction for tip
        tipTransaction := &entity.Transaction{
            Type:              entity.WITHDRAWAL,
            TransactionSource: entity.WITHDRAWAL_TIP,
            PhoneNumber:       *tx.TipeePhone,
            Amount:            *tx.TipAmount,
            Medium:            entity.TransactionMedium(*tx.TipMedium),
            Reference:         fmt.Sprintf("TIP-%s", uuid.New().String()[:8]),
            Status:           entity.INITIATED,
        }
        
        // Process withdrawal via withdrawal API
        s.paymentService.ProcessWithdrawal(ctx, "system-key", paymentReq)
        
        // Mark original transaction tip as processed
        s.transactionRepo.UpdateTipProcessing(ctx, tx.Id, tipTransaction.Id)
    }
}
```

### 4. Withdrawal API Processes Tip

The tip withdrawal is processed through the standard withdrawal API:

```http
POST /payment/withdrawal
Content-Type: application/json
X-API-Key: system-api-key

{
    "amount": 15.00,
    "currency": "ETB",
    "medium": "TELEBIRR",
    "phone_number": "+251922345678",
    "reference": "TIP-ABC12345",
    "callback_url": "https://api.socialpay.et/webhooks/tip-processed"
}
```

**Response:**
```json
{
    "success": true,
    "status": "PENDING",
    "message": "Withdrawal initiated successfully",
    "reference_id": "TIP-ABC12345",
    "socialpay_transaction_id": "999e8877-e89b-12d3-a456-426614174099"
}
```

### 5. Tip Withdrawal Completion

When the tip withdrawal completes, the system updates both transactions:

1. **Original Payment Transaction:**
   ```json
   {
       "tip_processed": true,
       "tip_transaction_id": "999e8877-e89b-12d3-a456-426614174099"
   }
   ```

2. **Tip Withdrawal Transaction:**
   ```json
   {
       "id": "999e8877-e89b-12d3-a456-426614174099",
       "type": "WITHDRAWAL",
       "transaction_source": "WITHDRAWAL",
       "status": "SUCCESS",
       "amount": 15.00,
       "phone_number": "+251922345678"
   }
   ```

## API Endpoints

### QR Payment with Tips
```http
POST /qr/merchant/payment
```

### Withdrawal API
```http
POST /payment/withdrawal
```

### Get Transaction (includes tip info)
```http
GET /payment/transactions/{id}
```

### Get Transactions by QR Link
```http
GET /qr/links/{qr_link_id}/transactions
```

## Repository Methods

The system provides comprehensive repository methods for QR and tip processing:

```go
type TransactionRepository interface {
    // QR Payment and Tip Processing methods
    CreateWithContext(ctx context.Context, tx *entity.Transaction) error
    UpdateTipProcessing(ctx context.Context, transactionID, tipTransactionID uuid.UUID) error
    GetTransactionsWithPendingTips(ctx context.Context) ([]entity.Transaction, error)
    GetTransactionsByQRLink(ctx context.Context, qrLinkID uuid.UUID, limit, offset int32) ([]entity.Transaction, error)
}
```

## Tip Processing Service

The tip processing service provides automated tip handling:

```go
type TipProcessingService interface {
    // Process all transactions with pending tips
    ProcessPendingTips(ctx context.Context) error
    
    // Process a tip for a specific transaction
    ProcessTipForTransaction(ctx context.Context, transactionID uuid.UUID) error
    
    // Create a withdrawal transaction for processed tips
    CreateTipWithdrawal(ctx context.Context, tipeePhone string, tipAmount float64, 
                       medium entity.TransactionMedium, merchantID uuid.UUID) (*entity.Transaction, error)
}
```

## Benefits

1. **Seamless Integration**: QR payments with tips automatically flow through to withdrawals
2. **Automatic Processing**: No manual intervention required for tip distribution
3. **Audit Trail**: Complete transaction history linking payments to tip withdrawals
4. **Real-time Processing**: Tips are processed as soon as payments complete
5. **Multi-Medium Support**: Works with all supported payment mediums (TELEBIRR, CBE, etc.)

## Error Handling

The system includes comprehensive error handling:

- **Insufficient Funds**: Withdrawal API validates merchant wallet balance
- **Failed Tip Processing**: Original transaction tip status remains unprocessed for retry
- **Network Issues**: Automatic retry mechanisms for failed withdrawals
- **Validation Errors**: Complete request validation before processing

## Monitoring and Logging

All operations are logged for monitoring:

```
[TIP-PROCESSING] Starting to process pending tips
[TIP-PROCESSING] Found transactions with pending tips count=5
[TIP-PROCESSING] Processing tip for transaction transaction_id=456e7890...
[TIP-PROCESSING] Creating tip withdrawal transaction tipee_phone=+251922345678
[TIP-PROCESSING] Successfully processed tip for transaction tip_amount=15.00
```

This integration provides a complete, automated tip processing system that leverages the existing withdrawal infrastructure while maintaining data integrity and providing full audit trails. 