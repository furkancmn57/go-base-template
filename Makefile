APP_NAME := go-base-template
MAIN_PATH := ./src/main.go

.PHONY: run build tidy openapi docker-up docker-down test vet

run: ## Run the API locally
	go run $(MAIN_PATH)

build: ## Build the single binary
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

tidy: ## Tidy go.mod/go.sum
	go mod tidy

openapi: ## Regenerate src/docs from controller annotations (commit the result)
	@command -v swag >/dev/null || { echo "install swag: go install github.com/swaggo/swag/cmd/swag@latest"; exit 1; }
	swag init -g src/main.go -o src/docs --parseDependency --parseInternal --outputTypes go,json

docker-up: ## Start Postgres, Redis and RabbitMQ for local development
	docker compose up -d

docker-down: ## Stop local infrastructure containers
	docker compose down

test: ## Run the test suite
	go test ./...

vet: ## Static analysis
	go vet ./...
