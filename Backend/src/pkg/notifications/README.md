# Notifications Service

A clean, extensible notification service for Social Pay that supports SMS, email (future), and in-app notifications (future).

## Architecture

The notification service follows a clean architecture pattern with clear separation of concerns:

### Core Layer
- **Entity**: Contains business entities and enums (`NotificationType`, `NotificationStatus`, `NotificationTemplate`)
- **Domain Logic**: Pure business logic without external dependencies

### Use Case Layer  
- **Interfaces**: Defines contracts for providers and services
- **NotificationService**: Main service for sending notifications
- **TransactionNotifier**: Specialized service for transaction notifications
- **UserService**: Service for fetching user/merchant information from database

### Adapter Layer
- **SMS Provider**: AfroSMS implementation for sending SMS messages
- **Database Integration**: User and merchant data retrieval

## Features

### SMS Notifications
- ✅ **AfroSMS Integration**: Production-ready SMS provider using AfroMessage API
- ✅ **Ethiopian Phone Number Support**: Handles multiple formats (9xxxxxxxx, 09xxxxxxxx, 2519xxxxxxxx, +2519xxxxxxxx)
- ✅ **Transaction Notifications**: Role-specific messages for payers, merchants, and tipees
- ✅ **Template System**: Supports OTP and general message templates

### Transaction Notifications
- ✅ **Multi-Recipient Support**: Automatically sends to payer, merchant, and tipee
- ✅ **Role-Based Messages**: Different message formats based on recipient role
- ✅ **Status-Aware**: Success/failure specific messages
- ✅ **Database Integration**: Fetches user information from auth and merchants tables

### Webhook Integration
- ✅ **Automatic SMS Notifications**: Integrated with webhook consumer to send transaction status updates
- ✅ **Non-Blocking**: Notification failures don't block transaction processing
- ✅ **Comprehensive Logging**: Detailed logs for debugging and monitoring

## Usage

### Basic SMS Sending
```go
// Initialize notification service
notificationService := notifications.NewNotificationService(log)

// Send simple SMS
err := notificationService.SendSMS(ctx, "+251912345678", "Your verification code is 123456")
```

### Transaction Notifications
```go
// Initialize transaction notifier with database
transactionNotifier := notifications.NewTransactionNotifier(db, log)

// Send transaction status notifications to all parties
err := transactionNotifier.NotifyTransactionStatus(ctx, transaction, "SUCCESS")
```

### Webhook Integration
The webhook system automatically sends SMS notifications when transaction status updates are received. This is integrated in `HandlePaymentStatusUpdate` method.

## Supported Providers

### SMS Providers
- **AfroSMS**: Uses AfroMessage API with Bearer token authentication

### Future Providers
- **Email**: Framework ready for email provider implementation
- **In-App**: Framework ready for push notification implementation

## Phone Number Formats

The system accepts and normalizes various Ethiopian phone number formats:
- `9xxxxxxxx` → `2519xxxxxxxx`
- `09xxxxxxxx` → `2519xxxxxxxx` 
- `2519xxxxxxxx` → `2519xxxxxxxx`
- `+2519xxxxxxxx` → `2519xxxxxxxx`

## Message Templates

### Transaction Messages
Role-specific transaction notification messages:

**For Payers/Customers:**
```
Dear [Name],
Payment of [Amount] [Currency] to [Merchant] successful.
Reference: [Reference]
Date: [Date] at [Time]

Winners choose and use Social Pay!
```

**For Merchants:**
```
Dear [Name],
You received [Amount] [Currency] from [Customer].
Reference: [Reference]  
Date: [Date] at [Time]

Social Pay - Your trusted payment partner!
```

**For Tipees:**
```
Dear [Name],
You received a tip of [Amount] [Currency].
Reference: [Reference]
Date: [Date] at [Time]

Thank you for your service!
```

## Error Handling

The service includes comprehensive error handling:
- Invalid phone number formats
- Provider API failures  
- Database connection issues
- Template processing errors

Errors are logged but don't prevent transaction processing from continuing.

## Database Integration

The service integrates with existing Social Pay database tables:
- `auth.users` - User information
- `auth.phones` & `auth.phone_identities` - Phone numbers
- `merchants.merchants` - Merchant information
- `merchants.merchant_additional_info` - Additional merchant contact info

## Configuration

SMS provider configuration is currently embedded in the code. Future versions will support configuration through environment variables or config files.

## Testing

To test the notifications package:
```bash
go build -o /dev/null ./src/pkg/notifications/...
```

To test webhook integration:
```bash  
go build -o /dev/null ./src/pkg/webhook/...
``` 