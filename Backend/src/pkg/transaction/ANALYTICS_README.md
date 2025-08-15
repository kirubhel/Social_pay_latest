# Transaction Analytics System

A comprehensive analytics system for SocialPay transactions that provides aggregated data, filtering capabilities, and chart generation for business intelligence.

## Features

### 📊 Transaction Analytics
- **Transaction Counts & Amounts**: Total transactions and amounts
- **Financial Breakdown**: Merchant net amounts
- **Transaction Type Analysis**: Deposits, withdrawals, tips breakdown
- **Period Comparison**: Percentage changes from previous period

### 📈 Chart Data
- **Time-based Aggregation**: Day, week, month, year
- **Dual Chart Types**: Transaction count or amount charts
- **Statistical Summary**: Min, max, average, total values

### 🔍 Advanced Filtering
- **Date Range**: Start and end date filtering
- **Status Filter**: Success, failed, pending, etc.
- **Transaction Type**: Sale, deposit, withdrawal, etc.
- **Payment Medium**: Telebirr, CBE, M-Pesa, etc.
- **Transaction Source**: QR, direct, hosted checkout
- **Amount Range**: Min/max amount filtering
- **QR Tags**: Filter by specific QR payment tags
- **Merchant Filter**: Admin-level merchant filtering

> **Note**: Source and status breakdowns are achieved through filtering rather than separate analytics fields. This provides more flexibility and cleaner responses.

## API Endpoints

### 1. Transaction Analytics
```http
POST /api/v2/transactions/analytics
```

**Request Body:**
```json
{
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-01-31T23:59:59Z",
  "status": ["SUCCESS", "PENDING"],
  "type": ["SALE", "DEPOSIT"],
  "medium": ["TELEBIRR", "CBE"],
  "source": ["QR_PAYMENT", "DIRECT"],
  "qr_tag": ["QR_SHOP_PAYMENT"],
  "amount_min": 100.0,
  "amount_max": 10000.0,
  "merchant_id": ["uuid1", "uuid2"]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_transactions": 1250,
    "total_amount": 125000.50,
    "total_merchant_net": 118750.48,
    "total_deposits": {
      "count": 300,
      "amount": 30000.00
    },
    "total_withdrawals": {
      "count": 150,
      "amount": 15000.00
    },
    "total_tips": {
      "count": 50,
      "amount": 2500.00
    },
    "period_comparison": {
      "transaction_count_change": 15.5,
      "amount_change": 12.3,
      "merchant_net_change": 11.8
    }
  }
}
```

### 2. Chart Data
```http
POST /api/v2/transactions/chart
```

**Request Body:**
```json
{
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-01-31T23:59:59Z",
  "date_unit": "day",
  "chart_type": "amount",
  "status": ["SUCCESS"],
  "type": ["SALE"],
  "medium": ["TELEBIRR"]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "chart_type": "amount",
    "date_unit": "day",
    "data": [
      {
        "date": "2024-01-01T00:00:00Z",
        "value": 5000.50,
        "label": "2024-01-01"
      },
      {
        "date": "2024-01-02T00:00:00Z",
        "value": 7500.75,
        "label": "2024-01-02"
      }
    ],
    "summary": {
      "total_value": 125000.50,
      "average_value": 4032.27,
      "max_value": 8500.00,
      "min_value": 1200.00,
      "data_points": 31
    }
  }
}
```

## Filter Parameters

### Date Filters
- `start_date` (required): Start date in ISO 8601 format
- `end_date` (required): End date in ISO 8601 format

### Transaction Filters
- `status`: Array of transaction statuses
  - `INITIATED`, `PENDING`, `SUCCESS`, `FAILED`, `REFUNDED`, `EXPIRED`, `CANCELED`
- `type`: Array of transaction types
  - `SALE`, `DEPOSIT`, `WITHDRAWAL`, `REFUND`, `SETTLEMENT`, `P2P`, `BILL_PAYMENT`
- `medium`: Array of payment mediums
  - `TELEBIRR`, `CBE`, `MPESA`, `CYBERSOURCE`, `SOCIALPAY`
- `source`: Array of transaction sources
  - `DIRECT`, `QR_PAYMENT`, `HOSTED_CHECKOUT`, `WITHDRAWAL_TIP`

### Amount Filters
- `amount_min`: Minimum transaction amount
- `amount_max`: Maximum transaction amount

### QR & Merchant Filters
- `qr_tag`: Array of QR payment tags
- `merchant_id`: Array of merchant UUIDs (admin only)

### Chart-Specific Filters
- `date_unit`: Time aggregation unit
  - `day`, `week`, `month`, `year`
- `chart_type`: Type of chart data
  - `amount`: Transaction amounts
  - `count`: Transaction counts

## Performance Optimizations

### Database Optimizations
- **Single Query Aggregation**: All analytics calculated in one optimized SQL query
- **Indexed Columns**: Proper indexing on `created_at`, `user_id`, `status`, `type`, `medium`
- **Conditional Aggregation**: Uses `CASE WHEN` for efficient categorization
- **Date Truncation**: PostgreSQL `DATE_TRUNC` for efficient time-based grouping

### Query Structure
```sql
-- Example analytics query structure
SELECT 
    COUNT(*) as total_transactions,
    COALESCE(SUM(amount), 0) as total_amount,
    COUNT(CASE WHEN type = 'DEPOSIT' THEN 1 END) as deposit_count,
    COALESCE(SUM(CASE WHEN type = 'DEPOSIT' THEN amount ELSE 0 END), 0) as deposit_amount,
    -- ... more aggregations
FROM transactions 
WHERE user_id = $1 
    AND created_at >= $2 
    AND created_at <= $3
    -- ... additional filters
```

## Architecture

### Clean Architecture Layers

#### 1. Entity Layer (`core/entity/analytics.go`)
- `AnalyticsFilter`: Filter parameters for analytics
- `ChartFilter`: Filter parameters for chart data
- `TransactionAnalytics`: Aggregated analytics response
- `ChartData`: Chart data response with summary

#### 2. Repository Layer (`core/repository/`)
- `GetTransactionAnalytics()`: Aggregated analytics query
- `GetChartData()`: Time-series chart data query
- Optimized SQL with proper indexing

#### 3. Use Case Layer (`usecase/`)
- `GetTransactionAnalytics()`: Business logic and validation
- `GetChartData()`: Chart data processing
- User context handling and authorization

#### 4. Controller Layer (`adapter/controller/gin/`)
- `GetTransactionAnalytics()`: HTTP handler for analytics
- `GetChartData()`: HTTP handler for chart data
- Request validation and response formatting

## Usage Examples

### Basic Analytics
```bash
curl -X POST "https://api.socialpay.co/api/v2/transactions/analytics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z"
  }'
```

### Filtered Analytics
```bash
curl -X POST "https://api.socialpay.co/api/v2/transactions/analytics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "status": ["SUCCESS"],
    "type": ["SALE"],
    "medium": ["TELEBIRR"],
    "amount_min": 100.0,
    "amount_max": 10000.0
  }'
```

### Get Specific Breakdowns Using Filters

#### Successful Transactions Only
```bash
curl -X POST "https://api.socialpay.co/api/v2/transactions/analytics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "status": ["SUCCESS"]
  }'
```

#### QR Payments Only
```bash
curl -X POST "https://api.socialpay.co/api/v2/transactions/analytics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "source": ["QR_PAYMENT"]
  }'
```

#### Sales Transactions Only
```bash
curl -X POST "https://api.socialpay.co/api/v2/transactions/analytics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "type": ["SALE"]
  }'
```

#### Deposits Only
```bash
curl -X POST "https://api.socialpay.co/api/v2/transactions/analytics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "type": ["DEPOSIT"]
  }'
```

#### Withdrawals Only
```bash
curl -X POST "https://api.socialpay.co/api/v2/transactions/analytics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "type": ["WITHDRAWAL"]
  }'
```

### Daily Amount Chart
```bash
curl -X POST "https://api.socialpay.co/api/v2/transactions/chart" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "date_unit": "day",
    "chart_type": "amount"
  }'
```

### Monthly Transaction Count Chart
```bash
curl -X POST "https://api.socialpay.co/api/v2/transactions/chart" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T23:59:59Z",
    "date_unit": "month",
    "chart_type": "count"
  }'
```

## Error Handling

### Validation Errors (400)
```json
{
  "success": false,
  "error": {
    "type": "validation_bad_input",
    "message": "end date must be after start date"
  }
}
```

### Authentication Errors (401)
```json
{
  "success": false,
  "error": {
    "type": "unauthorized",
    "message": "user ID not found in context"
  }
}
```

### Server Errors (500)
```json
{
  "success": false,
  "error": {
    "type": "db_read_failed",
    "message": "failed to get transaction analytics"
  }
}
```

## Security

- **JWT Authentication**: All endpoints require valid JWT tokens
- **User Isolation**: Analytics filtered by authenticated user ID
- **Input Validation**: Comprehensive validation of all filter parameters
- **SQL Injection Protection**: Parameterized queries with proper escaping

## Future Enhancements

- **Real-time Analytics**: WebSocket-based live analytics updates
- **Merchant Comparison**: Cross-merchant analytics for admin users
- **Export Functionality**: PDF/Excel export of analytics data
- **Caching Layer**: Redis caching for frequently accessed analytics
- **Advanced Visualizations**: More chart types and visualization options
- **Scheduled Reports**: Automated analytics reports via email
- **Machine Learning**: Predictive analytics and trend analysis

## Contributing

When adding new analytics features:

1. **Add Entity Fields**: Update `analytics.go` with new data structures
2. **Update Repository**: Add new aggregation queries in repository layer
3. **Extend Use Cases**: Add business logic in use case layer
4. **Update Controllers**: Add new endpoints in controller layer
5. **Add Tests**: Comprehensive unit and integration tests
6. **Update Documentation**: Keep this README updated

## Performance Monitoring

Monitor these metrics for optimal performance:

- **Query Execution Time**: Analytics queries should complete under 500ms
- **Database Load**: Monitor CPU and memory usage during analytics queries
- **Cache Hit Rate**: If caching is implemented, monitor cache effectiveness
- **API Response Time**: End-to-end response time should be under 1 second

---

**Built with ❤️ for SocialPay Analytics**