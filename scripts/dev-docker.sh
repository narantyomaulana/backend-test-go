#!/bin/bash

echo "🐳 Starting Full Docker Development Environment..."

# Build dan start semua services
docker-compose up -d --build

echo ""
echo "✅ Docker environment started"
echo "📊 PostgreSQL: localhost:5433"
echo "🐰 RabbitMQ: localhost:5672 (Management: http://localhost:15672)"
echo "🚀 Application: http://localhost:8080"
echo ""
echo "💡 To view logs: make docker-logs"
echo "💡 To stop: make docker-down"
