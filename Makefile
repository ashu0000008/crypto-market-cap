# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test


all: clean test build build-linux
build: 
				cd main && $(GOBUILD) -o ../output/crypto-collector -v main.go
				cd api && $(GOBUILD) -o ../output/crypto-api -v api.go
test:
				cd main && $(GOTEST) -v ./...
clean:
				rm -f ./output/*

# Cross compilation
build-linux:
				CGO_ENABLED=0 GOOS=linux GOARCH=amd64  cd main && $(GOBUILD) -o crypto-collector-linux -v main.go