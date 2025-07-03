#!/bin/bash

# === CONFIGURATION ===
SERVER_IP="196.190.251.194"
SERVER_USER="daftech1"
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

ssh $SERVER_USER@$SERVER_IP "mkdir -p $SERVER_PATH"

rsync -avz --exclude 'node_modules' --exclude '.next' --exclude '.git' --exclude 'logs' --exclude 'dist' --exclude 'build' --exclude 'coverage' --exclude '.DS_Store' --exclude 'backups' --exclude '*.log' ./ $SERVER_USER@$SERVER_IP:$SERVER_PATH/

if [ $? -ne 0 ]; then
  print_error "File copy failed. Aborting."
  exit 1
fi

print_status "Files copied successfully."

# === 2. RUN DEPLOYMENT ON SERVER WITH SUDO AND TTY ===
print_status "Running deployment script on server with su and TTY..."

ssh -t $SERVER_USER@$SERVER_IP "cd $SERVER_PATH && chmod +x $DEPLOY_SCRIPT && su -c './$DEPLOY_SCRIPT --production'"

if [ $? -eq 0 ]; then
  print_status "Deployment completed successfully!"
  print_status "Visit: http://$SERVER_IP/"
else
  print_error "Deployment failed. Check the server logs for details."
fi 