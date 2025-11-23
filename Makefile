BINARY_NAME=port_forwarder
CONFIG_FILE=config.yaml

.PHONY: all build run clean deps help

all: build ## Build the application

build: ## Build the binary
	go build -o $(BINARY_NAME) cmd/port_forwarder/main.go

run: build ## Run the application
	./$(BINARY_NAME) -config $(CONFIG_FILE)

clean: ## Remove the binary
	rm -f $(BINARY_NAME)

deps: ## Download dependencies
	go mod download

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'