# Cash Flow Payment Gateway API

A comprehensive payment gateway API built with Go, PostgreSQL, and RabbitMQ. Designed for merchants to process payments with automatic fee deduction (1%), callback notifications, and balance management.

##  Quick Start

### Prerequisites
- Docker and Docker Compose installed

### Start the Application
```bash
docker compose up 
```

This will start:
- **PostgreSQL 15** on port 5433 (internal: 5432)
- **RabbitMQ** on ports 5673 (AMQP) and 15673 (Management UI)
- **Go Application** on port 3074

### Stop the Application
```bash
docker compose down
```

##  API Documentation

### Swagger UI
Once the application is running, visit:
- **Swagger Documentation**: http://localhost:3074/swagger/index.html

##  Payment Gateway Flow

### 1. Merchant Registration
Merchants create accounts and receive unique API keys for authentication.

### 2. Payment Intent Creation
Merchants create payment intents using their API key. The system only accepts **ETB** and **USD** currencies.

### 3. Asynchronous Processing
Payment intents are processed asynchronously via RabbitMQ workers.

### 4. Fee Deduction & Balance Update
The gateway automatically deducts **1% fee** from each transaction and updates merchant balances.

### 5. Callback Notifications
Merchants receive callback notifications with transaction results.

##  API Endpoints

### Merchant Management

#### Create Merchant Account
```http
POST /cashflow_test/v1/account/create-merchant
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

**Response:**
```json
{
  "status": true,
  "merchant_id": "CASM-ABC123",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "api_key": "cash_test_abc123def456",
  "message": "Merchant created successfully"
}
```

#### Get Merchant Details
```http
GET /cashflow_test/v1/account/merchant?merchant_id=CASM-ABC123
```

**Response:**
```json
{
  "status": true,
  "merchant_id": "CASM-ABC123",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "merchant_status": "active",
  "api_key": "cash_test_abc123def456",
  "api_key_status": "active",
  "created_at": "2024-01-05T10:30:00Z",
  "api_key_created": "2024-01-05T10:30:00Z",
  "balances": [
    {
      "currency": "ETB",
      "available_balance": "99.49",
      "total_deposit": "100.50",
      "total_transaction_count": 1,
      "last_updated": "2024-01-05T10:35:00Z"
    }
  ],
  "transactions": [
    {
      "id": "d49a7dd0-95b9-4636-acf6-d06b87a8e525",
      "payment_intent_id": "PI-ABC123",
      "merchant_id": "CASM-ABC123",
      "amount": "100.50",
      "currency": "ETB",
      "status": "success",
      "third_party_reference": "TP123456789",
      "payment_method": "card",
      "fee_amount": "1.01",
      "account_number": "1234567890",
      "processed_at": "2024-01-05T10:35:00Z",
      "created_at": "2024-01-05T10:30:00Z",
      "updated_at": "2024-01-05T10:35:00Z"
    }
  ],
  "message": "Merchant details retrieved successfully"
}
```

### Payment Processing

#### Create Payment Intent
```http
POST /cashflow_test/v1/checkout/create-intent
X-API-KEY: your_merchant_api_key
Content-Type: application/json

{
  "amount": 100.50,
  "currency": "ETB",
  "description": "Payment for order #123",
  "callback_url": "https://example.com/callback",
  "nonce": "unique_nonce_123456789"
}
```

**Response:**
```json
{
  "status": true,
  "payment_intent_id": "PI-ABC123",
  "amount": 100.5,
  "currency": "ETB",
  "payment_status": "pending",
  "created_at": "2024-01-05T10:30:00Z",
  "expires_at": "2024-01-05T10:45:00Z",
  "message": "Payment intent created successfully"
}
```


### Health Check
```http
GET /health
```

##  Fee Structure

- **Transaction Fee**: 1% of the payment amount
- **Supported Currencies**: ETB (Ethiopian Birr) and USD (US Dollar)
- **Fee Deduction**: Automatic deduction before balance credit

##  Authentication

All payment-related endpoints require API key authentication:

```http
X-API-KEY: your_merchant_api_key
```

##  Database Schema

The system automatically initializes with the following tables:

- **`merchants`** - Merchant account information
- **`merchant_api_keys`** - API key management with status tracking
- **`payment_intents`** - Payment intent records with expiration
- **`payment_transactions`** - Transaction records with fee tracking
- **`merchant_balances`** - Balance management per currency

##  Message Queue

RabbitMQ is used for asynchronous processing:
- **Queue**: `payment_intents_queue`
- **Management UI**: http://localhost:15673 (guest/guest)

## Development

### Rebuild Application
```bash
docker compose up --build -d
```

### View Logs
```bash
# All services
docker compose logs

# Specific service
docker compose logs app
docker compose logs postgres
docker compose logs rabbitmq
```

### Access Services
- **API Server**: http://localhost:3074
- **Swagger UI**: http://localhost:3074/swagger/index.html
- **RabbitMQ Management**: http://localhost:15673 (guest/guest)
- **PostgreSQL**: localhost:5433 (cashflow_user/cashflow_pass, db: cashflow_dev)

##  Troubleshooting

### Port Conflicts
If you have local services running on standard ports, the Docker setup uses:
- PostgreSQL: 5433 (instead of 5432)
- RabbitMQ: 5673 (AMQP), 15673 (Management UI)

### Clean Restart
```bash
docker compose down -v  # Remove volumes too
docker compose up -d
```

##  Example Usage

### 1. Create a Merchant
```bash
curl -X POST http://localhost:3074/cashflow_test/v1/account/create-merchant \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Merchant",
    "email": "merchant@example.com"
  }'
```

### 2. Create Payment Intent
```bash
curl -X POST http://localhost:3074/cashflow_test/v1/checkout/create-intent \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: cash_test_abc123def456" \
  -d '{
    "amount": 100.00,
    "currency": "USD",
    "description": "Test payment",
    "callback_url": "https://example.com/webhook",
    "nonce": "unique_nonce_123"
  }'
```

##  Support

For API integration  refer to the Swagger documentation at `/swagger/index.html` when the application is running.
