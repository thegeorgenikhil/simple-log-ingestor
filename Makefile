hello:
	@echo "Hello World"

start:
	@echo "Starting the server"
	@docker compose up -d
	@go run http_server/*.go