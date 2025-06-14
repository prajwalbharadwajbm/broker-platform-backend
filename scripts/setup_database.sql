-- Broker Platform Database Setup Script
-- Run this script after creating the database and connecting to it

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create refresh tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create holdings table
CREATE TABLE IF NOT EXISTS holdings (
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
CREATE TABLE IF NOT EXISTS positions (
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
CREATE TABLE IF NOT EXISTS orderbook (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(4) NOT NULL CHECK (side IN ('BUY', 'SELL')),
    price NUMERIC(20,8) NOT NULL,
    quantity NUMERIC(20,8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


-- Insert test users
-- use POST /api/v1/users/signup to create users
-- curl -X POST http://localhost:8080/api/v1/users/signup -H "Content-Type: application/json" -d '{"email":"trader1@example.com","password":"password"}'
-- curl -X POST http://localhost:8080/api/v1/users/signup -H "Content-Type: application/json" -d '{"email":"trader2@example.com","password":"password"}'

-- use POST /api/v1/users/login to get access and refresh tokens
-- curl -X POST http://localhost:8080/api/v1/users/login -H "Content-Type: application/json" -d '{"email":"trader1@example.com","password":"password"}'
-- curl -X POST http://localhost:8080/api/v1/users/login -H "Content-Type: application/json" -d '{"email":"trader2@example.com","password":"password"}'

-- Insert sample holdings for Indian stocks
INSERT INTO holdings (user_id, symbol, quantity, average_price, current_price, total_value) VALUES 
-- User 1 holdings (Major Indian Blue Chips)
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'RELIANCE', 50.00000000, 2450.75000000, 2500.00000000, 125000.00000000),
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'TCS', 25.00000000, 3650.50000000, 3700.00000000, 92500.00000000),
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'INFY', 40.00000000, 1580.25000000, 1600.00000000, 64000.00000000),
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'HDFCBANK', 30.00000000, 1450.80000000, 1480.00000000, 44400.00000000),

-- User 2 holdings (Banking & Telecom Stocks)
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'ICICIBANK', 60.00000000, 950.75000000, 980.00000000, 58800.00000000),
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'SBIN', 100.00000000, 420.50000000, 435.00000000, 43500.00000000),
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'WIPRO', 80.00000000, 385.25000000, 390.00000000, 31200.00000000),
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'BHARTIARTL', 70.00000000, 865.90000000, 880.00000000, 61600.00000000)
ON CONFLICT DO NOTHING;

-- Insert sample positions
INSERT INTO positions (user_id, symbol, position_type, quantity, entry_price, current_price, unrealized_pnl, realized_pnl) VALUES 
-- User 1 positions
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'ADANIPORTS', 'LONG', 20.00000000, 750.00000000, 780.00000000, 600.00000000, 0.00000000),
((SELECT id FROM users WHERE email = 'trader1@example.com'), 'BAJFINANCE', 'LONG', 15.00000000, 6800.00000000, 6750.00000000, -750.00000000, 0.00000000),

-- User 2 positions
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'MARUTI', 'LONG', 10.00000000, 10500.00000000, 10650.00000000, 1500.00000000, 0.00000000),
((SELECT id FROM users WHERE email = 'trader2@example.com'), 'TATAMOTORS', 'SHORT', 50.00000000, 520.00000000, 510.00000000, 500.00000000, 0.00000000)
ON CONFLICT DO NOTHING;

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
('INFY', 'SELL', 1610.00000000, 35.00000000),

-- HDFCBANK orders
('HDFCBANK', 'BUY', 1475.00000000, 25.00000000),
('HDFCBANK', 'SELL', 1485.00000000, 30.00000000),

-- ICICIBANK orders
('ICICIBANK', 'BUY', 975.00000000, 40.00000000),
('ICICIBANK', 'SELL', 985.00000000, 35.00000000)
ON CONFLICT DO NOTHING;

-- Display summary of inserted data
SELECT 'Users created:' as info, COUNT(*) as count FROM users;
SELECT 'Holdings created:' as info, COUNT(*) as count FROM holdings;
SELECT 'Positions created:' as info, COUNT(*) as count FROM positions;
SELECT 'Order book entries:' as info, COUNT(*) as count FROM orderbook; 