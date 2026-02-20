#!/bin/sh

set -e

echo "running migrate"
./migratE -path ./migrate -database "$DB_SOURCE" -verbose up

echo "doen running migrate"

echo "Starting TeslaBank API..."
./main