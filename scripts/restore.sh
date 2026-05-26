#!/bin/sh
# Restore database from backup
# Usage: ./scripts/restore.sh <backup_file>

set -e

BACKUP_FILE="$1"
DB_URL="${DATABASE_URL:-postgres://livehouse:livehouse_dev@localhost:5432/livehouse_dev?sslmode=disable}"

if [ -z "$BACKUP_FILE" ]; then
  echo "Usage: $0 <backup_file>"
  exit 1
fi

if [ ! -f "$BACKUP_FILE" ]; then
  echo "Backup file not found: $BACKUP_FILE"
  exit 1
fi

echo "Restoring database from $BACKUP_FILE"
pg_restore "$DB_URL" \
  --no-owner \
  --no-acl \
  --clean \
  --if-exists \
  --dbname="$DB_URL" \
  "$BACKUP_FILE"

echo "Restore complete"
