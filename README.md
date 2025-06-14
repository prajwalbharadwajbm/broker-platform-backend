# Broker Platform Backend

A robust trading platform backend built with Go that provides RESTful APIs for managing user accounts, stock holdings, trading positions, and order book functionality. The platform supports real-time trading operations with secure authentication and comprehensive data management with circuit breaker for database service.

## Features

- **User Authentication**: JWT-based authentication with refresh tokens
- **Portfolio Management**: Track holdings, positions, and P&L
- **Order Book**: Real-time order matching and execution
- **Secure API**: Comprehensive middleware for security and logging
- **Database**: PostgreSQL with UUID primary keys and proper indexing
- **Logging**: Structured logging with configurable levels
- **Circuit Breaker**: Fault tolerance for database service

## Tech Stack

- **Database**: PostgreSQL with pgcrypto extension
- **Authentication**: JWT tokens
- **Router**: HttpRouter for fast HTTP routing
- **Logging**: Zerolog for structured logging
- **Security**: bcrypt for password hashing

## Prerequisites

- Go 1.23+ installed
- PostgreSQL 12+ installed and running

## Environment Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/prajwalbharadwajbm/broker.git
   cd broker
   ```

2. **Create environment file**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your configuration:
   ```env
   # Application Configuration
   APP_ENV=dev
   LOG_LEVEL=info
   PORT=8080
   
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   POSTGRES_USER=your_username
   POSTGRES_PASSWORD=your_password
   POSTGRES_DB=broker-platform
   
   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key-here
   ```

3. **Install dependencies**
   ```bash
   go mod download
   or 
   go mod tidy
   ```

## Database Setup

### 1. Create Database and Tables

Connect to PostgreSQL and run the following commands:

```sql
-- Create database
CREATE DATABASE "broker-platform";

-- Connect to the database
\c "broker-platform"

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create refresh tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create holdings table
CREATE TABLE holdings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    quantity NUMERIC(20,8) NOT NULL,
    average_price NUMERIC(20,8) NOT NULL,
    current_price NUMERIC(20,8) NOT NULL,
    total_value NUMERIC(20,8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create positions table
CREATE TABLE positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    position_type VARCHAR(10) NOT NULL CHECK (position_type IN ('LONG', 'SHORT')),
    quantity NUMERIC(20,8) NOT NULL,
    entry_price NUMERIC(20,8) NOT NULL,
    current_price NUMERIC(20,8) NOT NULL,
    unrealized_pnl NUMERIC(20,8) NOT NULL,
    realized_pnl NUMERIC(20,8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create orderbook table
CREATE TABLE orderbook (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(4) NOT NULL CHECK (side IN ('BUY', 'SELL')),
    price NUMERIC(20,8) NOT NULL,
    quantity NUMERIC(20,8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### 2. Insert Mock Data for Testing

```sql
-- Insert test users
-- use POST /api/v1/users/signup to create users
-- curl -X POST http://localhost:8080/api/v1/users/signup -H "Content-Type: application/json" -d '{"email":"trader1@example.com","password":"password"}'
-- curl -X POST http://localhost:8080/api/v1/users/signup -H "Content-Type: application/json" -d '{"email":"trader2@example.com","password":"password"}'

-- use POST /api/v1/users/login to get access and refresh tokens
-- curl -X POST http://localhost:8080/api/v1/users/login -H "Content-Type: application/json" -d '{"email":"trader1@example.com","password":"password"}'
-- curl -X POST http://localhost:8080/api/v1/users/login -H "Content-Type: application/json" -d '{"email":"trader2@example.com","password":"password"}'

-- use scripts/setup.sh to setup the database with mock data and insert test users automatically
-- ./scripts/setup.sh

-- or manually insert the data by running the following commands

-- Insert sample holdings for Indian stocks
INSERT INTO holdings (user_id, symbol, quantity, average_price, current_price, total_value) VALUES 
-- User 1 holdings
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'RELIANCE', 50.00000000, 2450.75000000, 2500.00000000, 125000.00000000),
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'TCS', 25.00000000, 3650.50000000, 3700.00000000, 92500.00000000),
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'INFY', 40.00000000, 1580.25000000, 1600.00000000, 64000.00000000),
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'HDFCBANK', 30.00000000, 1450.80000000, 1480.00000000, 44400.00000000),

-- User 2 holdings
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'ICICIBANK', 60.00000000, 950.75000000, 980.00000000, 58800.00000000),
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'SBIN', 100.00000000, 420.50000000, 435.00000000, 43500.00000000),
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'WIPRO', 80.00000000, 385.25000000, 390.00000000, 31200.00000000),
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'BHARTIARTL', 70.00000000, 865.90000000, 880.00000000, 61600.00000000);

-- Insert sample positions
INSERT INTO positions (user_id, symbol, position_type, quantity, entry_price, current_price, unrealized_pnl, realized_pnl) VALUES 
-- User 1 positions
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'ADANIPORTS', 'LONG', 20.00000000, 750.00000000, 780.00000000, 600.00000000, 0.00000000),
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'BAJFINANCE', 'LONG', 15.00000000, 6800.00000000, 6750.00000000, -750.00000000, 0.00000000),

-- User 2 positions
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'MARUTI', 'LONG', 10.00000000, 10500.00000000, 10650.00000000, 1500.00000000, 0.00000000),
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'TATAMOTORS', 'SHORT', 50.00000000, 520.00000000, 510.00000000, 500.00000000, 0.00000000);

-- Insert sample order book data
INSERT INTO orderbook (symbol, side, price, quantity) VALUES 
-- RELIANCE orders
('RELIANCE', 'BUY', 2495.00000000, 10.00000000),
('RELIANCE', 'BUY', 2490.00000000, 25.00000000),
('RELIANCE', 'SELL', 2505.00000000, 15.00000000),
('RELIANCE', 'SELL', 2510.00000000, 20.00000000),

-- TCS orders
('TCS', 'BUY', 3695.00000000, 12.00000000),
('TCS', 'BUY', 3690.00000000, 18.00000000),
('TCS', 'SELL', 3705.00000000, 8.00000000),
('TCS', 'SELL', 3710.00000000, 22.00000000),

-- INFY orders
('INFY', 'BUY', 1595.00000000, 30.00000000),
('INFY', 'BUY', 1590.00000000, 45.00000000),
('INFY', 'SELL', 1605.00000000, 20.00000000),
('INFY', 'SELL', 1610.00000000, 35.00000000);
```

## Running the Application

1. **Start the server**
   ```bash
   # From the project root
   go run ./cmd/server/main.go
   ```
   
   Or build and run:
   ```bash
   go build -o bin/broker-server ./cmd/server/main.go
   ./bin/broker-server
   ```

2. **The server will start on the configured port (default: 8080)**
   ```
   INFO Starting server on port 8080
   ```

## API Endpoints

The server provides REST API endpoints for:

### Public Endpoints (No Authentication Required)

- **Health Check**
  - `GET /health` — Check if the server is running mostly used for health check when deployed (kubernetes).

- **User Authentication**
  - `POST /api/v1/users/signup` — Register a new user.
  - `POST /api/v1/users/login` — Log in and receive access/refresh tokens.

- **Token Management**
  - `POST /api/v1/auth/refresh` — Refresh access token using a valid refresh token.
  - `POST /api/v1/auth/revoke` — Revoke refresh token (logout).

---

### Authenticated Endpoints (Require Access Token)

- **Holdings**
  - `POST /api/v1/holdings` — Add a new holding for the user.
  - `GET /api/v1/holdings` — Retrieve the user's holdings.

- **Positions**
  - `GET /api/v1/positions` — Get user's current trading positions with PNL summary.

- **Order Book**
  - `GET /api/v1/orderbook` — Fetch current order book data with PNL summary.


## Testing

Use the mock data created above to test the APIs:

**Test Credentials:**
- Email: `trader1@example.com`, Password: `password`
- Email: `trader2@example.com`, Password: `password`

**Sample API Calls:**
```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"trader1@example.com","password":"password"}'

# Get holdings (replace TOKEN with actual JWT)
curl -X GET http://localhost:8080/api/holdings \
  -H "Authorization: Bearer TOKEN"
```

## Database Cleanup

To reset the database for testing:

```bash
psql -h localhost -U your_username -d broker-platform -c "DELETE FROM refresh_tokens; DELETE FROM positions; DELETE FROM holdings; DELETE FROM orderbook; DELETE FROM users;"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 