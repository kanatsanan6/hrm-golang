#!/bin/sh

set -e

echo "run db migration"
echo "postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable"
/app/migrate -path /app/migration -database "postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable" -verbose up

echo "start the app"
exec "$@"
