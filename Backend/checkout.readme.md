# SocialPay Hosted Checkout API Documentation

## Overview

The SocialPay Hosted Checkout API provides a secure and flexible way for merchants to accept payments without handling sensitive payment information directly. Instead of integrating payment forms into their own applications, merchants can redirect customers to a SocialPay-hosted payment page.

## Key Benefits

- **Flexibility**: Support for multiple payment methods (MPESA, TELEBIRR, CBE)
- **Simplicity**: Minimal integration effort required
- **Customization**: Configurable redirect URLs and payment options

---

## 1. Creating a Hosted Checkout Session

### Create Hosted Checkout

Creates a new hosted checkout session and returns a payment URL where customers can complete their payment.

**Endpoint**: `POST /api/v2/payment/checkout`

**Authentication**: Required (API Key in header)

**Headers**:
```
Content-Type: application/json
X-API-Key: your-api-key
```

**Request Body**:
```json
{
  "amount": 1000.50,
  "currency": "ETB",
  "description": "Payment for Order #12345",
  "reference": "ORDER_12345_1699123456",
  "supported_mediums": ["MPESA", "TELEBIRR", "CBE"],
  "phone_number": "251911234567",
  "redirects": {
    "success": "https://yourstore.com/payment/success",
    "failed": "https://yourstore.com/payment/failed"
  },
  "callback_url": "https://yourstore.com/webhooks/socialpay"
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `amount` | number | Yes | Payment amount (must be > 0.01) |
| `currency` | string | Yes | 3-letter currency code (e.g., "ETB") |
| `description` | string | Yes | Payment description for customer |
| `reference` | string | Yes | Unique merchant reference (1-100 chars) |
| `supported_mediums` | array | Yes | List of payment methods to offer |
| `phone_number` | string | No | Pre-filled customer phone number |
| `redirects.success` | string | Yes | URL to redirect after successful payment |
| `redirects.failed` | string | Yes | URL to redirect after failed payment |
| `callback_url` | string | No | Webhook URL for payment notifications |

**Supported Payment Mediums**:
- `MPESA` - M-Pesa mobile money
- `TELEBIRR` - TeleBirr mobile money  
- `CBE` - Commercial Bank of Ethiopia
- `CYBERSOURCE` - Credit/debit cards

**Response**:
```json
{
  "success": true,
  "status": "PENDING",
  "message": "Hosted checkout created successfully",
  "payment_url": "https://checkout.socialpay.com/checkout/123e4567-e89b-12d3-a456-426614174000",
  "reference_id": "ORDER_12345_1699123456",
  "socialpay_transaction_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

**Response Fields**:

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Whether the operation was successful |
| `status` | string | Current status ("PENDING") |
| `message` | string | Human-readable status message |
| `payment_url` | string | URL where customer should complete payment |
| `reference_id` | string | Echo of the merchant reference |
| `socialpay_transaction_id` | string | SocialPay's internal checkout ID |

**Error Responses**:

| Status Code | Error | Description |
|-------------|-------|-------------|
| 400 | Bad Request | Invalid request parameters or validation errors |
| 401 | Unauthorized | Missing or invalid API key |
| 403 | Forbidden | Insufficient permissions |
| 409 | Conflict | Reference ID already exists |
| 500 | Internal Server Error | Server-side error |

**Example Error Response**:
```json
{
  "success": false,
  "message": "Reference ID already exists for this merchant"
}
```

---

## 2. Redirecting Customers

### Customer Payment Flow

After creating a hosted checkout session, redirect your customer to the `payment_url` returned in the response.

**Redirect Process**:

1. **Create Checkout**: Call the creation endpoint to get a `payment_url`
2. **Redirect Customer**: Send customer to the `payment_url`
3. **Customer Pays**: Customer completes payment on SocialPay's secure page
4. **Automatic Redirect**: Customer is automatically redirected to your success/failed URL
5. **Webhook Notification**: You receive payment status via webhook (if configured)

**Example Redirect**:
```javascript
// After creating checkout and getting response
const checkoutResponse = await createHostedCheckout(paymentData);
if (checkoutResponse.success) {
  // Redirect customer to payment page
  window.location.href = checkoutResponse.payment_url;
}
```

**Customer Experience**:
- Customer is taken to a secure SocialPay payment page
- They can select from the payment methods you specified
- They enter their phone number (if not pre-filled)
- They complete payment through their chosen provider
- They are automatically redirected back to your website

---

## 3. Handling Payment Results

### Redirect URLs

Configure your success and failed URLs to handle customers returning from payment:

**Success URL**: `https://yourstore.com/payment/success`
- Customer is redirected here after successful payment
- You can display a success message and order confirmation
- Payment confirmation will also come via webhook

**Failed URL**: `https://yourstore.com/payment/failed`
- Customer is redirected here if payment fails
- You can display an error message and retry options
- Consider offering alternative payment methods

**Webhook Security**:
- Use HTTPS for your webhook endpoint
- Implement idempotency to handle duplicate notifications

---

## Integration Flow

### Complete Integration Steps

1. **Create Checkout Session**
   ```
   POST /api/v2/payment/checkout
   → Returns payment_url
   ```

2. **Redirect Customer**
   ```
   Redirect to payment_url
   → Customer completes payment
   ```

3. **Handle Return**
   ```
   Customer redirected to success/failed URL
   → Display appropriate message
   ```

4. **Process Webhook**
   ```
   Receive webhook notification
   → Update order status in your system
   ```

### Sample Integration Code

**Step 1: Create Checkout**
```bash
curl -X POST https://api.socialpay.com/api/v2/payment/checkout \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "amount": 1000.50,
    "currency": "ETB",
    "description": "Payment for Order #12345",
    "reference": "ORDER_12345_1699123456",
    "supported_mediums": ["MPESA", "TELEBIRR", "CBE"],
    "redirects": {
      "success": "https://yourstore.com/payment/success",
      "failed": "https://yourstore.com/payment/failed"
    },
    "callback_url": "https://yourstore.com/webhooks/socialpay"
  }'
```

## Fee Structure

**For Deposits**:
- **Transaction Fee**: 2.75% of payment amount
- **VAT**: 15% of the transaction fee
- **Total Customer Pays**: Original Amount + Fee + VAT

**Example Calculation**:
- Order Amount: 1000 ETB
- Transaction Fee: 27.50 ETB (2.75%)
- VAT: 4.13 ETB (15% of fee)
- **Total Charged**: 1031.63 ETB
- **You Receive**: 1000 ETB

---


## Security Best Practices

- **API Keys**: Keep API keys secure and rotate regularly
- **HTTPS**: All endpoints must be accessed over HTTPS
- **Webhook Validation**: Always validate webhook signatures
- **Reference IDs**: Use unique, non-guessable reference IDs
- **Timeout Handling**: Checkout sessions expire after 24 hours

---



## Support

For technical support or questions about the Hosted Checkout API:

- **Documentation**: [https://docs.socialpay.com](https://docs.socialpay.com)
- **Support Email**: support@socialpay.com
- **Developer Portal**: [https://developer.socialpay.com](https://developer.socialpay.com) 