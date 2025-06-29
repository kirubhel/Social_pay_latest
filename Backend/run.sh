#!/bin/bash

# Generate Swagger documentation
echo "Generating Swagger documentation..."
~/go/bin/swag init -g src/cmd/v2/main.go --exclude ./src/pkg/wallet/adapter/gateway/repository/generated/,./src/pkg/transaction/core/repository/generated/,./src/pkg/v2_merchant/adapter/gateway/repo/sqlc/generated/,./src/pkg/apikey_mgmt/adapter/gateway/repo/sqlc/generated/ --parseDependency --parseInternal

# Check if Swagger generation was successful
if [ $? -eq 0 ]; then
    echo "Swagger documentation generated successfully."
else
    echo "Error generating Swagger documentation. Check the output above."
    exit 1
fi

# Run the server
echo "Starting the server..."
go run src/cmd/v2/main.go 