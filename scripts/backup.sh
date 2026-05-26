#!/bin/sh
# Database backup script
# Usage: ./scripts/backup.sh [output_dir]

set -e

BACKUP_DIR="${1:-./backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DB_URL="${DATABASE_URL:-postgres://livehouse:livehouse_dev@localhost:5432/livehouse_dev?sslmode=disable}"

mkdir -p "$BACKUP_DIR"

echo "Backing up database to $BACKUP_DIR/backup_$TIMESTAMP.sql"

# Extract connection details from URL
# postgres://user:pass@host:port/db?params
pg_dump "$DB_URL" \
  --no-owner \
  --no-acl \
  --format=custom \
  --file="$BACKUP_DIR/backup_$TIMESTAMP.dump"

echo "Backup complete: $BACKUP_DIR/backup_$TIMESTAMP.dump"
echo "Size: $(du -h "$BACKUP_DIR/backup_$TIMESTAMP.dump" | cut -f1)"

# Keep only last 7 backups
ls -t "$BACKUP_DIR"/backup_*.dump 2>/dev/null | tail -n +8 | xargs -r rm

echo "Cleaned up old backups"
