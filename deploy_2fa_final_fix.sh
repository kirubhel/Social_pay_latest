#!/bin/bash

# Deploy 2FA Final Fix
# This script fixes the final 2FA session issue by:
# 1. Returning existing session tokens instead of creating new ones
# 2. Avoiding the pre-session token validation issue
# 3. Simplifying the 2FA verification flow

set -e

echo "🚀 Deploying 2FA Final Fix..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    print_error "docker-compose.yml not found. Please run this script from the project root."
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

print_status "Building and deploying backend with 2FA final fix..."

# Build and restart backend
print_status "Building backend..."
docker-compose build backend

if [ $? -ne 0 ]; then
    print_error "Backend build failed!"
    exit 1
fi

print_status "Restarting backend..."
docker-compose up -d backend

if [ $? -ne 0 ]; then
    print_error "Backend restart failed!"
    exit 1
fi

# Wait for backend to be ready
print_status "Waiting for backend to be ready..."
sleep 10

# Check if backend is responding
print_status "Checking backend health..."
for i in {1..30}; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        print_success "Backend is responding!"
        break
    fi
    
    if [ $i -eq 30 ]; then
        print_error "Backend failed to start within 30 seconds!"
        docker-compose logs backend
        exit 1
    fi
    
    print_status "Waiting for backend... (attempt $i/30)"
    sleep 2
done

print_success "2FA Final Fix deployed successfully!"

echo ""
print_status "Summary of changes:"
echo "  ✅ Fixed session creation issue in 2FA verification"
echo "  ✅ Return existing session tokens instead of creating new ones"
echo "  ✅ Avoided pre-session token validation problems"
echo "  ✅ Simplified 2FA verification flow"

echo ""
print_status "Testing the fix:"
echo "  1. Try logging in with a user who has 2FA enabled"
echo "  2. Enter the 2FA code when prompted"
echo "  3. The login should complete successfully without session creation errors"

echo ""
print_warning "Note: SMS sending may still fail due to IP restrictions, but the 2FA flow should work with manual code entry"

echo ""
print_success "Deployment complete! The 2FA login should now work end-to-end." 