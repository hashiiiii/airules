# Define variables
BINARY      := airules
BUILD_DIR   := bin
ENTRY_PATH  := main.go
GO_VERSION  := $(shell go version | cut -d ' ' -f 3 | sed 's/go//')

# Get version from git tags (format: v1.2.3)
# If no tag exists, use 'unknown' as version
GIT_TAG     := $(shell git tag -l "v[0-9]*" | sort -V | tail -n 1)
VERSION     := $(if $(GIT_TAG),$(shell echo $(GIT_TAG) | sed 's/^v//'),unknown)
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE        ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS     := -ldflags "-X github.com/hashiiiii/airules/pkg/version.Version=$(VERSION) -X github.com/hashiiiii/airules/pkg/version.Commit=$(COMMIT) -X github.com/hashiiiii/airules/pkg/version.BuildDate=$(DATE)"
GOPATH      := $(shell go env GOPATH)

# Set default target to help
.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help
	@echo "-----------------------------------------------------------"
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_.-]+:.*?##/ { \
		printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 \
	}' $(MAKEFILE_LIST)
	@echo "-----------------------------------------------------------"

.PHONY: info
info: ## Display environment variables
	@echo "-----------------------------------------------------------"
	@echo "Environment variables:"
	@printf "  \033[36m%-16s\033[0m %s\n" "BINARY" "$(BINARY)"
	@printf "  \033[36m%-16s\033[0m %s\n" "BUILD_DIR" "$(BUILD_DIR)"
	@printf "  \033[36m%-16s\033[0m %s\n" "ENTRY_PATH" "$(ENTRY_PATH)"
	@printf "  \033[36m%-16s\033[0m %s\n" "GO_VERSION" "$(GO_VERSION)"
	@printf "  \033[36m%-16s\033[0m %s\n" "VERSION" "$(VERSION)"
	@printf "  \033[36m%-16s\033[0m %s\n" "COMMIT" "$(COMMIT)"
	@printf "  \033[36m%-16s\033[0m %s\n" "DATE" "$(DATE)"
	@printf "  \033[36m%-16s\033[0m %s\n" "GOPATH" "$(GOPATH)"
	@echo "-----------------------------------------------------------"

.PHONY: all
all: clean build ## Run clean build

.PHONY: clean
clean: ## Clean the build directory
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) 2>/dev/null || true

.PHONY: build
build: ## Build the binary
	@echo "Building $(BINARY) with Go $(GO_VERSION)..."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) $(ENTRY_PATH)

.PHONY: deps
deps: ## Update dependencies
	@echo "Updating dependencies..."
	go mod tidy

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	go test -v -parallel 4 ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html


.PHONY: lint
lint: ## Run linter and fix issues
	@echo "Running linter with auto-fix..."
	gofmt -w -s .
	golangci-lint run --fix ./...
