# Makefile for snaptrack

# The name of your final binary:
APP_NAME  = snapstack

# The path to the main package for your CLI/TUI (within cmd/):
CMD_PATH  = ./cmd/$(APP_NAME)

# Optional: path for the final built binary:
BIN_DIR   = bin

.PHONY: all build run test clean

all: build

## Build the Go binary and place it in ./bin
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_PATH)

## Run the built binary
run: build
	@echo "Running $(APP_NAME)..."
	./$(BIN_DIR)/$(APP_NAME)

## Run all Go tests in the project
test:
	@echo "Running tests..."
	go test ./... -v

## Clean removes the bin folder (the built binary)
clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
