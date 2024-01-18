.PHONY: all mac linux windows

GOFLAGS=-ldflags="-s -w"
BINARY_NAME=broterm
MAIN=cmd/broterm/main.go
OUTPUT_DIR=bin

all: mac linux windows

mac:
	mkdir -p $(OUTPUT_DIR)/OSX
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(OUTPUT_DIR)/OSX/$(BINARY_NAME) $(MAIN)

linux:
	mkdir -p $(OUTPUT_DIR)/Linux
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(OUTPUT_DIR)/Linux/$(BINARY_NAME) $(MAIN)

windows:
	mkdir -p $(OUTPUT_DIR)/Windows
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(OUTPUT_DIR)/Windows/$(BINARY_NAME).exe $(MAIN)