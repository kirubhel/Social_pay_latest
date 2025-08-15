# Database Migrations
.PHONY: migrate-create migrate-up migrate-down migrate-reset migrate-status migrate-version migrate-force migrate-fix run

migrate-create:
	@read -p "Enter migration name: " name; \
	./scripts/migrate.sh create $$name

migrate-up:
	@./scripts/migrate.sh up

migrate-down:
	@./scripts/migrate.sh down

migrate-reset:
	@./scripts/migrate.sh reset

migrate-status:
	@./scripts/migrate.sh status

migrate-version:
	@./scripts/migrate.sh version

migrate-force:
	@read -p "Enter version to force: " version; \
	./scripts/migrate.sh force "$$version"

migrate-fix:
	@version=$$(./scripts/migrate.sh version | grep -oE '[0-9]+' | head -n1); \
	if [ -n "$$version" ]; then \
		echo "Current version: $$version"; \
		./scripts/migrate.sh force "$$version"; \
	else \
		echo "Error: Could not determine current version"; \
		exit 1; \
	fi

# Run the project
run:
	@echo "Installing/updating dependencies..."
	@go mod tidy
	@if [ $$? -eq 0 ]; then \
		echo "Dependencies updated successfully."; \
	else \
		echo "Error updating dependencies. Check the output above."; \
		exit 1; \
	fi
	@echo "Checking swag installation..."
	@if ! which swag > /dev/null; then \
		echo "Error: swag tool is not installed. Please install it using:"; \
		echo "go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi
	@echo "Generating Swagger documentation..."
	@swag init -g src/cmd/v2/main.go --exclude ./src/pkg/wallet/adapter/gateway/repository/generated/,./src/pkg/transaction/core/repository/generated/,./src/pkg/v2_merchant/adapter/gateway/repo/sqlc/generated/,./src/pkg/apikey_mgmt/adapter/gateway/repo/sqlc/generated/ --parseDependency --parseInternal
	@if [ $$? -eq 0 ]; then \
		echo "Swagger documentation generated successfully."; \
	else \
		echo "Error generating Swagger documentation. Check the output above."; \
		exit 1; \
	fi
	@echo "Starting the server..."
	@go run src/cmd/v2/main.go 