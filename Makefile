
APP_NAME := pr-reviewer-service
CMD_DIR  := ./cmd/$(APP_NAME)
BIN_DIR  := ./bin
BIN_FILE := $(BIN_DIR)/$(APP_NAME)

.PHONY: fmt vet build clean run test lint dev

all: build

build:
	go build -o $(BIN_FILE) $(CMD_DIR)

run: build
	go run $(CMD_DIR)

# TODO TESTS
lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	@echo "→ Cleaning binary"
	@rm -rf $(BIN_DIR)
