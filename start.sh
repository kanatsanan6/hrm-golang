#!/bin/sh

set -e

echo "run db migration"
source .env
/app/migrate -path /app/migration -database "postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable" -verbose up

echo "start the app"
exec "$@"
