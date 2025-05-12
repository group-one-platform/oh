include .make/*.mk

.PHONY: run build build-win build-linux build-mac-apple build-mac-intel build-mac-universal

OUTPUT_DIR := bin
APP_NAME := oh
MAC_ARM64_BINARY := $(OUTPUT_DIR)/$(APP_NAME)-mac-arm64
MAC_AMD64_BINARY := $(OUTPUT_DIR)/$(APP_NAME)-mac-amd64
MAC_UNIVERSAL_BINARY := $(OUTPUT_DIR)/$(APP_NAME)-mac-universal
WIN_BINARY := $(OUTPUT_DIR)/$(APP_NAME)-win
LINUX_BINARY := $(OUTPUT_DIR)/$(APP_NAME)-linux

run:: ##@Run Run the application
	go run main.go $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

build:: ##@Build Build the binary for your current OS
	go build -o bin/oh

build-all:: build-linux build-win build-mac-universal ##@Build	Build all binaries (Linux, Windows, and macOS)

build-win:: ##@Build Build the binary for Windows platform
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(WIN_BINARY)

build-mac-apple:: ##@Build Build the binary for ARM64 (Apple Silicone)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -buildvcs=false -buildmode=pie -ldflags="-s -w" -o $(MAC_ARM64_BINARY) 

build-mac-intel:: ##@Build Build the binary for AMD64 (Intel)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -buildmode=pie -ldflags="-s -w" -o $(MAC_AMD64_BINARY)

build-mac-universal: build-mac-apple build-mac-intel ##@Build Build Apple Universal Binary
	lipo -create -output $(MAC_UNIVERSAL_BINARY) $(MAC_ARM64_BINARY) $(MAC_AMD64_BINARY)

build-linux:: ##@Build Build the binary for Linux platform
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -buildmode=pie -ldflags="-s -w" -o $(LINUX_BINARY)

clean:: ##@Cleanup Remove all build artifacts
	rm -rf $(OUTPUT_DIR)
