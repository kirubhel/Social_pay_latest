#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
fi

# Set the migrate command path
MIGRATE_CMD=~/go/bin/migrate

# Construct database URL from environment variables
DB_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${SSL_MODE}"

# Use the constructed URL or fallback to default if env vars are not set
MIGRATION_DB_URL=${DATABASE_URL:-$DB_URL}

# Function to display usage
usage() {
    echo "Usage: $0 <command> [args]"
    echo "Commands:"
    echo "  create <name>  - Create a new migration"
    echo "  up            - Run all pending migrations"
    echo "  down          - Rollback the last migration"
    echo "  reset         - Rollback all migrations and run them again"
    echo "  version       - Show current migration version"
    echo "  status        - Show migration status"
    echo "  force <version> - Force set database version"
    echo "  fix           - Fix dirty database state"
}

# Check if migrate tool is installed
if ! command -v $MIGRATE_CMD &> /dev/null; then
    echo "Error: migrate tool is not installed"
    echo "Install it with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Function to check if database is dirty
is_dirty() {
    local status
    status=$($MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations version 2>&1)
    if echo "$status" | grep -q "dirty"; then
        return 0  # true in bash
    else
        return 1  # false in bash
    fi
}

# Function to get current version
get_current_version() {
    local version
    version=$($MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations version 2>/dev/null)
    if [ $? -eq 0 ]; then
        echo "$version" | grep -oE '^[0-9]+'
    else
        echo ""
    fi
}

# Parse command line arguments
cmd=$1
case $cmd in
    "create")
        name=$2
        if [ -z "$name" ]; then
            echo "Error: Migration name is required"
            usage
            exit 1
        fi
        
        # Create migrations directory if it doesn't exist
        mkdir -p db/migrations
        
        
        $MIGRATE_CMD create -ext sql -dir db/migrations -seq $name
        
        # Verify creation
        if [ $? -eq 0 ]; then
            echo "Successfully created migration files:"
        else
            echo "Failed to create migration files"
            exit 1
        fi
        ;;
    "up")
        $MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations up
        ;;
    "down")
        $MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations down 1
        ;;
    "reset")
        $MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations down -all
        $MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations up
        ;;
    "version")
        $MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations version
        ;;
    "status")
        $MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations status
        ;;
    "force")
        version=$2
        if [ -z "$version" ]; then
            echo "Error: Version number is required"
            usage
            exit 1
        fi
        echo "Forcing migration version to: $version"
        $MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations force $version
        ;;
    "fix")
        if is_dirty; then
            current_version=$(get_current_version)
            if [ -n "$current_version" ]; then
                echo "Database is in dirty state at version: $current_version"
                echo "Attempting to fix by forcing current version..."
                $MIGRATE_CMD -database "${MIGRATION_DB_URL}" -path db/migrations force $current_version
                if [ $? -eq 0 ]; then
                    echo "Successfully fixed dirty state"
                    echo "You may now need to manually fix any partially applied migrations"
                    echo "Current schema version is now set to: $current_version"
                else
                    echo "Failed to fix dirty state"
                    exit 1
                fi
            else
                echo "Error: Could not determine current version"
                exit 1
            fi
        else
            echo "Database is not in a dirty state"
        fi
        ;;
    *)
        usage
        exit 1
        ;;
esac 