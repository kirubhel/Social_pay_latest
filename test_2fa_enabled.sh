#!/bin/bash

echo "🧪 Testing 2FA enabled user login..."

# Set error handling
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

print_status "Docker is running"

# Stop existing containers
print_status "Stopping existing containers..."
docker-compose down

# Build and start backend
print_status "Building and starting backend..."
cd Backend
docker-compose up --build -d

# Wait for backend to be ready
print_status "Waiting for backend to be ready..."
sleep 15

# Test the login endpoint
print_status "Testing login endpoint..."

echo ""
echo "Testing login with 2FA enabled user..."
echo "Expected: OTP_REQUIRED (because 2FA is enabled)"
echo ""

curl -X POST http://localhost:8004/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "prefix": "251",
    "number": "911123456",
    "password": "testpassword"
  }' | jq '.'

echo ""
print_status "✅ Test completed!"
print_status "Expected behavior:"
echo "  • If user has 2FA enabled → OTP_REQUIRED"
echo "  • If user has 2FA disabled → Direct login (token + user data)"

echo ""
print_status "📝 Logs:"
echo "  Backend logs: docker-compose logs -f backend" 