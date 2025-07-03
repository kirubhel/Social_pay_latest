#!/bin/bash

# SocialPay Deployment Script
# This script automates the deployment of the SocialPay application using Docker

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="socialpay"
BACKUP_DIR="./backups"
LOG_FILE="./deployment.log"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Function to log messages
log_message() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> "$LOG_FILE"
}

# Function to check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    print_status "Docker and Docker Compose are installed"
}

# Function to create backup
create_backup() {
    print_header "Creating Backup"
    
    mkdir -p "$BACKUP_DIR"
    
    BACKUP_FILE="$BACKUP_DIR/${APP_NAME}-backup-$(date +%Y%m%d-%H%M%S).tar.gz"
    
    if [ -f "docker-compose.yml" ] || [ -f "Backend/_env" ]; then
        tar -czf "$BACKUP_FILE" \
            docker-compose.yml \
            nginx.conf \
            Backend/_env \
            2>/dev/null || true
        
        print_status "Backup created: $BACKUP_FILE"
        log_message "Backup created: $BACKUP_FILE"
    else
        print_warning "No existing configuration found to backup"
    fi
}

# Function to stop existing services
stop_services() {
    print_header "Stopping Existing Services"
    
    if [ -f "docker-compose.yml" ]; then
        docker-compose down
        print_status "Existing services stopped"
        log_message "Existing services stopped"
    else
        print_warning "No docker-compose.yml found, skipping service stop"
    fi
}

# Function to clean up Docker resources
cleanup_docker() {
    print_header "Cleaning Up Docker Resources"
    
    # Remove unused images
    docker image prune -f
    
    # Remove unused volumes (be careful with this)
    # docker volume prune -f
    
    print_status "Docker cleanup completed"
    log_message "Docker cleanup completed"
}

# Function to build and start services
start_services() {
    print_header "Building and Starting Services"
    
    # Build and start services
    docker-compose build
    docker-compose up --build -d
    
    # Wait for services to be healthy
    print_status "Waiting for services to start..."
    sleep 30
    
    # Check service status
    docker-compose ps
    
    print_status "Services started successfully"
    log_message "Services started successfully"
}

# Function to run health checks
health_check() {
    print_header "Running Health Checks"
    
    # Check if services are running
    if docker-compose ps | grep -q "Up"; then
        print_status "Services are running"
    else
        print_error "Some services are not running"
        docker-compose logs --tail=50
        exit 1
    fi
    
    # Test endpoints
    print_status "Testing application endpoints..."
    
    # Wait a bit more for applications to fully start
    sleep 10
    
    # Test frontend (if accessible)
    if curl -f -s http://localhost:3000 > /dev/null; then
        print_status "Frontend is responding"
    else
        print_warning "Frontend is not responding on port 3000"
    fi
    
    # Test backend (if accessible)
    if curl -f -s http://localhost:8008 > /dev/null; then
        print_status "Backend is responding"
    else
        print_warning "Backend is not responding on port 8008"
    fi
    
    log_message "Health checks completed"
}

# Function to show deployment summary
show_summary() {
    print_header "Deployment Summary"
    
    echo "Application: $APP_NAME"
    echo "Timestamp: $(date)"
    echo ""
    echo "Services Status:"
    docker-compose ps
    echo ""
    echo "Access URLs:"
    echo "  Frontend: http://localhost:3000"
    echo "  Backend:  http://localhost:8008"
    echo "  Nginx:    http://localhost"
    echo ""
    echo "Useful Commands:"
    echo "  View logs:    docker-compose logs -f"
    echo "  Stop services: docker-compose down"
    echo "  Restart:      docker-compose restart"
    echo ""
    
    log_message "Deployment completed successfully"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --no-backup    Skip backup creation"
    echo "  --no-cleanup   Skip Docker cleanup"
    echo "  --production   Use production configuration"
    echo "  --help         Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Full deployment with backup and cleanup"
    echo "  $0 --no-backup       # Deploy without creating backup"
    echo "  $0 --production       # Deploy with production settings"
}

# Main deployment function
main() {
    print_header "SocialPay Deployment Script"
    
    # Parse command line arguments
    SKIP_BACKUP=false
    SKIP_CLEANUP=false
    PRODUCTION=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --no-backup)
                SKIP_BACKUP=true
                shift
                ;;
            --no-cleanup)
                SKIP_CLEANUP=true
                shift
                ;;
            --production)
                PRODUCTION=true
                shift
                ;;
            --help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Start deployment process
    log_message "Starting deployment with options: backup=$SKIP_BACKUP, cleanup=$SKIP_CLEANUP, production=$PRODUCTION"
    
    check_docker
    
    if [ "$SKIP_BACKUP" = false ]; then
        create_backup
    fi
    
    stop_services
    
    if [ "$SKIP_CLEANUP" = false ]; then
        cleanup_docker
    fi
    
    # Use production compose file if specified
    if [ "$PRODUCTION" = true ]; then
        if [ -f "docker-compose.prod.yml" ]; then
            export COMPOSE_FILE="docker-compose.yml:docker-compose.prod.yml"
            print_status "Using production configuration"
        else
            print_warning "Production compose file not found, using default"
        fi
    fi
    
    start_services
    health_check
    show_summary
    
    print_status "Deployment completed successfully!"
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi 