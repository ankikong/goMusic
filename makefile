# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -ldflags "-s -w"
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

RELEASE_VERSION=V0.2.0
BINARY_NAME=goMusic-$(RELEASE_VERSION)
LOCAL_BINARY_NAME=goMusic
BINARY_DIR=bin


all: test build
build: 
	$(GOBUILD)
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOGET) github.com/jedib0t/go-pretty/table

mybuild:
	$(GOBUILD) -o D:/cmd/$(LOCAL_BINARY_NAME).exe

# Cross compilation

# windows
build-win64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_win64.exe
build-win32:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_win32.exe

# Linux
build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_linux_amd64
build-linux-X86:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_linux_X86
build-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_linux_arm64
build-linux-armV7:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_linux_armV7

# Mac
build-darwin-X86:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_darwin_X86
build-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)_darwin_amd64

build-all: build-win64 build-win32 build-linux-amd64 build-linux-X86 build-linux-arm64 build-linux-armV7 build-darwin-X86 build-darwin-amd64
