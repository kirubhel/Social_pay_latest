# SocialPay QA Test Checklist

## 1. RBAC (Role-Based Access Control) - AuthV2 System
### User Permission Testing
- [ ] **Merchant-Scoped Permissions**
  - [ ] Verify users can have different roles for different merchants (using `merchant_id` in `auth.user_groups`)
  - [ ] Test `CheckUserPermissionForMerchant()` function with valid merchant context
  - [ ] Verify `X-MERCHANT-ID` header requirement for merchant-scoped operations
  - [ ] Test permission inheritance through group membership per merchant

- [ ] **Admin vs Merchant User Types**
  - [ ] Test `USER_TYPE_ADMIN` and `USER_TYPE_SUPER_ADMIN` global permissions
  - [ ] Verify admin users can access resources without merchant context
  - [ ] Test `RequirePermissionForAdmin()` middleware functionality
  - [ ] Verify merchant users require merchant context for operations

- [ ] **Resource-Operation Combinations**
  - [ ] Test specific permissions: `RESOURCE_TRANSACTION:OPERATION_READ`
  - [ ] Test wildcard permissions: `RESOURCE_ALL:OPERATION_ALL`
  - [ ] Test admin permissions: `RESOURCE_ADMIN_ALL:OPERATION_ADMIN_ALL`
  - [ ] Verify permission effect (`allow` vs `deny`)

- [ ] **JWT Token Validation**
  - [ ] Verify JWT contains `merchant_permissions` and `merchant_resource_operations`
  - [ ] Test token with multiple merchant contexts
  - [ ] Validate `groups` field grouped by merchant_id

## 2. Transaction Processing & Commission System
### Commission Calculation Engine
- [ ] **Default Commission Settings**
  - [ ] Test default rate: 2.75% + 0.00 cent (from `admin.settings`)
  - [ ] Verify VAT calculation: 15% of commission
  - [ ] Test commission tiers based on amount ranges

- [ ] **Merchant-Specific Commission**
  - [ ] Test custom merchant commission rates (`merchants.commission_percent`)
  - [ ] Verify `commission_active` flag functionality
  - [ ] Test fallback to default when merchant commission inactive
  - [ ] Validate commission persistence in `CommissionSettings` entity

- [ ] **Commission Calculation Logic**
  - [ ] Verify `feeAmount = amount * (percent / 100.0) + cent`
  - [ ] Test `vatAmount = feeAmount * 0.15`
  - [ ] Validate `adminNet = feeAmount - vatAmount`
  - [ ] Check `totalCommission = feeAmount + vatAmount`

### Transaction Amount Calculations
- [ ] **MerchantPaysFee Scenarios**
  - [ ] **DEPOSIT with MerchantPaysFee=true**:
    - [ ] `totalAmount = baseAmount + tipAmount`
    - [ ] `customerNet = totalAmount`
    - [ ] `merchantNet = baseAmount - feeAmount`
  - [ ] **DEPOSIT with MerchantPaysFee=false**:
    - [ ] `totalAmount = baseAmount + feeAmount + tipAmount`
    - [ ] `customerNet = totalAmount`
    - [ ] `merchantNet = baseAmount`
  - [ ] **WITHDRAWAL with MerchantPaysFee=true**:
    - [ ] `totalAmount = baseAmount + feeAmount`
    - [ ] `merchantNet = totalAmount` (positive - what merchant pays)
    - [ ] `customerNet = baseAmount`
  - [ ] **WITHDRAWAL with MerchantPaysFee=false**:
    - [ ] `totalAmount = baseAmount`
    - [ ] `merchantNet = baseAmount` (positive)
    - [ ] `customerNet = baseAmount - feeAmount`

- [ ] **Transaction Fields Validation**
  - [ ] Verify `base_amount`, `fee_amount`, `vat_amount` calculations
  - [ ] Test `admin_net`, `merchant_net`, `customer_net` values
  - [ ] Validate `total_amount` includes all components
  - [ ] Check `merchant_pays_fee` flag storage

### Tip Processing
- [ ] **Tip Validation**
  - [ ] Test tip amount > 0 requires `tipee_phone` and `tip_medium`
  - [ ] Verify tip inclusion in `total_amount` calculation
  - [ ] Test `has_tip`, `tip_processed` flags
  - [ ] Validate tip transaction creation for successful payments

## 3. Wallet Operations & Balance Management
### Admin Wallet (Single Instance)
- [ ] **Balance Calculations**
  - [ ] Verify admin wallet gets `admin_net` from successful transactions
  - [ ] Test `ProcessDepositSuccess()` - adds admin commission
  - [ ] Test `ProcessWithdrawalSuccess()` - adds admin commission
  - [ ] Validate single admin wallet constraint

- [ ] **Admin Wallet Health Check**
  - [ ] Test `/admin/health/wallet-balance` endpoint
  - [ ] Verify calculation: `SUM(admin_net) WHERE status='SUCCESS'`
  - [ ] Check balance discrepancy detection (< 0.01 tolerance)
  - [ ] Test wallet health report generation

### Merchant Wallet Operations
- [ ] **Deposit Processing**
  - [ ] Test `ProcessDepositSuccess()` atomic operation
  - [ ] Verify merchant balance increase by `merchant_net`
  - [ ] Check concurrent transaction handling

- [ ] **Withdrawal Processing**
  - [ ] Test withdrawal locking: `LockAmountForWithdrawal()`
  - [ ] Verify `ProcessWithdrawalSuccess()` unlocks amount
  - [ ] Test `ProcessWithdrawalFailure()` returns locked amount
  - [ ] Check `locked_amount` vs `amount` fields

- [ ] **Balance Reconciliation**
  - [ ] Test balance calculation: `DEPOSITS(merchant_net) - WITHDRAWALS(merchant_net)`
  - [ ] Verify `/admin/health/wallet-balance` for all merchant wallets
  - [ ] Test recalculation script: `scripts/recalculate_wallet_balances.sql`
  - [ ] Check transaction history consistency

## 4. Webhook System & Notifications
### Webhook Event Processing
- [ ] **Transaction Status Updates**
  - [ ] Test valid status transitions: `PENDING → INITIATED → SUCCESS/FAILED`
  - [ ] Verify invalid transition rejection
  - [ ] Test `HandlePaymentStatusUpdate()` function
  - [ ] Check provider data storage (`provider_tx_id`, `provider_data`)

- [ ] **Webhook Payload Structure**
  - [ ] Verify `WebhookEventMerchant` structure:
    ```json
    {
      "event": "DEPOSIT",
      "socialpayTxnId": "uuid",
      "status": "SUCCESS",
      "amount": "merchant_net_value",
      "callbackUrl": "merchant_callback",
      "providerTxId": "external_ref",
      "merchantId": "uuid",
      "userId": "uuid"
    }
    ```

- [ ] **Kafka Integration**
  - [ ] Test webhook dispatch to `WebhookDispatch` topic
  - [ ] Verify webhook send via `WebhookSend` topic
  - [ ] Test message grouping by merchant_id
  - [ ] Check retry mechanism for failed deliveries

### SMS Notifications
- [ ] **Transaction Notifications**
  - [ ] Test customer SMS for successful transactions
  - [ ] Test merchant SMS notifications
  - [ ] Verify notification templates: `TRANSACTION_SUCCESS`, `TRANSACTION_FAILED`
  - [ ] Check notification delivery status tracking

## 5. Transaction Status & Flow Management
### Status Validation
- [ ] **Valid Transaction Statuses**
  - [ ] Test: `INITIATED`, `PENDING`, `SUCCESS`, `FAILED`, `REFUNDED`, `EXPIRED`, `CANCELED`
  - [ ] Verify status enum consistency across system
  - [ ] Test status transition validation

- [ ] **Transaction Types**
  - [ ] Test `DEPOSIT` transactions (customer to merchant)
  - [ ] Test `WITHDRAWAL` transactions (merchant to customer)
  - [ ] Verify type-specific amount calculations
  - [ ] Check transaction medium support: `MPESA`, `TELEBIRR`, `CBE`, `AWASH`

### Payment Flow Testing
- [ ] **Direct API Payments**
  - [ ] Test `/api/v2/payment/direct` endpoint
  - [ ] Verify transaction creation with `TransactionCreationService`
  - [ ] Test payment processing with different mediums
  - [ ] Check callback URL handling

- [ ] **Hosted Checkout**
  - [ ] Test hosted payment creation and processing
  - [ ] Verify `merchant_pays_fee` flag inheritance
  - [ ] Test payment medium selection
  - [ ] Check redirect URL functionality

- [ ] **QR Payments**
  - [ ] Test QR link generation with fee configuration
  - [ ] Verify QR payment processing
  - [ ] Test QR-specific transaction tracking
  - [ ] Check QR tag functionality: `QR_SHOP_PAYMENT`, `QR_RESTAURANT_PAYMENT`

## 6. API Security & Authentication
### API Key Management
- [ ] **Key Validation**
  - [ ] Test API key authentication for payment endpoints
  - [ ] Verify merchant-specific API key restrictions
  - [ ] Test key revocation and expiration

### JWT Authentication
- [ ] **Token Structure**
  - [ ] Verify JWT contains user info and merchant permissions
  - [ ] Test token expiration and refresh mechanism
  - [ ] Check merchant context in token claims

## 7. Data Integrity & Performance
### Database Consistency
- [ ] **Transaction Atomicity**
  - [ ] Test concurrent transaction processing
  - [ ] Verify wallet updates are atomic with transaction status
  - [ ] Check rollback scenarios on failures

- [ ] **Index Performance**
  - [ ] Test query performance on transaction analytics indexes
  - [ ] Verify merchant_id and status index usage
  - [ ] Check wallet balance query performance

### Error Handling
- [ ] **Network Failures**
  - [ ] Test payment processor timeouts
  - [ ] Verify webhook retry mechanisms
  - [ ] Check transaction state consistency during failures

- [ ] **Validation Errors**
  - [ ] Test invalid amount values
  - [ ] Verify merchant ID validation
  - [ ] Check phone number format validation

## 8. Specific Test Scenarios
### Commission Edge Cases
- [ ] Test zero commission rates
- [ ] Verify negative amount handling
- [ ] Test very large transaction amounts
- [ ] Check decimal precision in calculations

### MerchantPaysFee Edge Cases
- [ ] Test fee calculation with tips
- [ ] Verify fee handling in refund scenarios
- [ ] Test commission calculation when merchant pays fees
- [ ] Check edge cases with zero fees

### Webhook Reliability
- [ ] Test webhook delivery during high load
- [ ] Verify duplicate webhook prevention
- [ ] Test webhook authentication if implemented
- [ ] Check webhook payload size limits

## Environment-Specific Tests
- [ ] **Development Environment**
  - [ ] Test with test transaction flags
  - [ ] Verify sandbox payment processor integration
  - [ ] Check debug logging functionality

- [ ] **Staging Environment**
  - [ ] Test with production-like data volumes
  - [ ] Verify database migration scripts
  - [ ] Check monitoring and alerting

- [ ] **Production Environment**
  - [ ] Verify all endpoints with rate limiting
  - [ ] Test backup and recovery procedures
  - [ ] Check compliance with security requirements

## Testing Notes
- Use the wallet health check endpoint to verify balance consistency after each test
- Monitor Kafka topics for webhook message delivery
- Check transaction logs for proper commission calculations
- Verify all monetary calculations maintain precision to 2 decimal places
- Test with various merchant commission configurations
- Ensure RBAC permissions are tested across different user types and merchant contexts 