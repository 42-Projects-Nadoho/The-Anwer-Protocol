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

The Answer Protocol (TAP) is built around a robust client-server architecture written in Go, specifically designed to handle concurrent connections and real-time state synchronization over TCP. 

- **Concurrency Model:** The server employs a `readPump` and `writePump` goroutine per client connection to handle non-blocking I/O.
- **State Synchronization:** To prevent race conditions, the server acts as the single source of truth. All game state modifications (movement, combat, items) are serialized through a centralized `Hub` using a single operation channel (`Hub.do()`), eliminating the need for complex, scattered mutex locks.
- **Command Dispatch:** Incoming client commands are parsed and routed inline via a `switch` in each connection's `readPump` goroutine, rather than through a separate dispatcher object — the subject explicitly allows either approach, and the command set is small enough that a switch stays readable without the extra indirection.

> [!NOTE]
> Read the detailed Architecture documentation [here](docs/001-architecture.md).

## Protocol Implementation

The server strictly adheres to RFC 42TAP for the initial handshake, greeting (`S: OK hello proto=1`), and ABNF parsing (implemented in [`tap/src/internal/protocol/`](tap/src/internal/protocol/)). We implemented a few deliberate deviations, such as introducing custom error codes in [`errors.go`](tap/src/internal/protocol/errors.go) (e.g., `404 ErrRoomNotFound`, `902 ErrNotAuthenticated`) and custom combat/quest events, to enhance error handling clarity and support our extended gameplay mechanics.

> [!NOTE]
> Read the detailed Protocol Implementation documentation [here](docs/002-protocol-implementation.md).

## Combat System
Our combat system employs a synchronous, real-time architecture where actions are processed instantly in [`tap/src/internal/server/handlers_combat.go`](tap/src/internal/server/handlers_combat.go). Players engage hostile NPCs using the `ATTACK <npc>` command. Damage is calculated on the server, followed immediately by an automatic NPC counter-attack if it survives. Upon dropping to 0 HP, players are instantly teleported to the safety of the starting room with 50 HP.

> [!NOTE]
> Read the detailed Combat System documentation [here](docs/003-combat-system.md).

## Quest System
The quest engine (implemented in [`tap/src/internal/server/handlers_quests.go`](tap/src/internal/server/handlers_quests.go)) is designed to handle multiple objective types to keep progression engaging. Quests are acquired from `quest_giver` NPCs via the `QUEST <npc>` command and tracked using the `QUESTS` command. We implemented two distinct quest types: `multi_stage` (e.g., interacting with specific entities) and `defeat` (e.g., slaying a specific enemy). The server securely validates objectives on the backend—preventing client-side cheating—and distributes rewards such as unlocking new progression paths or restoring health.

> [!NOTE]
> Read the detailed Quest System documentation [here](docs/004-quest-system.md).

## World Design
The game world is purely data-driven, defined entirely within [`tap/data/world.yaml`](tap/data/world.yaml). The map consists of an 8-room central loop (Destiny Islands, Traverse Town, Olympus Coliseum, etc.) with branching optional areas to encourage exploration. Rooms are populated dynamically with obtainable items (like Potions and Keyblades) and NPCs that serve distinct roles: `quest_giver`, `dialogue`, and `enemy`.

> [!NOTE]
> Read the detailed World Design documentation [here](docs/005-world-design.md).

## Server Logging
The server utilizes Go's built-in `log/slog` package to implement structured, JSON-formatted logging. This allows for centralized monitoring of all incoming commands, broadcasted events, and dispatched errors (handled in [`tap/src/internal/server/logging.go`](tap/src/internal/server/logging.go)). Regular traffic is logged at the `Info` level, while protocol violations and `ERR` responses are elevated to `Warn` for easy anomaly detection and abuse prevention.

> [!NOTE]
> Read the detailed Server Logging documentation [here](docs/006-server-logging.md).

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
- **Brainstorming:** Drafting initial ideas for world topology, quest structures, and Kingdom Hearts theme integration.
- **Research:** Deep dives into RFC 42TAP specifications, ABNF syntax edge cases, and Go's concurrency primitives (`log/slog`, channels, goroutines).
- **Bugfixing:** Identifying and resolving race conditions during concurrent client accesses and TCP packet fragmentation issues.
- **Linting:** Enforcing idiomatic Go styling and checking for potential memory leaks or missing error handling.
- **Audit:** Reviewing the entire codebase against the subject to guarantee 100% compliance with mandatory requirements.
- **Documentation Assistance:** Drafting the detailed Markdown files, Mermaid diagrams, and tables found in the `docs/` folder to ensure clean, professional project presentation.
