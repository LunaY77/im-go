# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
SERVER_BINARY_NAME=imServer
CLIENT_BINARY_NAME=imClient

# Build paths
SERVER_SOURCE=./server
CLIENT_SOURCE=./client

.PHONY: all build clean server client help deps test

all: build

build: server client

server:
	$(GOBUILD) -o $(SERVER_BINARY_NAME) $(SERVER_SOURCE)

client:
	$(GOBUILD) -o $(CLIENT_BINARY_NAME) $(CLIENT_SOURCE)

clean:
	$(GOCLEAN)
	rm -f ./$(SERVER_BINARY_NAME) 2>/dev/null || true
	rm -f ./$(CLIENT_BINARY_NAME) 2>/dev/null || true

deps:
	$(GOMOD) download
	$(GOMOD) tidy

test:
	$(GOTEST) -v ./...

help:
	@echo "Available targets:"
	@echo "  all     - Build both server and client (default)"
	@echo "  build   - Build both server and client"
	@echo "  server  - Build only server"
	@echo "  client  - Build only client"
	@echo "  clean   - Clean build artifacts"
	@echo "  deps    - Download and tidy dependencies"
	@echo "  test    - Run tests"
	@echo "  help    - Show this help"
