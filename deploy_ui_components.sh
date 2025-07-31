#!/bin/bash

echo "🚀 Deploying UI Components and Settings Improvements..."

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

echo "✅ UI Components deployment completed!"
echo ""
echo "🎨 What's new:"
echo "   • Modern InputField component with animations"
echo "   • Professional TextareaField component"
echo "   • Drag & drop FileUpload component"
echo "   • Updated General Settings with new components"
echo "   • Improved Password form with consistent styling"
echo "   • Better focus states and transitions"
echo ""
echo "🌐 Frontend should be available at: http://localhost:3000" 