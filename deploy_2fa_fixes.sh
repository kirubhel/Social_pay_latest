#!/bin/bash

echo "🚀 Deploying 2FA fixes for registration and login flow..."

# Set error handling
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
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
sleep 10

# Check if backend is running
if ! curl -f http://localhost:8004/health > /dev/null 2>&1; then
    print_warning "Backend health check failed, but continuing..."
fi

# Build and start frontend
print_status "Building and starting frontend..."
cd ../frontend
docker-compose up --build -d

# Wait for frontend to be ready
print_status "Waiting for frontend to be ready..."
sleep 15

# Check if frontend is running
if ! curl -f http://localhost:3000 > /dev/null 2>&1; then
    print_warning "Frontend health check failed, but continuing..."
fi

print_status "✅ Deployment completed!"
print_status "Frontend: http://localhost:3000"
print_status "Backend: http://localhost:8004"

echo ""
print_status "🔧 Changes deployed:"
echo "  • Registration now requires phone prefix (default: +251)"
echo "  • Login flow now checks 2FA status after successful authentication"
echo "  • 2FA verification is required if enabled for user"
echo "  • New backend endpoint: /auth/2fa/send-login"

echo ""
print_status "🧪 Testing instructions:"
echo "  1. Register a new user with phone prefix"
echo "  2. Enable 2FA in settings"
echo "  3. Log out and log back in"
echo "  4. Verify 2FA code is sent and required"

echo ""
print_status "📝 Logs:"
echo "  Frontend logs: docker-compose logs -f frontend"
echo "  Backend logs: docker-compose logs -f backend" 