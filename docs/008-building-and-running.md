# Building and Running

The project utilizes a centralized `Makefile` for dependency management, compilation, testing, and execution. This document details how each rule functions under the hood.

## Environment Variables
The `Makefile` supports overriding default variables to customize the build and run environment:
- `GO`: The Go binary to use (default: `go`).
- `ADDR`: The TCP address for the server and clients (default: `127.0.0.1:4242`).
- `HTTP`: The HTTP address for the GUI client's web server (default: `127.0.0.1:8080`).
- `WORLD`: The path to the world configuration file (default: `tap/data/world.yaml`).

## Mandatory Rules
These rules satisfy the core project requirements for building and running the ecosystem:

- **`make install`** (or `make deps`):
  Executes `go mod tidy` to securely resolve, download, and verify all required Go modules.
  
- **`make run-server`**:
  Depends on the `server` build target. Compiles `tap/src/cmd/server` to `bin/tap-server` and immediately executes it, passing the `-addr` and `-world` flags.

- **`make run-client`**:
  Depends on the `cli` build target. Compiles `tap/src/cmd/cli` to `bin/tap-cli` and immediately executes it, passing the `-addr` flag to connect to the server.

- **`make run-client-gui`**:
  Depends on the `gui` build target. Compiles `tap/src/cmd/gui` to `bin/tap-gui` and immediately executes it, passing both `-addr` (for the TAP TCP connection) and `-http` (for serving the web-based GUI).

- **`make lint`**:
  Runs a rigorous quality check. It first calls `go vet ./...` to detect suspicious constructs, followed by `gofmt -l` to ensure all files in `tap/src/cmd` and `tap/src/internal` adhere strictly to standard Go formatting. If unformatted files are found, the build fails.

- **`make clean`**:
  Forcefully removes the `bin/` directory and all compiled artifacts using `rm -rf bin`, ensuring a pristine workspace for the next build.

## Helper Rules
- **`make all`** (or just `make`): The default target. Calls `make build` to compile all three binaries (Server, CLI, GUI) into the `bin/` directory without running them.
- **`make fmt`**: Automatically applies standard formatting using `gofmt -w` to all source files.
- **`make vet`**: Runs standard `go vet` static analysis.
- **`make test`**: Runs the standard Go test suite (`go test ./...`) for unit testing.

## Testing Scripts
We included specific rules to execute our automated integration and edge-case testing scripts. Each rule dynamically compiles and runs a specific script located in the `scripts/` directory:
- **`make test-concurrency`**: Runs the high-concurrency stress test (`scripts/test_concurrency.go`).
- **`make test-race`**: Runs the simultaneous movement race condition test (`scripts/test_race_conditions.go`).
- **`make test-group`**: Runs the group state volatility test (`scripts/test_group_volatility.go`).
- **`make test-coalescing`**: Runs the TCP packet coalescing test (`scripts/test_coalescing.go`).
- **`make test-fragmentation`**: Runs the TCP packet fragmentation test (`scripts/test_fragmentation.go`).
- **`make test-abuse`**: Runs the command flooding and abuse test (`scripts/test_abuse.go`).
