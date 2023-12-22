hello:
	@echo "Hello World"

start_mq:
	@echo "Starting up rabbitmq..."
	@docker compose up -d

start_server:
	@echo "Starting the server..."
	@go run http_server/*.go

start_log_consumer:
	@echo "Starting the consumer..."
	@go run log_consumer/*.go