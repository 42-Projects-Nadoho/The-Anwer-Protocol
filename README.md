*This project has been created as part of the 42 curriculum by nadoho, spacotto.*


![TAP_KH_Images](tap/data/images/image_for_tap.jpg)

## Description
The Answer Protocol (TAP) is a collaborative multiplayer text adventure game, inspired by the classic Multi-User Dungeons (MUDs) of the early internet era. It features a persistent virtual world where players can connect in real-time to explore interconnected rooms, interact with characters, complete quests, and battle enemies together. Behind the scenes, the project showcases robust network programming through a custom-built server that manages the shared world and communicates seamlessly with players. Users can experience the adventure through two distinct interfaces: a nostalgic command-line client or a more accessible graphical application. Ultimately, TAP demonstrates the ability to design and build a complex, real-time multiplayer system from the ground up, blending technical architecture with engaging game design.

> [!IMPORTANT]
> This project is implemented in **Go**, strictly adhering to the language constraints of the subject (C, C++, Rust, Go, Zig are allowed; Python is forbidden). **Go (or  Golang)** has been chosen for its native concurrency model and excellent standard library for TCP networking.

## Instructions: Building and Running

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
  Checks formatting (`gofmt`) and runs static analysis (`go vet`).
  ```bash
  make lint
  ```

- **Clean build artifacts (`make clean`):**
  Removes compiled binaries and cleans the workspace.
  ```bash
  make clean
  ```

> [!TIP]
> To simply build all binaries without running them, use `make build` or `make`.
> Read the detailed Building and Running documentation [here](docs/008-building-and-running.md).

## Architecture
section explaining your server design choices (dispatcher/router vs inline handling, concurrency model, etc.).

> [!NOTE]
> Read the detailed Architecture documentation [here](docs/001-architecture.md).

## Protocol Implementation

section documenting any deviations from RFC 42TAP and justifying your choices.

> [!NOTE]
> Read the detailed Protocol Implementation documentation [here](docs/002-protocol-implementation.md).

## Combat System
section describing your turn-based combat mechanics, damage formulas, initiative order, and additional combat commands (DEFEND, FLEE, etc.).

## Quest System
section explaining your quest progression mechanics, completion validation, and reward systems.

## World Design
section describing your world layout, room connections, NPC roles, and item distribution.\

> [!NOTE]
> Read the detailed World Design documentation [here](docs/005-world%20-design.md).

## Server Logging

section documenting your logging implementation, including log format, event types, output destinations, and how to monitor server behavior and detect abuse patterns.

## Group Contributions
| Memeber | Role |
| :--- | :--- |
| nadoho | Server Implementation and CLI Client |
| spacotto | GUI Client and World Design |

> [!NOTE]
> Read the detailed Group Contributions documentation [here](docs/007-group-contributions.md).

## Testing

Our testing documentation covers how to manually verify the RFC protocol handshake, test multiplayer interactions, and validate the combat and quest systems. It also details our extensive suite of automated tests designed to verify high-concurrency stability, race-condition safety, group volatility, network-layer robustness (TCP coalescing and fragmentation), and security measures against abuse (command flooding and connection cycling).

> [!NOTE]
> Read the detailed Testing documentation [here](docs/009-testing.md) for step-by-step verification instructions.

## Resources
- **RFC 2119:** Key words for use in RFCs to Indicate Requirement Levels
- **Go Standard Library:** `net`, `log/slog`, `bufio`, `sync`

### AI usage
- 
