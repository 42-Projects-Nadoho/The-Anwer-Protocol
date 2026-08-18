# ============================================================
#  VARIABLES
# ============================================================

GO      ?= go
RUN	:= $(GO) run
BINDIR  := bin
ADDR    ?= 127.0.0.1:4242
HTTP    ?= 127.0.0.1:8080
WORLD   ?= tap/data/world.yaml

# ------------------------------------------------------------
#  Ansi Colors
# ------------------------------------------------------------

RESET    := \033[0m
BLUE     := \033[1;94m
CYAN     := \033[1;96m
GREEN    := \033[1;92m
YELLOW   := \033[1;93m
MAGENTA  := \033[1;95m
ECHO     := echo -e

# ============================================================
#  RULES
# ============================================================

.PHONY: all install lint clean help \
	run-server run-client run-client-gui \
	deps build server cli gui run fmt vet \
	test-concurrency test-race test-group \
	test-coalescing test-fragmentation test-abuse

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
	$(GO) build -o $(BINDIR)/tap-server ./tap/src/cmd/server

cli:
	@$(ECHO) ">>> $(YELLOW)Building TAP CLI client...$(RESET)"
	$(GO) build -o $(BINDIR)/tap-cli ./tap/src/cmd/cli

gui:
	@$(ECHO) ">>> $(YELLOW)Building TAP GUI client...$(RESET)"
	$(GO) build -o $(BINDIR)/tap-gui ./tap/src/cmd/gui

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
	@unformatted=$$(gofmt -l tap/src/cmd tap/src/internal); \
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
	gofmt -w tap/src/cmd tap/src/internal
	@$(ECHO) ">>> $(CYAN)Done.$(RESET)"

# ------------------------------------------------------------
#  vet — run go vet
# ------------------------------------------------------------

vet:
	@$(ECHO) ">>> $(YELLOW)Running go vet...$(RESET)"
	$(GO) vet ./tap/...
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
	@$(ECHO) " $(MAGENTA)MANDATORY RULES$(RESET)"
	@$(ECHO) "     $(MAGENTA)install$(RESET)		Resolve module dependencies"
	@$(ECHO) "     $(MAGENTA)build$(RESET)            	Compile the server, CLI client and GUI client"
	@$(ECHO) "     $(MAGENTA)run-server$(RESET)       	Build and run the server"
	@$(ECHO) "     $(MAGENTA)run-client$(RESET)       	Build and run the CLI client"
	@$(ECHO) "     $(MAGENTA)run-client-gui$(RESET)	Build and run the GUI client"
	@$(ECHO) "     $(MAGENTA)lint$(RESET)             	Vet the code and check gofmt compliance"
	@$(ECHO) "     $(MAGENTA)clean$(RESET)            	Remove build artifacts"
	@$(ECHO) ""
	@$(ECHO) " $(YELLOW)HELPER RULES$(RESET)"
	@$(ECHO) "     $(YELLOW)all$(RESET)              	Default target (build)"
	@$(ECHO) "     $(YELLOW)deps$(RESET)             	Alias for install"
	@$(ECHO) "     $(YELLOW)fmt$(RESET)              	Format all Go source"
	@$(ECHO) "     $(YELLOW)vet$(RESET)              	Run go vet"
	@$(ECHO) "     $(YELLOW)test$(RESET)             	Run the standard Go test suite"
	@$(ECHO) ""
	@$(ECHO) " $(CYAN)TESTING SCRIPTS$(RESET)"
	@$(ECHO) "     $(CYAN)test-concurrency$(RESET)   Run the high-concurrency stress test"
	@$(ECHO) "     $(CYAN)test-race$(RESET)          Run the simultaneous movement race condition test"
	@$(ECHO) "     $(CYAN)test-group$(RESET)         Run the group state volatility test"
	@$(ECHO) "     $(CYAN)test-coalescing$(RESET)    Run the TCP packet coalescing test"
	@$(ECHO) "     $(CYAN)test-fragmentation$(RESET) Run the TCP packet fragmentation test"
	@$(ECHO) "     $(CYAN)test-abuse$(RESET)         Run the command flooding and abuse test"
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

test-coalescing:
	@$(ECHO) ">>> $(YELLOW)Running TCP Coalescing Test...$(RESET)"
	$(RUN) scripts/test_coalescing.go

test-fragmentation:
	@$(ECHO) ">>> $(YELLOW)Running TCP Fragmentation Test...$(RESET)"
	$(RUN) scripts/test_fragmentation.go

test-abuse:
	@$(ECHO) ">>> $(YELLOW)Running Server Abuse Test...$(RESET)"
	$(RUN) scripts/test_abuse.go
