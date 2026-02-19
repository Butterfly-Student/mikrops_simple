#!/bin/bash

# Database Migration Script for GEMBOK
# Usage: ./migrate.sh [up|down|status]

set -e

# Database configuration (can be overridden by environment variables)
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_NAME=${DB_NAME:-gembok_db}
DB_USER=${DB_USER:-gembok}
DB_PASS=${DB_PASSWORD:-root}

# Migration directory
MIGRATIONS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATION_FILE="${MIGRATIONS_DIR}/20240119120000_initial_schema.sql"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Check if migration file exists
check_migration_file() {
    if [ ! -f "$MIGRATION_FILE" ]; then
        log_error "Migration file not found: $MIGRATION_FILE"
        exit 1
    fi
    log_info "Migration file found: $MIGRATION_FILE"
}

# Check database connection
check_connection() {
    log_info "Checking database connection..."
    if ! mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" -e "SELECT 1" &> /dev/null; then
        log_error "Cannot connect to database. Please check your credentials."
        exit 1
    fi
    log_info "Database connection successful."
}

# Check if database exists
check_database() {
    log_info "Checking if database '$DB_NAME' exists..."
    if ! mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" -e "USE $DB_NAME" &> /dev/null; then
        log_warn "Database '$DB_NAME' does not exist."
        return 1
    fi
    log_info "Database '$DB_NAME' exists."
    return 0
}

# Create database if not exists
create_database() {
    log_info "Creating database '$DB_NAME'..."
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
    log_info "Database created successfully."
}

# Run UP migration
migrate_up() {
    log_info "Starting UP migration..."
    check_database || create_database
    
    log_info "Running migration: $(basename "$MIGRATION_FILE")"
    
    # Extract and run only the Up section
    awk '/^-- Up/,/^-- Down/' "$MIGRATION_FILE" | mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME"
    
    if [ $? -eq 0 ]; then
        log_info "Migration UP completed successfully!"
    else
        log_error "Migration UP failed!"
        exit 1
    fi
}

# Run DOWN migration
migrate_down() {
    log_warn "Starting DOWN migration (this will drop all tables)..."
    check_database || { log_error "Database '$DB_NAME' does not exist."; exit 1; }
    
    # Confirm before running
    read -p "Are you sure you want to drop all tables? (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        log_info "Migration DOWN cancelled."
        exit 0
    fi
    
    log_info "Running DOWN migration..."
    
    # Extract and run only the Down section
    awk '/^-- Down/,0' "$MIGRATION_FILE" | mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME"
    
    if [ $? -eq 0 ]; then
        log_info "Migration DOWN completed successfully!"
    else
        log_error "Migration DOWN failed!"
        exit 1
    fi
}

# Show migration status
show_status() {
    log_info "Migration Status:"
    log_info "Database: $DB_NAME@$DB_HOST:$DB_PORT"
    log_info "Migration File: $MIGRATION_FILE"
    
    if check_database; then
        echo ""
        log_info "Current Tables:"
        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "SHOW TABLES;" | tail -n +2 | while read table; do
            if [ -n "$table" ]; then
                count=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "SELECT COUNT(*) FROM \`$table\`" | tail -n 1)
                echo "  - $table ($count records)"
            fi
        done
    else
        log_warn "Database does not exist."
    fi
}

# Show help
show_help() {
    echo "Database Migration Script for GEMBOK"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  up      Run UP migration (create tables)"
    echo "  down    Run DOWN migration (drop tables)"
    echo "  status  Show migration status"
    echo "  help    Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  DB_HOST        Database host (default: localhost)"
    echo "  DB_PORT        Database port (default: 3306)"
    echo "  DB_NAME        Database name (default: gembok_db)"
    echo "  DB_USER        Database user (default: root)"
    echo "  DB_PASSWORD    Database password (default: rootpassword)"
    echo ""
    echo "Examples:"
    echo "  $0 up                    # Run UP migration"
    echo "  $0 down                  # Run DOWN migration"
    echo "  DB_HOST=192.168.1.10 $0 up  # Run UP on remote server"
}

# Main script logic
case "${1:-help}" in
    up)
        check_migration_file
        check_connection
        migrate_up
        ;;
    down)
        check_migration_file
        check_connection
        migrate_down
        ;;
    status)
        check_migration_file
        check_connection
        show_status
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        log_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac

exit 0
