#!/bin/bash

# === CONFIGURATION ===
SERVER_IP="196.190.251.194"
SERVER_USER="daftech1"
SERVER_PASSWORD="DAFTech@2025new"
SU_PASSWORD="DAFTech@2024"
SERVER_PATH="~/socialpay"
LOCAL_PATH="$(pwd)"
DEPLOY_SCRIPT="deploy.sh"

# === FUNCTIONS ===
print_status() {
  echo -e "\033[1;32m[INFO]\033[0m $1"
}
print_error() {
  echo -e "\033[1;31m[ERROR]\033[0m $1"
}

# === 1. COPY FILES TO SERVER ===
print_status "Copying project files to $SERVER_USER@$SERVER_IP:$SERVER_PATH ..."

sshpass -p "$SERVER_PASSWORD" ssh $SERVER_USER@$SERVER_IP "mkdir -p $SERVER_PATH"

# Copy all files including V2 backend
sshpass -p "$SERVER_PASSWORD" rsync -avz --exclude 'node_modules' --exclude '.next' --exclude '.git' --exclude 'logs' --exclude 'dist' --exclude 'build' --exclude 'coverage' --exclude '.DS_Store' --exclude 'backups' --exclude '*.log' ./ $SERVER_USER@$SERVER_IP:$SERVER_PATH/

if [ $? -ne 0 ]; then
  print_error "File copy failed. Aborting."
  exit 1
fi

print_status "Files copied successfully."

# === 2. RUN DEPLOYMENT ON SERVER ===
print_status "Running deployment script on server..."

# Create a temporary script with sudo commands
sshpass -p "$SERVER_PASSWORD" ssh $SERVER_USER@$SERVER_IP "cd $SERVER_PATH && chmod +x $DEPLOY_SCRIPT && echo '$SU_PASSWORD' | sudo -S bash -c 'cd $SERVER_PATH && ./$DEPLOY_SCRIPT --production'"

if [ $? -eq 0 ]; then
  print_status "Deployment completed successfully!"
  print_status "Visit: http://$SERVER_IP/"
  print_status "V2 API: http://$SERVER_IP:8082"
  print_status "Swagger: http://$SERVER_IP:8082/swagger/index.html"
  print_status ""
  print_status "All services deployed:"
  print_status "  ✓ Frontend (Next.js)"
  print_status "  ✓ Backend V1 (Mux)"
  print_status "  ✓ Backend V2 (Gin) with Swagger"
  print_status "  ✓ Nginx (Reverse Proxy)"
else
  print_error "Deployment failed. Check the server logs for details."
fi 