#!/bin/bash

echo "🚀 Starting Local Development Environment..."

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo "📦 Installing air..."
    go install github.com/cosmtrek/air@latest
fi

# Check if local PostgreSQL is running on port 5432
if ! pg_isready -h localhost -p 5432 -U postgres > /dev/null 2>&1; then
    echo "❌ Local PostgreSQL is not running on port 5432"
    echo "💡 Options:"
    echo "   1. Start local PostgreSQL service"
    echo "   2. Use 'make dev-hybrid' to use Docker PostgreSQL"
    echo "   3. Use 'make dev-docker' for full Docker environment"
    exit 1
fi

# Ensure we're using local environment
if [ ! -f .env ]; then
    echo "❌ .env file not found"
    exit 1
fi

echo "✅ Environment configured for local development"
echo "📊 PostgreSQL: localhost:5432 (local)"
echo "🐰 RabbitMQ: localhost:5672"
echo "🌐 App will run on: http://localhost:8080"
echo ""

# Start the application with air
echo "🔥 Starting application with hot reload..."
air -c .air.toml
