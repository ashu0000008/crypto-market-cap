# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test


all: clean test build build-linux
build: 
				cd main && $(GOBUILD) -o crypto-collector -v main.go
test:
				cd main && $(GOTEST) -v ./...
clean:
				rm -f ./main/crypto-collector*

# Cross compilation
build-linux:
				CGO_ENABLED=0 GOOS=linux GOARCH=amd64  cd main && $(GOBUILD) -o crypto-collector-linux -v main.go