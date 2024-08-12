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
STATICCHECK=staticcheck
GOFUMPT=gofumpt
GOIMPORTS=goimports
GOLINES=golines

.PHONY: all build clean test deps lint format build-all

all: deps lint format test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

clean:
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

deps:
	$(GOMOD) download
	$(GOGET) honnef.co/go/tools/cmd/staticcheck
	$(GOGET) mvdan.cc/gofumpt
	$(GOGET) golang.org/x/tools/cmd/goimports
	$(GOGET) github.com/segmentio/golines

lint:
	$(STATICCHECK) ./...

format:
	$(GOFUMPT) -l -w .
	$(GOIMPORTS) -w .
	$(GOLINES) -w .

# Probably not needed now

# build-all:
# 	mkdir -p $(BUILD_DIR)
# 	$(foreach PLATFORM,$(PLATFORMS),\
# 		$(eval GOOS=$(word 1,$(subst /, ,$(PLATFORM))))\
# 		$(eval GOARCH=$(word 2,$(subst /, ,$(PLATFORM))))\
# 		$(eval EXTENSION=$(if $(filter $(GOOS),windows),.exe,))\
# 		$(eval CGO_ENABLED=$(if $(filter $(GOOS),windows),1,0))\
# 		$(eval CC=$(if $(filter $(GOOS),windows),x86_64-w64-mingw32-gcc,))\
# 		GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) CC=$(CC) $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(EXTENSION) .;\
# 	)

# Version information
version:
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Install the application
install:
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) .

uninstall:
	@rm $(GOPATH)/bin/$(BINARY_NAME) 
