# Makefile для проекта AI-Driven Ops Copilot

# --- Переменные ---
# Используем go run с флагом -mod=mod, чтобы он правильно работал с Go Workspaces
GO_RUN = go run -mod=mod

# --- Основные цели для запуска сервисов ---

.PHONY: run-auditor run-collector run-gateway run-user run-reasoner
.PHONY: run-all stop-all

# Запускает Auditor Service, который слушает Kafka и предоставляет gRPC API
run-auditor:
	@echo ">> Starting Auditor Service..."
	@$(GO_RUN) github.com/waste3d/ai-ops/services/auditor/cmd/auditor

# Запускает Collector Service, который принимает HTTP запросы
run-collector:
	@echo ">> Starting Collector Service..."
	@$(GO_RUN) github.com/waste3d/ai-ops/services/collector/cmd/collector

# Запускает API Gateway Service, который предоставляет внешний REST API
run-gateway:
	@echo ">> Starting API Gateway Service..."
	@$(GO_RUN) github.com/waste3d/ai-ops/services/api_gateway/cmd/api_gateway

# Запускает User Service, который предоставляет gRPC API для пользователей
run-user:
	@echo ">> Starting User Service..."
	@$(GO_RUN) github.com/waste3d/ai-ops/services/user_service/cmd/user_service

# Запускает AI Reasoner/Responser Service, который анализирует тикеты
# Используем имя папки ai_responser, как в вашем проекте
run-ai:
	@echo ">> Starting AI Responser Service..."
	@$(GO_RUN) github.com/waste3d/ai-ops/services/ai_responser/cmd/ai_responser

# --- Утилиты для разработки ---

.PHONY: proto-gen

# Генерирует Go-код из всех .proto файлов
proto-gen:
	@echo ">> Generating Protobuf code for auth.proto..."
	@protoc --proto_path=protos \
	  --go_out=./gen/go --go_opt=paths=source_relative \
	  --go-grpc_out=./gen/go --go-grpc_opt=paths=source_relative \
	  protos/auth.proto

	@echo ">> Generating Protobuf code for ticket.proto..."
	@protoc --proto_path=protos \
	  --go_out=./gen/go --go_opt=paths=source_relative \
	  --go-grpc_out=./gen/go --go-grpc_opt=paths=source_relative \
	  protos/ticket.proto
	
	@echo ">> Protobuf generation complete."


.PHONY: tidy

# Выполняет 'go mod tidy' для всех сервисов, чтобы привести зависимости в порядок
tidy:
	@echo ">> Tidying Go modules for all services..."
	@(cd services/auditor && go mod tidy)
	@(cd services/collector && go mod tidy)
	@(cd services/api_gateway && go mod tidy)
	@(cd services/user_service && go mod tidy)
	@(cd services/ai_responser && go mod tidy)
	@echo ">> Go modules tidied."

# --- Управление Docker Compose ---

.PHONY: up down logs

# Запускает все инфраструктурные сервисы (Kafka, Postgres, etc.) в фоновом режиме
up:
	@echo ">> Starting Docker infrastructure..."
	@docker-compose up -d

# Останавливает и удаляет все инфраструктурные контейнеры
down:
	@echo ">> Stopping Docker infrastructure..."
	@docker-compose down

# Показывает логи всех инфраструктурных сервисов
logs:
	@echo ">> Tailing Docker infrastructure logs..."
	@docker-compose logs -f

# --- Помощь ---

.PHONY: help

# Отображает список доступных команд
help:
	@echo "Available commands:"
	@echo "  make up              - Start all Docker services (Kafka, Postgres, etc.)"
	@echo "  make down            - Stop and remove all Docker services"
	@echo "  make logs            - Tail logs from Docker services"
	@echo ""
	@echo "  make run-auditor     - Run the Auditor service"
	@echo "  make run-collector   - Run the Collector service"
	@echo "  make run-gateway     - Run the API Gateway service"
	@echo "  make run-user        - Run the User service"
	@echo "  make run-reasoner    - Run the AI Reasoner service"
	@echo ""
	@echo "  make proto-gen       - Generate Go code from .proto files"
	@echo "  make tidy            - Run 'go mod tidy' for all services"