# Database Design - Broker Platform

## Overview

The Broker Platform uses PostgreSQL as its primary database with a focus on financial data integrity, security, and performance. The design follows normalized relational principles while optimizing for trading operations and real-time data access.

## Database Schema

### Architecture Principles

- **UUID Primary Keys**: All tables use UUID for primary keys to ensure global uniqueness and security
- **Timestamp Tracking**: Created/updated timestamps for audit trails. All timestamps are in UTC timezone.
- **Referential Integrity**: Foreign key constraints ensure data consistency
- **Precision Decimals**: Financial amounts use `NUMERIC(20,8)` for precision
- **Data Constraints**: Check constraints for business rule enforcement

## Tables

### 1. Users Table

**Purpose**: Core user authentication and account management

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

**Design Decisions**:
- **UUID Primary Key**: Prevents user enumeration attacks
- **Email Uniqueness**: Enforced at database level for data integrity
- **Password Hashing**: Stores bcrypt hashes, never plaintext passwords
- **Timestamps**: Track account creation and modification

**Indexes**:
- Primary key on `id`
- Unique index on `email` for fast lookups

### 2. Refresh Tokens Table

**Purpose**: JWT refresh token management for secure authentication

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

**Design Decisions**:
- **Cascade Delete**: When user is deleted, all tokens are automatically removed
- **Token Uniqueness**: Prevents token duplication
- **Expiration Tracking**: Explicit expiry time for security
- **Cleanup Ready**: Structure supports automated cleanup of expired tokens

**Relationships**:
- `user_id` → `users.id` (Many-to-One)

### 3. Holdings Table

**Purpose**: Track user's stock portfolio and current positions

```sql
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
```

**Design Decisions**:
- **High Precision Decimals**: `NUMERIC(20,8)` handles fractional shares and precise pricing
- **Symbol Storage**: VARCHAR(20) accommodates various stock symbol formats
- **Calculated Fields**: `total_value` stored for performance (denormalized for speed)
- **Price Tracking**: Both average purchase price and current market price

**Financial Calculations**:
- Total Value = Quantity × Current Price
- Profit/Loss = (Current Price - Average Price) × Quantity

**Relationships**:
- `user_id` → `users.id` (Many-to-One)

### 4. Positions Table

**Purpose**: Active trading positions with P&L tracking

```sql
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
```

**Design Decisions**:
- **Position Types**: Check constraint ensures only 'LONG' or 'SHORT' positions
- **P&L Tracking**: Separate unrealized and realized profit/loss
- **Entry Price**: Original trade execution price
- **Real-time Updates**: Current price and unrealized P&L updated frequently

**P&L Calculations**:
- **LONG Position**: Unrealized P&L = (Current Price - Entry Price) × Quantity
- **SHORT Position**: Unrealized P&L = (Entry Price - Current Price) × Quantity

**Relationships**:
- `user_id` → `users.id` (Many-to-One)

### 5. Orderbook Table

**Purpose**: Market depth and order matching data

```sql
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

**Design Decisions**:
- **Order Side**: Check constraint for 'BUY' or 'SELL' only
- **Market Data**: Not tied to specific users (market-wide data)
- **Price Levels**: Each row represents a price level in the order book
- **High Frequency**: Designed for frequent updates during trading hours

**Order Book Structure**:
- **Bid Side**: BUY orders (price descending)
- **Ask Side**: SELL orders (price ascending)
- **Spread**: Difference between highest bid and lowest ask

## Indexes and Performance

### Recommended Indexes to be created for better performance as its high frequency data

```sql
-- User lookup optimization
CREATE INDEX idx_users_email ON users(email);

-- Token cleanup and validation
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- Holdings queries by user
CREATE INDEX idx_holdings_user_id ON holdings(user_id);
CREATE INDEX idx_holdings_symbol ON holdings(symbol);

-- Position queries and P&L calculations
CREATE INDEX idx_positions_user_id ON positions(user_id);
CREATE INDEX idx_positions_symbol ON positions(symbol);

-- Order book queries (market data)
CREATE INDEX idx_orderbook_symbol_side ON orderbook(symbol, side);
CREATE INDEX idx_orderbook_symbol_price ON orderbook(symbol, price);
```

## Data Types Rationale

### NUMERIC(20,8) for Financial Data
- **Precision**: 8 decimal places for fractional shares
- **Range**: 20 total digits handle large monetary values
- **Accuracy**: Avoids floating-point rounding errors

### UUID for Primary Keys
- **Security**: Prevents ID enumeration attacks
- **Performance**: Good for insert-heavy workloads

### TIMESTAMPTZ for Timestamps
- **Timezone Aware**: Handles global users correctly by using UTC timezone.
- **Precision**: Microsecond accuracy for trading timestamps

## Business Rules Enforced

### Database Constraints
1. **Email Uniqueness**: `UNIQUE(email)` in users table
2. **Position Types**: `CHECK (position_type IN ('LONG', 'SHORT'))`
3. **Order Sides**: `CHECK (side IN ('BUY', 'SELL'))`
4. **Referential Integrity**: All foreign keys with CASCADE DELETE
