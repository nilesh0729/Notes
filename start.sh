#!/bin/sh

set -e

echo "run db migration"
/usr/bin/migrate -path /app/internal/db/migrate_files -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
