#!/bin/bash

echo "🚀 Deploying Settings Integration with Backend..."

# Navigate to frontend directory
cd frontend

echo "📦 Installing dependencies..."
npm install

echo "🔨 Building frontend..."
npm run build

echo "🐳 Building Docker image..."
docker build -t socialpay-frontend:latest .

echo "🔄 Restarting frontend container..."
cd ..
docker-compose up -d --build frontend

echo "✅ Settings Integration deployment completed!"
echo ""
echo "🔗 What's connected:"
echo "   • User profile update via /user-update endpoint"
echo "   • Business information via /api/get/merchant/details"
echo "   • Merchant details update via /api/merchant/update"
echo "   • Real-time data loading and saving"
echo "   • Loading states and error handling"
echo "   • Toast notifications for user feedback"
echo ""
echo "🌐 Frontend should be available at: http://localhost:3000"
echo "📝 Test the settings page to verify backend integration" 