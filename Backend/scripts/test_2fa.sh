#!/bin/bash

# 2FA Test Script
# This script tests the 2FA functionality

echo "🧪 Testing 2FA Functionality"
echo "=============================="

# Set your API base URL
API_BASE_URL=${API_BASE_URL:-"http://localhost:8004"}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test function
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${YELLOW}Testing: $description${NC}"
    echo "Endpoint: $method $API_BASE_URL$endpoint"
    
    if [ -n "$data" ]; then
        echo "Data: $data"
        response=$(curl -s -X $method \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer YOUR_TOKEN_HERE" \
            -d "$data" \
            "$API_BASE_URL$endpoint")
    else
        response=$(curl -s -X $method \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer YOUR_TOKEN_HERE" \
            "$API_BASE_URL$endpoint")
    fi
    
    echo "Response: $response"
    
    # Check if response contains success
    if echo "$response" | grep -q '"success":true'; then
        echo -e "${GREEN}✅ PASS${NC}"
    else
        echo -e "${RED}❌ FAIL${NC}"
    fi
}

echo -e "\n${YELLOW}Note: Replace 'YOUR_TOKEN_HERE' with a valid authentication token${NC}"

# Test 2FA Status
test_endpoint "GET" "/auth/2fa/status" "" "Get 2FA Status"

# Test Enable 2FA
test_endpoint "POST" "/auth/2fa/enable" "" "Enable 2FA"

# Test Verify 2FA Setup (with a sample code)
test_endpoint "POST" "/auth/2fa/verify-setup" '{"code":"123456"}' "Verify 2FA Setup"

# Test Resend 2FA Code
test_endpoint "POST" "/auth/2fa/resend" "" "Resend 2FA Code"

# Test Disable 2FA (with password)
test_endpoint "POST" "/auth/2fa/disable" '{"password":"your_password"}' "Disable 2FA"

echo -e "\n${GREEN}🎉 2FA Testing Complete!${NC}"
echo -e "\n${YELLOW}Next Steps:${NC}"
echo "1. Run the database migration: ./Backend/scripts/migrate_2fa.sh"
echo "2. Restart your backend application"
echo "3. Test with real authentication tokens"
echo "4. Verify SMS delivery" 