# Build stage
FROM golang:1.23-alpine AS builder

# Install ca-certificates for HTTPS requests and git for Go modules
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/cmd/v1/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh appuser

WORKDIR /home/appuser/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy the environment file (.env only)
COPY --from=builder /app/.env /home/appuser/.env

# Set proper permissions
RUN chown appuser:appuser /home/appuser/main /home/appuser/.env
RUN chmod +x /home/appuser/main

# Expose port 8004
EXPOSE 8004

# Switch to non-root user
USER appuser

# Command to run
CMD ["./main"]
