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
