#!/bin/bash

# Script to convert existing migrations to golang-migrate format
# Converts from single SQL files to up/down pairs

set -e

SOURCE_DIR="migrations"
TARGET_DIR="migrations_new"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Create target directory
mkdir -p $TARGET_DIR

log_info "Converting migrations from $SOURCE_DIR to $TARGET_DIR..."

# Convert each migration
for file in $SOURCE_DIR/*.sql; do
    if [ -f "$file" ]; then
        # Extract number and name
        basename_file=$(basename "$file")
        number=$(echo "$basename_file" | sed 's/\([0-9]*\)_.*/\1/')
        name=$(echo "$basename_file" | sed 's/[0-9]*_\(.*\)\.sql/\1/')
        
        # Format number with leading zeros (6 digits for golang-migrate)
        formatted_number=$(printf "%06d" "$number")
        
        # Create up migration
        up_file="$TARGET_DIR/${formatted_number}_${name}.up.sql"
        cp "$file" "$up_file"
        
        # Create down migration based on content analysis
        down_file="$TARGET_DIR/${formatted_number}_${name}.down.sql"
        create_down_migration "$file" "$down_file" "$name"
        
        log_success "Converted: $basename_file → ${formatted_number}_${name}.{up,down}.sql"
    fi
done

log_success "All migrations converted successfully!"
log_info "Review the down migrations and test them before using in production"

# Function to create down migration based on up migration content
create_down_migration() {
    local up_file="$1"
    local down_file="$2"
    local name="$3"
    
    # Start down migration file
    cat > "$down_file" << EOF
-- Down migration: Revert $name
-- This undoes the changes made in the corresponding up migration

EOF
    
    # Analyze the up migration and create appropriate down commands
    if grep -q "CREATE TABLE" "$up_file"; then
        # Extract table names and create DROP statements
        grep "CREATE TABLE" "$up_file" | sed 's/.*CREATE TABLE[^A-Za-z]*\([A-Za-z_]*\).*/DROP TABLE IF EXISTS \1 CASCADE;/' >> "$down_file"
    fi
    
    if grep -q "ALTER TABLE.*ADD COLUMN" "$up_file"; then
        # Extract column additions and create DROP COLUMN statements
        grep "ALTER TABLE.*ADD COLUMN" "$up_file" | while read line; do
            table=$(echo "$line" | sed 's/.*ALTER TABLE[^A-Za-z]*\([A-Za-z_]*\).*/\1/')
            column=$(echo "$line" | sed 's/.*ADD COLUMN[^A-Za-z]*\([A-Za-z_]*\).*/\1/')
            echo "ALTER TABLE $table DROP COLUMN IF EXISTS $column;" >> "$down_file"
        done
    fi
    
    if grep -q "CREATE INDEX" "$up_file"; then
        # Extract index names and create DROP statements
        grep "CREATE INDEX" "$up_file" | sed 's/.*CREATE INDEX[^A-Za-z]*\([A-Za-z_]*\).*/DROP INDEX IF EXISTS \1;/' >> "$down_file"
    fi
    
    if grep -q "CREATE TYPE" "$up_file"; then
        # Extract enum types and create DROP statements
        grep "CREATE TYPE" "$up_file" | sed 's/.*CREATE TYPE[^A-Za-z]*\([A-Za-z_]*\).*/DROP TYPE IF EXISTS \1 CASCADE;/' >> "$down_file"
    fi
    
    if grep -q "CREATE.*FUNCTION" "$up_file"; then
        # Extract function names and create DROP statements
        grep "CREATE.*FUNCTION" "$up_file" | sed 's/.*FUNCTION[^A-Za-z]*\([A-Za-z_]*\).*/DROP FUNCTION IF EXISTS \1() CASCADE;/' >> "$down_file"
    fi
    
    if grep -q "CREATE TRIGGER" "$up_file"; then
        # Extract trigger names and create DROP statements
        grep "CREATE TRIGGER" "$up_file" | while read line; do
            trigger=$(echo "$line" | sed 's/.*CREATE TRIGGER[^A-Za-z]*\([A-Za-z_]*\).*/\1/')
            table=$(echo "$line" | sed 's/.*ON[^A-Za-z]*\([A-Za-z_]*\).*/\1/')
            echo "DROP TRIGGER IF EXISTS $trigger ON $table;" >> "$down_file"
        done
    fi
    
    if grep -q "CREATE.*VIEW" "$up_file"; then
        # Extract view names and create DROP statements
        grep "CREATE.*VIEW" "$up_file" | sed 's/.*VIEW[^A-Za-z]*\([A-Za-z_]*\).*/DROP VIEW IF EXISTS \1 CASCADE;/' >> "$down_file"
    fi
    
    # Add comment if no operations were detected
    if [ ! -s "$down_file" ] || [ $(wc -l < "$down_file") -le 3 ]; then
        cat >> "$down_file" << EOF
-- TODO: Add appropriate down migration commands
-- Review the up migration and manually add the reverse operations
EOF
    fi
}

# Export function for use in subshell
export -f create_down_migration