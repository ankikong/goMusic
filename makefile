# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

RELEASE_VERSION=V0.0.1
BINARY_NAME=goMusic-$(RELEASE_VERSION)
BINARY_DIR=bin


all: test build
build: 
		$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME).exe -v
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

deps:
	$(GOGET) github.com/jedib0t/go-pretty/table

# Cross compilation

# windows
build-win64:
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_win64.exe -v
build-win32:
		CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_win32.exe -v

# Linux
build-linux-amd64:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_linux_amd64 -v
build-linux-X86:
		CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_linux_X86 -v
build-linux-arm64:
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_linux_arm64 -v
build-linux-armV7:
		CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_linux_armV7 -v

# Mac
build-darwin-X86:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_darwin_X86 -v
build-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_darwin_amd64 -v

build-all:
	make build-win64
	make build-win32
	make build-linux-amd64
	make build-linux-X86
	make build-linux-arm64
	make build-linux-armV7
	make build-darwin-X86
	make build-darwin-amd64
