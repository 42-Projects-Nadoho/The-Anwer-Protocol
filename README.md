*This project has been created as part of the 42 curriculum by nadoho, spacotto.*


![TAP_KH_Images](/data/images/image_for_tap.jpg)

## Description
The Answer Protocol (TAP) is a collaborative multiplayer text adventure game, inspired by the classic Multi-User Dungeons (MUDs) of the early internet era. It features a persistent virtual world where players can connect in real-time to explore interconnected rooms, interact with characters, complete quests, and battle enemies together. Behind the scenes, the project showcases robust network programming through a custom-built server that manages the shared world and communicates seamlessly with players. Users can experience the adventure through two distinct interfaces: a nostalgic command-line client or a more accessible graphical application. Ultimately, TAP demonstrates the ability to design and build a complex, real-time multiplayer system from the ground up, blending technical architecture with engaging game design.

> [!IMPORTANT]
> This project is implemented in **Go**, strictly adhering to the language constraints of the subject (C, C++, Rust, Go, Zig are allowed; Python is forbidden). **Go (or  Golang)** has been chosen for its native concurrency model and excellent standard library for TCP networking.

## Instructions

This project uses `make` as its build tool. You can compile all components (server, CLI, GUI) simply by running:
```bash
make
```
For detailed, component-specific commands and targets, see the **Building and Running** section below.

## Architecture

The Answer Protocol (TAP) is built around a robust client-server architecture written in Go, specifically designed to handle concurrent connections and real-time state synchronization over TCP. The core server acts as the single source of truth, employing an event-driven dispatcher to manage player actions, combat, and world state across a shared environment. By isolating responsibilities into distinct modules—such as protocol serialization, world management, and network I/O—the system remains highly maintainable and scalable. This backend seamlessly supports two independent client implementations: a fast, text-based CLI and a richer, interactive GUI.

The primary components of the system include:

| Component | Path(s) | Description |
| :--- | :--- | :--- |
| **TCP server implementation** | `cmd/server/`, `internal/server/` | Handles concurrent connections and orchestrates the shared game state. |
| **CLI client** | `cmd/cli/` | A lightweight, text-based terminal interface. |
| **GUI client** | `cmd/gui/` | A richer graphical interface offering enhanced accessibility. |
| **Static world data** | `data/world.yaml` | Defines the rooms, NPCs, and items that construct the game world. |

> [!NOTE]
> Read the detailed Architecture documentation [here](docs/001-architecture.md).

## Protocol Implementation
Protocol Implementation section documenting any RFC deviations

## Combat System
Combat System section describing combat mechanics and design choices

## Quest System
Quest System section explaining quest progression and implementation

## World Design

The game world in The Answer Protocol is meticulously designed as a fully interconnected, non-linear environment that encourages deep exploration and cooperative gameplay. Moving away from simple linear paths, the layout features a central hub with branching loops and secret optional areas, ensuring players can freely traverse the world without hitting dead ends. This rich environment is populated by a diverse cast of NPCs—ranging from helpful dialogue characters and quest-givers to hostile enemies—and is scattered with unique items to discover, collect, and use. The deliberate distribution of these elements not only breathes life into the world but also seamlessly integrates with our dynamic combat and questing systems.

> [!NOTE]
> Read the detailed World Design documentation [here](docs/005-world%20-design.md).

## Server Logging
Server Logging section documenting logging implementation and monitoring capabilities

## Group Contributions
Group Contributions section indicating each member's responsibilities

## Building and Running

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

> [!TIP]
> To simply build all binaries without running them, use `make build` or `make`.

> [!NOTE]
> Read the detailed Building and Running documentation [here](docs/008-building-and-running.md).

## Testing
Testing section explaining how to test functionality

## Resources
Resources section listing references and describing AI usage (if any)

### AI usage
