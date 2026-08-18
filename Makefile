# ============================================================
#  VARIABLES
# ============================================================

GO      ?= go
RUN	:= $(GO) run
BINDIR  := bin
ADDR    ?= 127.0.0.1:4242
HTTP    ?= 127.0.0.1:8080
WORLD   ?= data/world.yaml

# ------------------------------------------------------------
#  Ansi Colors
# ------------------------------------------------------------

RESET    := \033[0m
BLUE     := \033[1;94m
CYAN     := \033[1;96m
GREEN    := \033[1;92m
YELLOW   := \033[1;93m
ECHO     := echo -e

# ============================================================
#  RULES
# ============================================================

.PHONY: all install lint clean help \
	run-server run-client run-client-gui \
	deps build server cli gui run fmt vet \
	test_concurrency test_race_conditions test_group_volatility

# ------------------------------------------------------------
#  all — default target
# ------------------------------------------------------------

all: build

# ------------------------------------------------------------
#  deps — resolve module dependencies
# ------------------------------------------------------------

deps:
	@$(ECHO) ">>> $(YELLOW)Resolving module dependencies...$(RESET)"
	$(GO) mod tidy
	@$(ECHO) ">>> $(CYAN)Dependencies resolved.$(RESET)"

# ------------------------------------------------------------
#  install — install dependencies (alias for deps)
# ------------------------------------------------------------

install: deps

# ------------------------------------------------------------
#  build — compile the server, CLI client and GUI client
# ------------------------------------------------------------

build: server cli gui
	@$(ECHO) ">>> $(CYAN)All binaries built successfully.$(RESET)"

server:
	@$(ECHO) ">>> $(YELLOW)Building TAP server...$(RESET)"
	$(GO) build -o $(BINDIR)/tap-server ./cmd/server

cli:
	@$(ECHO) ">>> $(YELLOW)Building TAP CLI client...$(RESET)"
	$(GO) build -o $(BINDIR)/tap-cli ./cmd/cli

gui:
	@$(ECHO) ">>> $(YELLOW)Building TAP GUI client...$(RESET)"
	$(GO) build -o $(BINDIR)/tap-gui ./cmd/gui

# ------------------------------------------------------------
#  run-server — build and run the server
# ------------------------------------------------------------

run-server: server
	@$(ECHO) ">>> $(YELLOW)Starting TAP server...$(RESET)"
	$(BINDIR)/tap-server -addr $(ADDR) -world $(WORLD)

# ------------------------------------------------------------
#  run-client — build and run the CLI client
# ------------------------------------------------------------

run-client: cli
	@$(ECHO) ">>> $(YELLOW)Starting TAP CLI client...$(RESET)"
	$(BINDIR)/tap-cli -addr $(ADDR)

# ------------------------------------------------------------
#  run-client-gui — build and run the GUI client
# ------------------------------------------------------------

run-client-gui: gui
	@$(ECHO) ">>> $(YELLOW)Starting TAP GUI client...$(RESET)"
	$(BINDIR)/tap-gui -addr $(ADDR) -http $(HTTP)

# ------------------------------------------------------------
#  lint — vet the code and check gofmt compliance
# ------------------------------------------------------------

lint: vet
	@$(ECHO) ">>> $(YELLOW)Checking gofmt compliance...$(RESET)"
	@unformatted=$$(gofmt -l cmd internal); \
	if [ -n "$$unformatted" ]; then \
		$(ECHO) ">>> $(YELLOW)needs gofmt:$(RESET)"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	@$(ECHO) ">>> $(CYAN)gofmt: clean$(RESET)"

# ------------------------------------------------------------
#  fmt — format all Go source
# ------------------------------------------------------------

fmt:
	@$(ECHO) ">>> $(YELLOW)Formatting Go source files...$(RESET)"
	gofmt -w cmd internal
	@$(ECHO) ">>> $(CYAN)Done.$(RESET)"

# ------------------------------------------------------------
#  vet — run go vet
# ------------------------------------------------------------

vet:
	@$(ECHO) ">>> $(YELLOW)Running go vet...$(RESET)"
	$(GO) vet ./...
	@$(ECHO) ">>> $(CYAN)go vet completed.$(RESET)"

# ------------------------------------------------------------
#  test — run the Go test suite
# ------------------------------------------------------------

test:
	@$(ECHO) ">>> $(YELLOW)Running Go test suite...$(RESET)"
	$(GO) test ./...
	@$(ECHO) ">>> $(CYAN)Tests completed.$(RESET)"

# ------------------------------------------------------------
#  clean — remove build artifacts
# ------------------------------------------------------------

clean:
	@$(ECHO) ">>> $(YELLOW)Cleaning build artifacts...$(RESET)"
	rm -rf $(BINDIR)
	@$(ECHO) ">>> $(CYAN)Done.$(RESET)"

# ------------------------------------------------------------
#  help — list available targets
# ------------------------------------------------------------

help:
	@$(ECHO) ""
	@$(ECHO) " $(BLUE)AVAILABLE RULES$(RESET)"
	@$(ECHO) ""
	@$(ECHO) "     $(BLUE)all$(RESET)              Default target (build)"
	@$(ECHO) "     $(BLUE)deps$(RESET)             Resolve module dependencies"
	@$(ECHO) "     $(BLUE)install$(RESET)          Alias for deps"
	@$(ECHO) "     $(BLUE)build$(RESET)            Compile the server, CLI client and GUI client"
	@$(ECHO) "     $(BLUE)run-server$(RESET)       Build and run the server"
	@$(ECHO) "     $(BLUE)run-client$(RESET)       Build and run the CLI client"
	@$(ECHO) "     $(BLUE)run-client-gui$(RESET)   Build and run the GUI client"
	@$(ECHO) "     $(BLUE)lint$(RESET)             Vet the code and check gofmt compliance"
	@$(ECHO) "     $(BLUE)fmt$(RESET)              Format all Go source"
	@$(ECHO) "     $(BLUE)vet$(RESET)              Run go vet"
	@$(ECHO) "     $(BLUE)test$(RESET)             Run the Go test suite"
	@$(ECHO) "     $(BLUE)clean$(RESET)            Remove build artifacts"
	@$(ECHO) ""

# ------------------------------------------------------------
# Edge Case Test Scripts
# ------------------------------------------------------------

test-concurrency:
	@$(ECHO) ">>> $(YELLOW)Running Concurrency Test...$(RESET)"
	$(RUN) scripts/test_concurrency.go

test-race:
	@$(ECHO) ">>> $(YELLOW)Running Race Condition Test...$(RESET)"
	$(RUN) scripts/test_race_conditions.go

test-group:
	@$(ECHO) ">>> $(YELLOW)Running Group Volatility Test...$(RESET)"
	$(RUN) scripts/test_group_volatility.go
