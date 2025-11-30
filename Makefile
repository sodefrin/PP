.PHONY: run build gen deps clean

# Default target
run:
	go run main.go

# Build the binary
build:
	mkdir -p bin
	go build -o bin/server main.go

# Generate code (sqlc)
gen:
	sqlc generate

# Install dependencies
deps:
	go mod download

# Clean build artifacts
clean:
	rm -rf bin
