# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=clido
PACKAGE=github.com/d4r1us-drk/clido

# Build directory
BUILD_DIR=build

# Git information
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags
LDFLAGS=-ldflags "-X $(PACKAGE)/internal/version.GitCommit=$(GIT_COMMIT) -X $(PACKAGE)/internal/version.BuildDate=$(BUILD_DATE)"

# Platforms to build for
PLATFORMS=windows/amd64 darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

# Tools
GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)
GOFUMPT := $(shell command -v gofumpt 2> /dev/null)
GOIMPORTS := $(shell command -v goimports 2> /dev/null)
GOLINES := $(shell command -v golines 2> /dev/null)

.PHONY: all build clean deps lint format build-all version install uninstall help

all: deps lint format build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

clean:
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

deps:
	@echo "Checking and updating dependencies..."
	@go mod tidy
	@if [ -z "$$(git status --porcelain go.mod go.sum)" ]; then \
		echo "No missing dependencies. All modules are up to date."; \
	else \
		echo "Dependencies updated. Please review changes in go.mod and go.sum."; \
	fi
ifndef GOLANGCI_LINT
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.60.1
endif
ifndef GOFUMPT
	@echo "Installing gofumpt..."
	@go install mvdan.cc/gofumpt@latest
endif
ifndef GOIMPORTS
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
endif
ifndef GOLINES
	@echo "Installing golines..."
	@go install github.com/segmentio/golines@latest
endif

lint:
	@echo "Running linter..."
	@golangci-lint run --fix -c .golangci.yml ./...


format:
	@echo "Formatting code..."
	@gofumpt -l -w .
	@golines -l -m 120 -t 4 -w .
	@golines -w .
	echo "Code formatted."; \

build-all:
	mkdir -p $(BUILD_DIR)
	$(foreach PLATFORM,$(PLATFORMS),\
		$(eval GOOS=$(word 1,$(subst /, ,$(PLATFORM))))\
		$(eval GOARCH=$(word 2,$(subst /, ,$(PLATFORM))))\
		$(eval EXTENSION=$(if $(filter $(GOOS),windows),.exe,))\
		$(eval CGO_ENABLED=$(if $(filter $(GOOS),windows),1,0))\
		$(eval CC=$(if $(filter $(GOOS),windows),x86_64-w64-mingw32-gcc,))\
		GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) CC=$(CC) $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(EXTENSION) .;\
	)

# Version information
version:
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Install the application
install:
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) .

# Uninstall the application
uninstall:
	@rm $(GOPATH)/bin/$(BINARY_NAME) 

# Installation help
help:
	@echo "Available commands:"
	@echo "  make              - Run deps, lint, format, test, and build"
	@echo "  make build        - Build for the current platform"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make deps         - Download dependencies and install tools"
	@echo "  make lint         - Run golangci-lint for linting"
	@echo "  make format       - Format code using gofumpt, goimports, and golines"
	@echo "  make build-all    - Build for all specified platforms"
	@echo "  make version      - Display the current git commit and build date"
	@echo "  make install      - Install the application to GOPATH/bin"
	@echo "  make help         - Display this help information"
