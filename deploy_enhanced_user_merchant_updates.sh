#!/bin/bash

echo "🚀 Deploying Enhanced User and Merchant Data Updates..."

# Navigate to backend directory
cd Backend

echo "🔨 Building backend..."
go build -o bin/v1/main src/cmd/v1/main.go

echo "🐳 Building backend Docker image..."
docker build -t socialpay-backend:latest .

echo "🔄 Restarting backend container..."
cd ..
docker-compose up -d --build backend

# Navigate to frontend directory
cd frontend

echo "📦 Installing frontend dependencies..."
npm install

echo "🔨 Building frontend..."
npm run build

echo "🐳 Building frontend Docker image..."
docker build -t socialpay-frontend:latest .

echo "🔄 Restarting frontend container..."
cd ..
docker-compose up -d --build frontend

echo "✅ Enhanced User and Merchant Data Updates deployment completed!"
echo ""
echo "🔗 What's enhanced:"
echo "   • Comprehensive user profile updates (first_name, last_name, sir_name, phone_number)"
echo "   • New /user-profile endpoint to fetch user data"
echo "   • Enhanced /user-update endpoint with authentication"
echo "   • Improved database queries with proper joins"
echo "   • Better error handling and validation"
echo "   • Real-time data synchronization between user and business info"
echo "   • Loading states and user feedback"
echo ""
echo "🌐 Frontend should be available at: http://localhost:3000"
echo "🔧 Backend should be available at: http://localhost:8080"
echo "📝 Test the settings page to verify enhanced functionality" 