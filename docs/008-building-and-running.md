# Building and Running

The project utilizes `make` for dependency management, compilation, and execution. The following targets verify our building tool meets all mandatory project requirements:

- **Install dependencies (`make install`):**
  Resolves and installs all required Go modules.
  ```bash
  make install
  ```

- **Run the server (`make run-server`):**
  Builds and starts the authoritative TAP server.
  ```bash
  make run-server
  ```

- **Run the CLI client (`make run-client`):**
  Builds and starts the command-line interface.
  ```bash
  make run-client
  ```

- **Run the GUI client (`make run-client-gui`):**
  Builds and starts the graphical user interface.
  ```bash
  make run-client-gui
  ```

- **Lint the code (`make lint`):**
  Checks formatting (`gofmt`) and runs static analysis (`go vet`) to ensure code quality.
  ```bash
  make lint
  ```

- **Clean build artifacts (`make clean`):**
  Removes compiled binaries and cleans the workspace.
  ```bash
  make clean
  ```

To simply build all binaries without running them, use `make build` or `make`.
