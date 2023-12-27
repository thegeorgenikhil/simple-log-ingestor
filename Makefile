hello:
	@echo "Hello World"

format:
	@echo "Formatting code..."
	@go fmt ./...

start_mq:
	@echo "Starting up rabbitmq..."
	@docker compose up -d

start_server:
	@echo "Starting the server..."
	@go run http_server/*.go

start_alert_service:
	@echo "Starting the alert service..."
	@go run alert_service/*.go

start_logger_service:
	@echo "Starting the logger service..."
	@go run logger_service/*.go