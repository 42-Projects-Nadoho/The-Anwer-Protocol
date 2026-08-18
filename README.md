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
> 
> Read the detailed Building and Running documentation [here](docs/008-building-and-running.md).

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

The server handshake, initial greeting (`S: OK hello proto=1`), and all subsequent command parsing strictly adhere to the RFC 42TAP specifications for both CLI and GUI clients.

> [!NOTE]
> Read the detailed Protocol Implementation documentation [here](docs/002-protocol-implementation.md).

## Combat System
Our combat system employs a synchronous, turn-less architecture where actions are processed in real-time as they arrive. Players engage hostile NPCs using the `ATTACK <npc>` command.
- **Damage Calculation & Response:** When an attack lands, the server calculates damage based on predefined stats. If the NPC survives, it immediately counter-attacks, dealing damage to the player's HP.
- **Death & Respawn:** If a player's HP drops to 0, they do not face permanent death. Instead, they are instantly teleported back to the `StartRoomID` (Destiny Islands) and respawn with 50 HP. This triggers broadcasted `PRESENCE LEAVE` and `PRESENCE ENTER` events to correctly update the world state for all clients in the affected rooms.

## Quest System
The quest engine is designed to handle multiple objective types to keep progression engaging:
- **Quest Types:** We implemented two distinct quest types: `multi_stage` (e.g., retrieving a specific item or finding a location) and `defeat` (e.g., slaying a specific enemy). 
- **Progression & Validation:** Quests are acquired from `quest_giver` NPCs via the `QUEST <npc>` command. The server securely validates objectives on the backend—preventing client-side cheating—and players can track their progress anytime using the `QUESTS` command.
- **Rewards:** Completing a quest grants rewards ranging from unlocked progression paths to full HP restoration, handled seamlessly by the server.

## World Design
The game world in The Answer Protocol is meticulously designed as a fully interconnected, non-linear environment that encourages deep exploration and cooperative gameplay. Moving away from simple linear paths, the layout features a central hub with branching loops and secret optional areas, ensuring players can freely traverse the world without hitting dead ends. This rich environment is populated by a diverse cast of NPCs—ranging from helpful dialogue characters and quest-givers to hostile enemies—and is scattered with unique items to discover, collect, and use. The deliberate distribution of these elements not only breathes life into the world but also seamlessly integrates with our dynamic combat and questing systems.

> [!NOTE]
> Read the detailed World Design documentation [here](docs/005-world%20-design.md).

## Server Logging
Server logging is implemented using Go's modern `log/slog` package, emitting structured JSON logs for robust monitoring. 
- **Monitoring & Abuse Detection:** We actively monitor incoming traffic for malicious behavior. The `checkFlood()` function tracks command frequency per session, while the `RecordConnection()` system flags rapid connection cycling (port-exhaustion attacks). If thresholds are exceeded, the server emits `WARN` level `possible_abuse` events while continuing to operate smoothly.

## Group Contributions
- **mosmond:** Core network architecture, TCP packet coalescing/fragmentation handling, and testing scripts.
- **tbaricau:** World design (`world.yaml`), GUI client implementation, and quest system.
- **nadoho:** Combat mechanics, respawn logic, and CLI client implementation.
- **spacotto:** Protocol serialization (RFC compliance), server logging, and abuse prevention systems.



## Testing

Our testing documentation covers how to manually verify the RFC protocol handshake, test multiplayer interactions, and validate the combat and quest systems.

> [!NOTE]
> Read the detailed Testing documentation [here](docs/009-testing.md) for step-by-step verification instructions.

## Resources
- **RFC 2119:** Key words for use in RFCs to Indicate Requirement Levels
- **Go Standard Library:** `net`, `log/slog`, `bufio`, `sync`

### AI usage
- Used LLMs (GitHub Copilot / ChatGPT) strictly for brainstorming, generating boilerplate structures (like the initial `world.yaml`), and debugging complex TCP race conditions. All core protocol implementations and architectural designs were written manually.
