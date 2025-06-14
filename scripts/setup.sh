#!/bin/bash

# Broker Platform Database Setup Script
# This script creates the database and sets up all tables with mock data

set -e

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${POSTGRES_USER:-postgres}
DB_NAME="broker-platform"

echo "Setting up Broker Platform Database..."

# Check if PostgreSQL is running
if ! command -v psql &> /dev/null; then
    echo "PostgreSQL is not installed or not in PATH"
    exit 1
fi

# Create database
echo "Creating database '$DB_NAME'..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE \"$DB_NAME\";" 2>/dev/null || echo "Database already exists"

# Run setup script
echo "Setting up tables and inserting mock data..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/setup_database.sql

echo "Database setup completed successfully!"
echo ""
echo "Test credentials:"
echo "   Email: trader1@example.com | Password: password"
echo "   Email: trader2@example.com | Password: password"
echo ""
echo "You can now start the server with:"
echo "   go run ./cmd/server/main.go" 