.PHONY: help dev-local dev-hybrid dev-docker services-only docker-up docker-down docker-logs clean test install-air

help:
	@echo "🚀 Available Development Commands:"
	@echo ""
	@echo "  dev-local      - Local development (PostgreSQL:5432 + RabbitMQ local + Air)"
	@echo "  dev-hybrid     - Hybrid development (PostgreSQL:5433 Docker + RabbitMQ Docker + Air local)"
	@echo "  dev-docker     - Full Docker development with hot reload"
	@echo "  services-only  - Start only PostgreSQL:5433 and RabbitMQ with Docker"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  docker-up      - Start production Docker environment"
	@echo "  docker-down    - Stop all Docker services"
	@echo "  docker-logs    - Show Docker logs"
	@echo ""
	@echo "🧹 Utility Commands:"
	@echo "  clean          - Clean up Docker resources"
	@echo "  test           - Run tests"
	@echo "  install-air    - Install air for hot reload"
	@echo ""
	@echo "📊 Port Summary:"
	@echo "  Local PostgreSQL:  5432"
	@echo "  Docker PostgreSQL: 5433"
	@echo "  RabbitMQ:          5672"
	@echo "  RabbitMQ UI:       15672"
	@echo "  Application:       8080"

# Local development dengan PostgreSQL lokal (port 5432)
dev-local:
	@chmod +x scripts/dev-local.sh
	@./scripts/dev-local.sh

# Hybrid: PostgreSQL & RabbitMQ di Docker, app dengan air
dev-hybrid:
	@chmod +x scripts/dev-hybrid.sh
	@./scripts/dev-hybrid.sh

# Full Docker development dengan hot reload
dev-docker:
	@chmod +x scripts/dev-docker.sh
	@./scripts/dev-docker.sh

# Start hanya services (PostgreSQL:5433 & RabbitMQ) untuk local development
services-only:
	@echo "🐳 Starting PostgreSQL and RabbitMQ with Docker..."
	docker-compose up -d postgres rabbitmq
	@echo ""
	@echo "✅ Services started:"
	@echo "📊 PostgreSQL: localhost:5433"
	@echo "🐰 RabbitMQ: localhost:5672 (Management: http://localhost:15672)"
	@echo ""
	@echo "💡 To use these services with your local app:"
	@echo "   1. Copy .env.hybrid to .env"
	@echo "   2. Run 'air' or 'go run cmd/main.go'"
	@echo ""
	@echo "   Or simply run: make dev-hybrid"

# Production Docker
docker-up:
	docker-compose up -d --build

# Stop Docker services
docker-down:
	docker-compose down

# Show Docker logs
docker-logs:
	docker-compose logs -f

# Clean up Docker resources
clean:
	docker-compose down -v
	docker system prune -f

# Run tests
test:
	go test ./...

# Install air
install-air:
	go install github.com/cosmtrek/air@latest
	@echo "✅ Air installed successfully"

# Quick status check
status:
	@echo "🔍 Service Status:"
	@echo ""
	@echo "Local PostgreSQL (5432):"
	@pg_isready -h localhost -p 5432 -U postgres 2>/dev/null && echo "  ✅ Running" || echo "  ❌ Not running"
	@echo ""
	@echo "Docker PostgreSQL (5433):"
	@pg_isready -h localhost -p 5433 -U postgres 2>/dev/null && echo "  ✅ Running" || echo "  ❌ Not running"
	@echo ""
	@echo "RabbitMQ (5672):"
	@curl -s http://localhost:15672 >/dev/null 2>&1 && echo "  ✅ Running" || echo "  ❌ Not running"
	@echo ""
	@echo "Docker containers:"
	@docker-compose ps 2>/dev/null || echo "  No Docker containers running"
