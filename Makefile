BIN        := api
BUILD_DIR  := bin
COVER_OUT  := coverage.out
COVER_HTML := coverage.html

.PHONY: build test cover swagger mock tools clean help

build: ## Compile the binary to $(BUILD_DIR)/$(BIN)
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BIN) ./cmd/api

test: ## Run all tests
	go test ./...

cover: ## Run tests with coverage and print a summary
	go test -coverprofile=$(COVER_OUT) ./...
	go tool cover -func=$(COVER_OUT)

cover-html: ## Open an HTML coverage report in the browser
	go test -coverprofile=$(COVER_OUT) ./...
	go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	open $(COVER_HTML)

swagger: ## Regenerate docs/ from swag annotations
	swag init -g cmd/api/main.go -o docs

mock: ## Regenerate mocks/ via mockgen
	go generate ./...

tools: ## Install development tools (mockgen, swag)
	go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest

clean: ## Remove build artefacts
	rm -rf $(BUILD_DIR) $(COVER_OUT) $(COVER_HTML)

help: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  %-14s %s\n", $$1, $$2}'
