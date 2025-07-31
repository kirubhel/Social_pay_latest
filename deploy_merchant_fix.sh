#!/bin/bash

echo "🚀 Deploying Merchant Not Found Fix..."

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

echo "✅ Merchant Not Found Fix deployment completed!"
echo ""
echo "🔧 What was fixed:"
echo "   • Fixed syntax error in merchant.go"
echo "   • Updated merchant endpoint to handle missing merchant profiles gracefully"
echo "   • Added createMerchant API function for new users"
echo "   • Updated frontend to handle both existing and new merchant profiles"
echo "   • Fixed field name mappings between frontend and backend"
echo ""
echo "🎯 Now the system will:"
echo "   • Return success with null data instead of error when no merchant exists"
echo "   • Allow users to create merchant profiles from settings"
echo "   • Handle both update and create scenarios seamlessly"
echo ""
echo "🚀 Ready to test!" 