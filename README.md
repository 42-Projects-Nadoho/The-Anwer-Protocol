*This project has been created as part of the 42 curriculum by nadoho, spacotto.*


![TAP_KH_Images](/data/images/image_for_tap.jpg)

## Description
The Answer Protocol (TAP) is a collaborative multiplayer text adventure game, inspired by the classic Multi-User Dungeons (MUDs) of the early internet era. It features a persistent virtual world where players can connect in real-time to explore interconnected rooms, interact with characters, complete quests, and battle enemies together. Behind the scenes, the project showcases robust network programming through a custom-built server that manages the shared world and communicates seamlessly with players. Users can experience the adventure through two distinct interfaces: a nostalgic command-line client or a more accessible graphical application. Ultimately, TAP demonstrates the ability to design and build a complex, real-time multiplayer system from the ground up, blending technical architecture with engaging game design.

## Instructions

This project uses `make` as its build tool. The following commands are available to compile, install, and execute the project:

- **Install dependencies:**
  ```bash
  make install
  ```

- **Run the server:**
  ```bash
  make run-server
  ```

- **Run the CLI client:**
  ```bash
  make run-client
  ```

- **Run the GUI client:**
  ```bash
  make run-client-gui
  ```

- **Lint the code:**
  Checks formatting (`gofmt`) and runs static analysis (`go vet`).
  ```bash
  make lint
  ```

- **Clean build artifacts:**
  Removes compiled binaries and cleans the workspace.
  ```bash
  make clean
  ```

> [!TIP]
> To simply build all binaries without running them, use `make build` or `make`)

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
Building and Running section with detailed instructions

## Testing
Testing section explaining how to test functionality

## Resources
Resources section listing references and describing AI usage (if any)

### AI usage
