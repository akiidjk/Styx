# Variables
APP_NAME := styx
GO_CMD := go
GO_BUILD := $(GO_CMD) build
GO_TEST := $(GO_CMD) test
GO_FMT := $(GO_CMD) fmt
GO_LINT := golangci-lint run
GO_GENERATE := $(GO_CMD) generate
GO_CLEAN := $(GO_CMD) clean
GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')
BUILD_DIR := bin
LOG_DIR := /var/log/$(APP_NAME)
MAIN_FILE := main.go
GENERATE_DIR := internal/ebpf/generated
CONVERT_SCRIPT := convert_to_public.sh
OS := $(shell uname -s)

# Colors
RESET := \033[0m
BOLD := \033[1m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RED := \033[31m

# Default target
all: run

# Build the application
build: generate
	@echo -e "$(BLUE)🔨 Building $(APP_NAME)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@$(GO_BUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo -e "$(GREEN)✅ Build complete!$(RESET)"

# Run tests
test:
	@echo -e "$(YELLOW)🧪 Running tests...$(RESET)"
	@$(GO_TEST) ./...
	@echo -e "$(GREEN)✅ All tests passed!$(RESET)"

# Run linter
lint:
	@echo -e "$(BLUE)🔍 Running linter...$(RESET)"
	@$(GO_LINT)
	@echo -e "$(GREEN)✅ Linting complete!$(RESET)"

# Format the code
fmt:
	@echo -e "$(YELLOW)✨ Formatting code...$(RESET)"
	@$(GO_FMT) $(GO_FILES)
	@echo -e "$(GREEN)✅ Code formatted!$(RESET)"

# Generate code
generate:
	@echo -e "$(BLUE)⚙️  Running go generate...$(RESET)"
	@$(GO_GENERATE) ./$(GENERATE_DIR)
	@./scripts/$(CONVERT_SCRIPT)
	@echo -e "$(GREEN)✅ Code generation complete!$(RESET)"

# Clean build files
clean:
	@echo -e "$(RED)🧹 Cleaning up...$(RESET)"
	@$(GO_CLEAN)
	@rm -rf $(BUILD_DIR)
	@echo -e "$(GREEN)✅ Clean complete!$(RESET)"

# Clean logs
clean_logs:
	@echo -e "$(RED)🧹 Cleaning up logs...$(RESET)"
	@sudo rm -rf $(LOG_DIR)/*
	@echo -e "$(GREEN)✅ Clean complete!$(RESET)"

# Run the application
run: build
	@echo -e "$(BLUE)🚀 Running $(APP_NAME)...$(RESET)"
	@sudo $(BUILD_DIR)/$(APP_NAME) $(ARGS)

# Vendor dependencies
vendor:
	@echo -e "$(BLUE)📦 Tidying and downloading vendor dependencies...$(RESET)"
	@$(GO_CMD) mod tidy
	@$(GO_CMD) mod vendor
	@echo -e "$(GREEN)✅ Vendor dependencies ready!$(RESET)"

# Help menu
help:
	@echo -e "$(BOLD)Usage:$(RESET)"
	@echo -e "$(YELLOW)  make$(RESET)            Build the application (default)"
	@echo -e "$(YELLOW)  make build$(RESET)      Build the application"
	@echo -e "$(YELLOW)  make run$(RESET)        Run the application"
	@echo -e "$(YELLOW)  make test$(RESET)       Run tests"
	@echo -e "$(YELLOW)  make lint$(RESET)       Run linter"
	@echo -e "$(YELLOW)  make fmt$(RESET)        Format code"
	@echo -e "$(YELLOW)  make generate$(RESET)   Run go generate"
	@echo -e "$(YELLOW)  make clean$(RESET)      Clean build files"
	@echo -e "$(YELLOW)  make vendor$(RESET)     Vendor dependencies"

.PHONY: all build test lint fmt generate clean run vendor help
