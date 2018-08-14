GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=fochoc
BINARY_WIN=$(BINARY_NAME)_win
BINARY_MAC=$(BINARY_NAME)_mac
BINARY_LINUX=$(BINARY_NAME)_linux

all: dep test build
build:
		env GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WIN) -v
		env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MAC) -v
		env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v
test:
		$(GOTEST) -v ./
clean:
		$(GOCLEAN)
		rm -f $(BINARY_WIN)
		rm -f $(BINARY_MAC)
		rm -f $(BINARY_LINUX)
dep:
		dep ensure