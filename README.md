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
Architecture section explaining server design choices

```
add tree
```

## Protocol Implementation
Protocol Implementation section documenting any RFC deviations

## Combat System
Combat System section describing combat mechanics and design choices

## Quest System
Quest System section explaining quest progression and implementation

## World Design
World Design section describing world layout and NPC/item distribution

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
