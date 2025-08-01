#!/bin/bash

# Simple deployment script with sudo
SERVER_IP="196.190.251.194"
SERVER_USER="daftech1"
SERVER_PASSWORD="DAFTech@2025new"
SU_PASSWORD="DAFTech@2024"
SERVER_PATH="~/socialpay"

echo "Deploying to server..."

# Copy files
sshpass -p "$SERVER_PASSWORD" rsync -avz --exclude 'node_modules' --exclude '.next' --exclude '.git' --exclude 'logs' --exclude 'dist' --exclude 'build' --exclude 'coverage' --exclude '.DS_Store' --exclude 'backups' --exclude '*.log' ./ $SERVER_USER@$SERVER_IP:$SERVER_PATH/

# Run deployment with su - ensuring latest code is deployed
echo "Stopping all existing socialpay containers..."
sshpass -p "$SERVER_PASSWORD" ssh $SERVER_USER@$SERVER_IP "cd $SERVER_PATH && echo '$SU_PASSWORD' | su -c 'docker-compose down'"

echo "Stopping old containers that might be running..."
sshpass -p "$SERVER_PASSWORD" ssh $SERVER_USER@$SERVER_IP "echo '$SU_PASSWORD' | su -c 'docker stop socialpay-backend socialpay-backend-v1 socialpay-backend-v2 socialpay-frontend 2>/dev/null || true'"
sshpass -p "$SERVER_PASSWORD" ssh $SERVER_USER@$SERVER_IP "echo '$SU_PASSWORD' | su -c 'docker rm socialpay-backend socialpay-backend-v1 socialpay-backend-v2 socialpay-frontend 2>/dev/null || true'"

echo "Building latest images..."
sshpass -p "$SERVER_PASSWORD" ssh $SERVER_USER@$SERVER_IP "cd $SERVER_PATH && echo '$SU_PASSWORD' | su -c 'docker-compose build --no-cache'"

echo "Starting services with latest code..."
sshpass -p "$SERVER_PASSWORD" ssh $SERVER_USER@$SERVER_IP "cd $SERVER_PATH && echo '$SU_PASSWORD' | su -c 'docker-compose up -d'"

echo "Deployment completed!"
echo "Check services at:"
echo "  Frontend: http://$SERVER_IP:3000"
echo "  Backend V1: http://$SERVER_IP:8004"
echo "  Backend V2: http://$SERVER_IP:8082"
echo "  Swagger: http://$SERVER_IP:8082/swagger/index.html" 