# Variables
APP_NAME=personal_finance_management

# Build the app
mac:
	go build -o $(APP_NAME) .

windows:
	GOOS=windows GOARCH=amd64 go build -o $(APP_NAME).exe .

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

