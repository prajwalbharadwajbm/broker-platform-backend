# Refresh Token API Documentation

This document describes the refresh token mechanism implemented in the broker platform backend.

## Overview

The refresh token mechanism is critical for the broker platform as it enables secure, uninterrupted trading sessions for users who need to maintain active connections to financial markets. In a trading environment, users often keep their applications open for extended periods to monitor positions, execute trades, and respond to market movements.

My broker platform uses a dual-token system:

- **Access Token**: Short-lived (10 minutes) JWT token containing user permissions and account information for API authentication
- **Refresh Token**: Long-lived (7 days) cryptographically secure random token stored in the database to obtain new access tokens

## Why Refresh Tokens are Essential for Broker Platforms

### 1. **Continuous Market Access**
- Traders need uninterrupted access to live market data and trading capabilities
- Short access token lifespans (10 minutes) provide security without disrupting active trading sessions
- Automatic token refresh ensures users don't lose connection during critical market moments

### 2. **Multi-Device Trading Support**
- Traders often use multiple devices (desktop, mobile, tablet) simultaneously
- Each device maintains its own refresh token, allowing independent session management
- Users can trade from their phone while monitoring positions on their desktop

### 3. **Enhanced Security for Financial Data**
- Financial platforms are high-value targets for attackers
- Token rotation on every refresh minimizes the window of vulnerability if a token is compromised
- Database-stored refresh tokens can be instantly revoked if suspicious activity is detected

## How It Works in My Broker Platform

### Authentication Flow
1. User logs in through the login endpoint with broker credentials
2. Backend validates user against the user database and broker account permissions
3. System generates both tokens and stores refresh token linked to user's broker account
4. Client receives tokens and can immediately begin trading operations

### Session Management During Trading
1. Client makes API calls to trading endpoints (orders, positions, market data) using access token
2. When access token expires (10 minutes), client automatically refreshes using refresh token
3. New token pair is issued, old refresh token is invalidated (token rotation)
4. Trading continues seamlessly without user intervention

### Multi-Session Support
- Each login creates a separate refresh token entry in the database
- Users can have active sessions on multiple devices simultaneously
- Each session is tracked independently for security and audit purposes

## API Endpoints

### 1. Login (Modified for Broker Platform)
**POST** `/api/v1/users/login`

Authenticates broker platform users and returns tokens for accessing trading APIs.

**Request Body:**
```json
{
  "email": "dummy@dummy.com",
  "password": "dummy"
}
```

**Response:**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "a1b2c3d4e5f6789012345678901234567890abcdef",
    "token_type": "Bearer",
    "expires_in": 900,
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

The access token contains user permissions for trading operations, account access levels, and broker account identifiers.

### 2. Refresh Token (Trading Session Continuity)
**POST** `/api/v1/auth/refresh`

Maintains active trading sessions by exchanging refresh tokens for new access tokens. This endpoint is called automatically by trading clients to ensure uninterrupted market access.

**Request Body:**
```json
{
  "refresh_token": "a1b2c3d4e5f6789012345678901234567890abcdef"
}
```

**Response:**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "b2c3d4e5f6789012345678901234567890abcdef1a",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

### 3. Revoke Refresh Token (Secure Logout)
**POST** `/api/v1/auth/revoke`

Securely terminates trading sessions and revokes access to all trading APIs. Critical for ending sessions when users finish trading or when security incidents occur.

**Request Body:**
```json
{
  "refresh_token": "a1b2c3d4e5f6789012345678901234567890abcdef"
}
```

**Response:**
```json
{
  "data": {
    "message": "Token revoked successfully"
  }
}
```

## Client Implementation for Trading Applications

### Initial Authentication Flow
1. Trading client calls `/api/v1/users/login` with broker credentials
2. Store refresh tokens securely 
3. Generate Access Token using refresh token when needed and use it for all trading API requests: orders, positions, market data, account info

### Seamless Token Refresh During Trading
1. When any trading API returns 401 (token expired), immediately call `/api/v1/auth/refresh`
2. Update refresh tokens with new values
3. Retry the original trading operation with the new access token generated when refreshed.
4. Implement this flow (client side) to be invisible to the user to maintain uninterrupted trading.

### Secure Session Termination
1. Call `/api/v1/auth/revoke` when user explicitly logs out or closes trading application
2. Clear all refresh tokens from client memory and secure storage
3. Consider implementing automatic logout after extended inactivity for security


## Configuration for Broker Platform

The system uses environment variables specific to our broker platform:

- `JWT_SECRET`: Used for signing access tokens containing trading permissions
- `REFRESH_TOKEN_EXPIRY_DAYS`: Set to 7 days for balance between security and user experience
- `ACCESS_TOKEN_EXPIRY_MINUTES`: Set to 10 minutes for frequent rotation in trading environment