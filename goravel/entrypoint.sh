#!/bin/sh
set -e

echo "🔑 Generating app key..."
go run . artisan key:generate || true

echo "🗄️ Running migrations..."
go run . artisan migrate || true

echo "🚀 Starting application..."
exec go run main.go
