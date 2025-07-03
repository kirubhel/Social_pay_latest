#!/bin/bash

# 2FA Database Migration Script
# This script runs the migration to add 2FA support to the database

echo "Starting 2FA database migration..."

# Set database connection details (update these as needed)
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_NAME=${DB_NAME:-"socialpay"}
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD=${DB_PASSWORD:-"password"}

# Migration file path
MIGRATION_FILE="Backend/db/migrations/000030_add_2fa_support.up.sql"

# Check if migration file exists
if [ ! -f "$MIGRATION_FILE" ]; then
    echo "Error: Migration file not found: $MIGRATION_FILE"
    exit 1
fi

echo "Running migration: $MIGRATION_FILE"

# Run the migration using psql
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$MIGRATION_FILE"

if [ $? -eq 0 ]; then
    echo "✅ 2FA migration completed successfully!"
    echo "The following changes have been applied:"
    echo "  - Added two_factor_enabled column to auth.users table"
    echo "  - Added two_factor_verified_at column to auth.users table"
    echo "  - Created auth.two_factor_codes table for storing verification codes"
    echo "  - Added indexes for better performance"
    echo "  - Added trigger for automatic updated_at timestamp updates"
else
    echo "❌ Migration failed!"
    exit 1
fi

echo ""
echo "Next steps:"
echo "1. Restart your backend application"
echo "2. Test the 2FA functionality in the frontend"
echo "3. Verify that SMS codes are being sent correctly" 