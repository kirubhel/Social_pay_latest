#!/bin/bash

echo "🎨 Deploying Modern UI Components..."

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

echo "✅ Modern UI Components deployment completed!"
echo ""
echo "✨ What's new in the modern UI:"
echo "   • Enhanced InputField with gradient effects and animations"
echo "   • Modern TextareaField with character counter and progress bar"
echo "   • Improved FileUpload with drag & drop and visual feedback"
echo "   • Better focus states with ring effects and shadows"
echo "   • Smooth hover animations and transitions"
echo "   • Professional backdrop blur effects"
echo "   • Password strength indicators"
echo "   • Enhanced error and success states"
echo "   • Modern rounded corners (xl) and spacing"
echo ""
echo "🌐 Frontend should be available at: http://localhost:3000"
echo "🎯 Check out the General Settings tab to see the new components!" 