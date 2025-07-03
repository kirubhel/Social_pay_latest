#!/bin/bash

echo "🔧 Fixing phone data issues..."

# Run the debug script first to see the current state
echo "📊 Checking current phone data state..."
docker-compose exec -T postgres psql -U postgres -d socialpay -f /scripts/debug_phone_data.sql

echo ""
echo "🔧 Applying phone data fixes..."
docker-compose exec -T postgres psql -U postgres -d socialpay -f /scripts/fix_phone_data.sql

echo ""
echo "✅ Phone data fixes applied. Restarting backend..."
docker-compose restart backend

echo ""
echo "🔄 Backend restarted. Testing 2FA functionality..."
echo "You can now try enabling 2FA again for user: 9b29fcdb-3326-4604-8a31-8ee844ec3fef"
echo ""
echo "To test, use the frontend or make a POST request to:"
echo "POST /auth/2fa/enable"
echo "Authorization: Bearer <your-session-token>" 