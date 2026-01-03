# Variables
APP_NAME=backend

# Build the app
build:
	go build -o $(APP_NAME) .

# Run the app
run:
	go run .

# Tidy dependencies
tidy:
	go mod tidy

# Clean build output
clean:
	rm -f $(APP_NAME)

cleandb:
	docker compose down -v

initdb:
	docker compose up -d

	