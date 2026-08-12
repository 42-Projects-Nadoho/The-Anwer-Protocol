<<<<<<< HEAD
# TAP — The Answer Protocol
# Build tool: GNU Make driving the Go toolchain (standard library only).

GO      ?= go
BINDIR  := bin
ADDR    ?= 127.0.0.1:4242
HTTP    ?= 127.0.0.1:8080
WORLD   ?= data/world.json

.PHONY: all deps build server cli gui run run-server run-cli run-gui lint fmt vet test clean help

all: build

## deps: resolve module dependencies (none beyond the standard library)
deps:
	$(GO) mod tidy

## build: compile the server, CLI client and GUI client into ./bin
build: server cli gui

server:
	$(GO) build -o $(BINDIR)/tap-server ./cmd/server

cli:
	$(GO) build -o $(BINDIR)/tap-cli ./cmd/cli

gui:
	$(GO) build -o $(BINDIR)/tap-gui ./cmd/gui

## run-server: build and run the server (ADDR, WORLD overridable)
run-server: server
	$(BINDIR)/tap-server -addr $(ADDR) -world $(WORLD)

## run-cli: build and run the CLI client (ADDR overridable)
run-cli: cli
	$(BINDIR)/tap-cli -addr $(ADDR)

## run-gui: build and run the GUI client (ADDR, HTTP overridable)
run-gui: gui
	$(BINDIR)/tap-gui -addr $(ADDR) -http $(HTTP)

## lint: vet the code and check gofmt compliance
lint: vet
	@echo "checking gofmt..."
	@unformatted=$$(gofmt -l cmd internal); \
	if [ -n "$$unformatted" ]; then echo "needs gofmt:"; echo "$$unformatted"; exit 1; fi
	@echo "gofmt: clean"

## fmt: format all Go source
fmt:
	gofmt -w cmd internal

vet:
	$(GO) vet ./...

## test: run the Go test suite
test:
	$(GO) test ./...

## clean: remove build artifacts
clean:
	rm -rf $(BINDIR)

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'
=======
BINARY_DIR := bin
SERVER_BIN := $(BINARY_DIR)/tap-server
CLI_BIN    := $(BINARY_DIR)/tap-cli
GUI_BIN    := $(BINARY_DIR)/tap-gui

.PHONY: all deps build build-server build-cli build-gui \
        run-server run-client run-client-gui \
        lint fmt vet test clean

all: build

deps:
	go mod download
	go mod tidy

build: build-server build-cli build-gui

build-server:
	go build -o $(SERVER_BIN) ./cmd/server

build-cli:
	go build -o $(CLI_BIN) ./cmd/cli

build-gui:
	go build -o $(GUI_BIN) ./cmd/gui

run-server: build-server
	./$(SERVER_BIN)

run-client: build-cli
	./$(CLI_BIN)

run-client-gui: build-gui
	./$(GUI_BIN)

lint: vet
	@test -z "$$(gofmt -l .)" || (echo "gofmt: files need formatting:" && gofmt -l . && exit 1)

fmt:
	gofmt -w .

vet:
	go vet ./...

test:
	go test ./...

clean:
	rm -rf $(BINARY_DIR)
>>>>>>> 578eb0dc3e417098fe3f87155ee8dabc00ab5696
