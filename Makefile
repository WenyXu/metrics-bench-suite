# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

.PHONY: all fmt lint

BINS := $(patsubst ./cmd/%/main.go,%,$(wildcard ./cmd/*/main.go))

# Build all binaries
all: $(BINS)

lint-all: fmt lint

# Pattern rule to build each binary
$(BINS):
	$(GOGET) ./cmd/$@
	$(GOBUILD) -o bin/$@ ./cmd/$@

# Format code
fmt:
	$(GOFMT) ./...

# Lint code
lint:
	golint ./...