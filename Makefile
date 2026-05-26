.PHONY: dev backend frontend db sqlc migrate-up migrate-down test

# Start all services in dev mode
dev:
	docker compose -f docker-compose.dev.yml up --build

# Start only backend
backend:
	docker compose -f docker-compose.dev.yml up --build backend

# Start only frontend
frontend:
	docker compose -f docker-compose.dev.yml up --build frontend

# Start only database
db:
	docker compose -f docker-compose.dev.yml up -d postgres redis

# Generate sqlc Go code
sqlc:
	cd backend && PATH="$$(go env GOPATH)/bin:$$PATH" sqlc generate

# Run database migrations
migrate-up:
	cd backend && go run cmd/migrate/main.go up

migrate-down:
	cd backend && go run cmd/migrate/main.go down

# Run tests
test:
	cd backend && go test ./... -v
