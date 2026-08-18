# Testing Guide

This document serves as the comprehensive testing manual for The Answer Protocol (TAP). It outlines both the automated tests designed to verify network robustness and abuse prevention, as well as the manual verification steps needed to audit protocol compliance, multiplayer interactions, combat, and quest mechanics.

## Automated Testing

The project includes a suite of automated Go scripts located in the `scripts/` directory. These can be easily executed via the `Makefile` to verify complex network features, concurrency edge-cases, and security measures:
- `make test-concurrency`: Simulates massive numbers of simultaneous connections and actions to ensure the server remains stable under load.
- `make test-race`: Aggressively triggers specific event overlaps to audit the server's goroutine safety and `hub` mutex synchronization.
- `make test-group`: Repeatedly creates, joins, and disbands groups with rapid volatility to ensure state consistency.
- `make test-coalescing`: Verifies the server correctly parses multiple commands arriving simultaneously in a single TCP packet.
- `make test-fragmentation`: Verifies the server correctly buffers incomplete commands split across multiple TCP packets.
- `make test-abuse`: Verifies the server's stability against rapid TCP connection cycling (port exhaustion) and its ability to detect and log application-layer command flooding.

# Protocol Compliance 
This section focuses on verifying that the server strictly adheres to the RFC 42TAP specification. This includes ensuring that the server strictly parses ABNF syntax, outputs the exact expected string formats for events, and strictly returns the correct standard error codes for all client interactions.

## Protocol Handshake

To verify that the initial greeting matches the RFC 42TAP specification:
1. Start the server in one terminal: 
   ```bash
   make run-server
   ```
2. In a new terminal, start the CLI client: 
   ```bash
   make run-client
   ```
3. In another terminal, start the GUI client: 
   ```bash
   make run-client-gui
   ```
4. Check the logs/output of both clients to verify that the very first message received is exactly:
   ```text
   S: OK hello proto=1
   ```

## ABNF Syntax Compliance

To verify that message formats strictly follow the [ABNF syntax definitions](https://en.wikipedia.org/wiki/Augmented_Backus%E2%80%93Naur_form), test the server's error handling by sending malformed commands from the CLI client:

1. **Unknown Command:** Send an undefined command (e.g., `GIBBERISH`) and verify the server responds with a protocol-compliant error.
2. **Missing Arguments:** Send a command that requires arguments without them (e.g., `MOVE` or `CHAT`) and verify it is rejected.
3. **Invalid Formatting:** Send commands with invalid characters to test the robustness of the ABNF parser.

The server should gracefully return `ERR` messages and never crash, proving that the ABNF syntax definitions are strictly enforced.

## Commands Testing

| Category | Command | Description |
| :--- | :--- | :--- |
| Core | `CONNECT <username>` | Connects a player to the server. |
| Core | `LOOK` | Displays the current room description, exits, items, and NPCs (JSON Structure). |
| Core | `MOVE <direction>` | Moves the player to an adjacent room. |
| Core | `QUIT` | Safely disconnects from the server. |
| Communication | `CHAT <scope> <message>` | Sends a message globally, to a room, or to a group. Available scopes: `GLOBAL`, `ROOM`, `GROUP`. |
| Communication | `WHO` | Lists players currently online or in the room. |
| Group Management | `GROUP CREATE` | Create a new player group. |
| Group Management | `GROUP INVITE <user>` | Invite a player to the current group. |
| Group Management | `GROUP JOIN <groupleader>` | Join an existing group. | 
| Group Management | `GROUP LEAVE` | Leave current group. |
| Resource Interaction | `TAKE <item>` | Picks up an obtainable item from the current room. |
| Resource Interaction | `DROP <item>` | Drops an item from the inventory into the room. |
| Resource Interaction | `INVENTORY` | Lists items currently held by the player. |
| NPC Interaction | `TALK <npc>` | Initiates dialogue with an NPC. |
| NPC Interaction | `ATTACK <target>` | Initiates or continues combat with an enemy NPC. |
| NPC Interaction | `STATUS` | Displays current health and combat status. |
| NPC Interaction | `QUEST <npc>` | Manages specific quest interactions. |
| NPC Interaction | `QUESTS` | Lists active and completed quests. |

## Events Testing

### Mandatory Events
| Category | Type | Event | Description |
| :--- | :--- | :--- | :--- |
| `ROOM` | `PRESENCE ENTER` | `EVT ROOM PRESENCE ENTER <username>` | Broadcasted when a player enters your room. |
| `ROOM` | `PRESENCE LEAVE` | `EVT ROOM PRESENCE LEAVE <username>` | Broadcasted when a player leaves your room. |
| `ROOM` | `CHAT` | `EVT ROOM CHAT <username> <message>` | Broadcasted to the room. |
| `GLOBAL` | `CHAT` | `EVT GLOBAL CHAT <username> <message>` | Broadcasted globally. |
| `GROUP` | `INVITE` | `EVT GROUP INVITE <username>` | Broadcasted when invited to a group. |
| `GROUP` | `JOIN` | `EVT GROUP JOIN <username>` | Broadcasted when a player joins the group. |
| `GROUP` | `LEAVE` | `EVT GROUP LEAVE <username>` | Broadcasted when a player leaves the group. |
| `GROUP` | `CHAT` | `EVT GROUP CHAT <username> <message>` | Broadcasted to group members. |
| `ROOM` | `COMBAT` | `EVT ROOM COMBAT DEFEAT <username> <npc_id>` | Custom event broadcasted when a hostile NPC is defeated. |
| `STATS` | `PLAYERS` | `EVT STATS players=<count>` | Updated server player count. |

### Custom Events
| Category | Type | Event | Description |
| :--- | :--- | :--- | :--- |
| `ROOM` | `CUSTOM` | `EVT ROOM COMBAT <details>` | Broadcasted during attack rounds. |
| `ROOM` | `CUSTOM` | `EVT ROOM RESPAWN The air shifts... <npc_id> has respawned!` | Broadcasted when a defeated NPC respawns after 30 seconds. |

## Error Testing

| Code | Name | Triggers |
| :--- | :--- | :--- |
| `201` | `NAME_IN_USE` | Attempt to `CONNECT` with a username that is already taken by another active player. |
| `301` | `NO_EXIT` | Attempt to `MOVE` in a direction that does not exist in the current room
| `301` | `DIRECTION_REQUIRED`| Attempt to `MOVE` omitting the direction argument. |
| `401` | `NOT_IN_GROUP` | Attempt to use group commands (`GROUP INVITE`, `CHAT GROUP`, `GROUP LEAVE`) while not currently in a group. |
| `402` | `ALREADY_IN_GROUP` | Attempt to `GROUP INVITE` or `GROUP JOIN` while already being part of a group. |
| `404` | `ITEM_NOT_FOUND` | Attempt to `TAKE` an item that does not exist in the current room. |
| `404` | `ITEM_NOT_IN_INVENTORY` | Attempt to `DROP` an item that does not exist in the player's inventory. |
| `404` | `NPC_NOT_FOUND` | Attempt to `TALK` or `ATTACK` an NPC that does not exist in the current room. |
| `405` | `NPC_NOT_HOSTILE` | Attempt to `ATTACK` a friendly NPC. |
| `406` | `NO_QUEST_AVAILABLE` | Attempt to `QUEST <npc>` when no quest is available from the NPC or conditions are unmet. |
| `900` | `CONNECTION_FAILED` | Dial the TCP server fails. |
| `901` | `SEND_FAILED` | Triggered internally by the server if internal JSON serialisation fails while packaging complex data structures. |
| `902` | `NOT_AUTHENTICATED` | Attempt to execute any gameplay command before successfully connecting via `CONNECT`. |
| `903` | `UNKNOWN_COMMAND` | Send a command format that the server's protocol parser does not recognise. |

# Server & Client Testing
This section focuses on testing the internal game logic, mechanics, and state management of the server. While protocol compliance ensures we speak the right language, server behaviour testing ensures the actual "game" functions correctly—validating that events don't leak across boundaries, configuration files are structurally sound, and gameplay mechanics work as intended.

## Server: World Data Validation Testing
The server actively validates the integrity of `data/world.yaml` during boot-up to ensure that no references point to missing entities. 

To test that this validation is working:
1. Open `data/world.yaml` in your editor.
2. Find any room (e.g., `destiny_islands`) and intentionally corrupt one of its exits by pointing it to a fake room ID (e.g., change `exits: north: traverse_town` to `exits: north: fake_room`).
3. Attempt to start the server with `make run-server`.
4. **Expected Result**: The server should immediately crash and refuse to start, printing an error such as: `Fatal error loading world: world data invalid: room "destiny_islands" exit "north" points to unknown room "fake_room"`.
5. Restore the broken exit back to normal after confirming.

You can perform similar tests by adding fake items to a room's `items` list or a fake NPC to a room's `spawns` list. The server will catch and reject all of them before booting up.

## Server: Event Isolation Testing
To ensure that room presence and room chat events do not "leak" to players in other locations across the world, perform the following multi-client test:

1. **Setup**: Start the server and launch 3 separate client instances (e.g., three CLI tabs).
2. **Connect**: Connect all 3 players: `CONNECT Alice`, `CONNECT Bob`, and `CONNECT Charlie`. By default, they will all spawn in the exact same starting room.
3. **Isolate**: Have Charlie leave the room by typing `MOVE north` (or any valid exit). Charlie is now in a different room than Alice and Bob.
4. **Test Room Chat**: Have Alice type `CHAT ROOM Hello everyone!`. 
   - **Expected Result**: Bob should receive `EVT ROOM CHAT Alice Hello everyone!`. Charlie should receive **nothing**.
5. **Test Presence Leak**: Have Alice type `MOVE south` (or any exit that Charlie is not in).
   - **Expected Result**: Bob should receive `EVT ROOM PRESENCE LEAVE Alice`. Charlie should receive **nothing**.
6. **Test Presence Arrival**: Have Alice move into the room Charlie is currently standing in.
   - **Expected Result**: Charlie should receive `EVT ROOM PRESENCE ENTER Alice`. Bob should receive **nothing**.

## Server: Disconnection Resilience Testing
To verify that the server's broadcast loop does not block or crash if a client abruptly disconnects during heavy event traffic:

1. **Setup**: Start the server and launch 2 separate CLI client tabs. Connect `Alice` and `Bob`.
2. **Trigger Broadcasts**: Have `Alice` initiate combat with an enemy (e.g., `ATTACK shadow_heartless_1`). This will trigger a recurring 3-second `EVT ROOM COMBAT` broadcast to both players.
3. **Abrupt Disconnect**: While the combat is running (and events are rapidly broadcasting), abruptly kill `Bob`'s terminal (e.g., press `Ctrl+C` or completely close the terminal window) instead of typing `QUIT`.
4. **Expected Result**: 
   - `Alice` should continue receiving combat events without any lag, interruption, or server crash. 
   - The server logs should show `Bob` disconnecting and being cleanly unregistered, proving that dead sockets do not stall the global event loop.

## CLI Client: Interactive Commands
To verify that players can organically send commands in real-time, test the CLI client's interactive prompt:
1. **Setup**: Start the server (`make run-server`) and the CLI client (`make run-client`).
2. **Connect**: Type `CONNECT <username>` and press Enter.
3. **Interact**: Type any command at your own pace.
4. **Expected Result**: Every time you press Enter, the CLI should instantly forward your command to the server.

## CLI Client: Immediate Responses & Asynchronous Events
To verify that the CLI correctly handles the background goroutine reading from the TCP socket:
1. **Setup**: Start the server (`make run-server`) and the CLI client (`make run-client`).
1. **Connect**: Have `Alice` and `Bob` connected via two separate CLI tabs.
2. **Test Immediate Responses**: Have `Alice` type `LOOK`. She should immediately receive the room description.
3. **Test Asynchronous Events**: While `Alice` is idling and waiting for input at her terminal prompt, have `Bob` type `CHAT ROOM Hello!`.
4. **Expected Result**: `Alice` should immediately see `EVT ROOM CHAT Bob Hello!` pop up on her screen asynchronously, without it interrupting her own pending prompt input.

## CLI Client: Full Flow Integration
To verify that a full standard gameplay loop functions seamlessly from start to finish without breaking the client:

**Setup**: Start the server (`make run-server`) and the CLI client (`make run-client`).
**Flow**: Execute the following commands in order:
```bash
> CONNECT Alice
> LOOK
> MOVE north
> CHAT GLOBAL Testing full flow
> WHO
> QUIT
```

**Expected Result**: 
```bash
make run-client
>>> Building TAP CLI client...
go build -o bin/tap-cli ./cmd/cli
>>> Starting TAP CLI client...
bin/tap-cli -addr x.x.x.x:xxxx
=== Welcome to TAP ===
Connecting to x.x.x.x:xxxx...
Connected! You can now type your commands.
OK hello proto=1
> CONNECT Alice
OK connected
> LOOK
OK {"room":{"id":"destiny_islands","name":"Destiny Islands","description":"A beautiful tropical island where journeys begin. The sun is shining brightly.","exits":{"east":"traverse_town","north":"secret_cave","south":"hollow_bastion"}},"players":["Alice"],"items":["potion"],"npcs":["yen_sid"]}
> MOVE north
OK room=secret_cave
> CHAT GLOBAL Testing full flow   
OK
EVT GLOBAL CHAT Alice Testing full flow
> WHO
OK players=1 (Alice)
> QUIT
Goodbye!
OK bye
> %
```

## GUI Client
To ensure the web-based graphical client (`bin/tap-gui`) correctly parses JSON and visually renders the game state:

1. **Setup**: Start the server (`make run-server`) and the GUI client (`make run-client-gui`).
2. **Connect**: Open your browser to the local GUI address and use the interface to connect to the game.
3. **Elements Checklist**: 
   - [x] **Room details, items, NPCs, and exits** are accurately displayed in the UI panels.
   - [x] **Chat** is visually separated by scope (Global, Room, Group) and renders properly.
   - [x] **Buttons** for actions (like moving or looking) send the correct commands to the server and update the UI.
   - [x] **Player counts** (both in the room and on the server globally) update in real-time as other clients connect and move around.
4. **Performance**: Ensure the GUI remains fully responsive, scrollable, and clickable even while actively receiving heavy bursts of events (like combat).

## Edge Cases
All the edge cases apply to the server and both clients. So, you need to start the server (`make run-server`) and both clients (`make run-client` and `make run-client-gui`)

### Abrupt Client Disconnection
The server's broadcast loop shall not block or crash if a client abruptly disconnects during heavy event traffic.
2. **Setup**: Connect `Alice` and `Bob`.
3. **Trigger Heavy Broadcasts**: Have `Alice` initiate combat with an enemy. This triggers a recurring 3-second broadcast to both players.
4. **Abrupt Disconnect**: While the combat is running, aggressively kill `Bob`'s terminal (e.g., `Ctrl+C` or closing the window) instead of typing `QUIT`.
5. **Expected Result**: `Alice` should continue receiving combat events without lag. The server logs should show `Bob` disconnecting and being cleanly unregistered, proving dead sockets do not stall the global event loop.

### High-Concurrency Input Handling
The server shall correctly queue and process rapid-fire commands from multiple clients simultaneously without corrupting state or dropping responses. 
- **Test:** Run `make test-concurrency` in a separate terminal while the server is running.
- **Expected Result:** The server should not crash. Both clients should receive the correct number of responses back in exactly the order they were received. 

### Simultaneous Movement Race Conditions
The server shall correctly broadcast presence events without race conditions when multiple clients move at the exact same time.
- **Test:** Run `make test-race` in a separate terminal while the server is running.
- **Expected Result**: `Charlie` (still in the starting room) should receive both `EVT ROOM PRESENCE LEAVE Alice` and `EVT ROOM PRESENCE LEAVE Bob` in rapid succession. Neither `Alice` nor `Bob` should receive each other's leave events because they exited the room at the same time.

### Group State Volatility
The server shall ensure the group management data structure does not break or panic when members join and leave chaotically.
- **Test:** Run `make test-group` in a separate terminal while the server is running.
- **Expected Result**: The server should handle the locks cleanly. If `Alice` leaves first, `Bob`'s join command should return `ERR 401 NOT_IN_GROUP` (or similar) because the group dissolved. If `Bob` joins first, he should receive the `EVT GROUP JOIN Bob` followed immediately by `EVT GROUP LEAVE Alice`. The server should not panic or encounter a map-read concurrent exception.

# Features

## Server Features

### TCP Packet Coalescing
The server shall correctly parse and execute multiple newline-terminated commands that arrive concatenated within a single TCP packet.
- **Test:** Run `make test-coalescing` in a separate terminal while the server is running.
- **Expected Result**: The server processes both commands sequentially and returns two distinct responses.

### TCP Packet Fragmentation
The server shall correctly buffer and reconstruct commands that are split across multiple TCP packets.
- **Test:** Run `make test-fragmentation` in a separate terminal while the server is running.
- **Expected Result**: The server waits until the newline `\n` is received, then correctly executes `CHAT GLOBAL test frag` without throwing a parse error.

### Unicode Encoding Resilience
The server shall safely process, store, and broadcast Unicode characters (like emojis or non-Latin alphabets) without mangling the text or crashing.
- **Test:** Connect to the server, login with the username `ユーザー`, and send the command `CHAT GLOBAL 🌍 Hello`.
- **Expected Result**: Other connected clients receive `EVT GLOBAL_CHAT ユーザー 🌍 Hello` perfectly intact.

### Control Character Sanitization
The server shall reject or safely strip unprintable ASCII control characters (like `\x00` null bytes or ANSI escape sequences) to prevent terminal injection attacks.
- **Test:** Send a chat message containing raw escape codes (e.g., `\x1b[31mRedText`) or null bytes.
- **Expected Result**: The server either sanitizes the characters out, or immediately responds with an `ERR` (e.g., `ERR 400 INVALID_COMMAND_FORMAT`). It must never crash or forward raw escape sequences that corrupt the recipient's CLI.

# Gameplay Features

## Inventory Interaction Features

### TAKE Command
The server shall allow clients to pick up items that are present in the current room.
- **Test:** Connect via the CLI client and send `TAKE potion` in a room that contains a potion.
- **Expected Result:** The server returns an `OK` confirming the item was taken. The item is removed from the room, and all other players in the room receive an `EVT ROOM ITEM_TAKEN <username> potion`.

### DROP Command
The server shall allow clients to drop items from their inventory into the current room.
- **Test:** Send `DROP potion` while possessing a potion in your inventory.
- **Expected Result:** The server returns an `OK` confirming the drop. The item is added back to the room, and all other players in the room receive an `EVT ROOM ITEM_DROPPED <username> potion`.

### INVENTORY Command
The server shall accurately return a list of all items currently held by the player.
- **Test:** Send the `INVENTORY` command.
- **Expected Result:** The server returns an `OK` accompanied by a JSON payload listing the items (e.g., `OK {"items":["potion"]}`). If the inventory is empty, it should return an empty list or appropriate empty state.

### Dynamic Item Management
The server shall dynamically manage item ownership and room states flawlessly to prevent item duplication and ensure robust parsing.
- **Test:** 
  1. Send `LOOK` to confirm an item exists. Send `TAKE <item>`, then send `LOOK` again to verify it is removed from the room payload.
  2. Attempt to send `TAKE <item>` again for the same item.
  3. Have a second connected player attempt to send `TAKE <item>` for that item.
  4. Send `DROP <item>`, and have the second player send `LOOK` to verify it has reappeared.
  5. Repeat the above interactions explicitly testing both the item's raw ID (e.g., `item.herbs`), its capitalized display name (e.g., `Herbs`), and a multi-word item name (e.g., `Loaf of Bread`).
- **Expected Result:** The first `TAKE` succeeds. Any subsequent `TAKE` (by either player) returns an error (e.g., `ERR 404 ITEM_NOT_FOUND`). Dropping the item cleanly restores it to the room. The command parser successfully maps and interacts with the item regardless of whether the player inputted the ID, the exact name, or a multi-word string.

## NPC Interaction Features

### TALK Command
The server shall allow clients to interact with friendly or neutral NPCs.
- **Test:** Connect via the CLI client and send `TALK yen_sid` in a room that contains the NPC.
- **Expected Result:** The server returns an `OK` followed by the NPC's dialog string in the payload.

### ATTACK Command
The server shall allow clients to initiate combat with enemy NPCs.
- **Test:** Send `ATTACK heartless` in a room containing the enemy NPC.
- **Expected Result:** The server returns an `OK` acknowledging combat has started, and begins broadcasting periodic `EVT COMBAT` events detailing the ongoing battle until the enemy or player is defeated.

### STATUS Command
The server shall accurately return the player's health and combat status.
- **Test:** Send the `STATUS` command before engaging in combat, and again while actively in combat.
- **Expected Result:** The server returns an `OK` along with a JSON payload indicating current HP, max HP, and whether the player is currently engaged in an active battle.

### Combat System
The server shall robustly handle combat mechanics, damage tracking, and respawning.
- **Test:** 
  1. Send `STATUS` to verify players start with exactly 100 HP.
  2. Send `ATTACK <enemy>` and observe the combat loop to confirm damage is dealt to the enemy.
  3. Allow the combat loop to run and verify the enemy NPC successfully counter-attacks, reducing your HP.
  4. Allow the player's HP to reach 0 to verify the player respawns at a safe location (e.g., `destiny_islands`) with appropriately restored/reduced health.
  5. Attempt to send `ATTACK yen_sid` (or another non-hostile NPC).
  6. Verify server logs and other clients in the room to ensure `EVT COMBAT` is correctly broadcast.
  7. Check the project's root `README.md` to ensure the group's combat mechanics and design decisions are explicitly documented with clear justifications.
- **Expected Result:** The player starts at 100 HP. Combat dynamically reduces both participant's health pools over time. Dying instantly teleports the player and resets combat state. Attacking friendly NPCs safely returns an `ERR 400 NPC_NOT_HOSTILE` (or similar) without starting combat. All broadcasts and logs perfectly reflect the battle.

### Quest System
The server shall robustly handle quest progression, completion logic, and reward distribution.
- **Test:**
  1. Send `QUEST <npc>` on an eligible quest-giver NPC to receive a quest.
  2. Send the `QUESTS` command to verify the active quest is accurately listed in your log.
  3. Verify the game world configuration explicitly implements at least two distinct quest types (e.g., "fetch item", "defeat NPC", or "deliver item").
  4. Perform the necessary quest objectives and interact with the quest giver again to test completion validation and ensure the reward system triggers correctly.
  5. Attempt to send `QUEST` to a standard NPC that does not offer quests.
  6. Check the project's root `README.md` to verify the group's quest progression mechanics and implementation approach are clearly documented.
- **Expected Result:** The player successfully requests, tracks, and completes dynamic quests. Rewards are correctly allocated upon completion validation. Interacting with standard NPCs securely returns an `ERR` (e.g., `NO_QUEST_AVAILABLE`). The required technical documentation is present in the README.

# Data Integrity

### JSON Validation and Consistent IDs
The server shall ensure that all `LOOK` outputs are perfectly valid JSON, containing correct string and list identifiers that precisely match the active world configuration.
- **Test:** Connect to the server, send the `LOOK` command, and inspect the JSON payload.
- **Expected Result:** The server returns `OK` accompanied by a flawless JSON payload. The `id`, `name`, `exits`, `players`, `items`, and `npcs` keys must be present, and all IDs correctly reflect the actual runtime game state.

### World Data Referential Integrity
The server shall ensure that all items, NPCs, and exit destinations advertised in a room's configuration correspond to valid, defined entities in the global world data schema.
- **Test:** Send `LOOK` to get the list of entities in a room. Proceed to interact with them (e.g., `TAKE` every item, `TALK` to every NPC).
- **Expected Result:** The server correctly recognizes every item and NPC without throwing internal errors. The game does not panic, and it does not return `ERR 404 NOT_FOUND` for anything that was explicitly advertised as being present in the room.

### Spatial Consistency
The server shall ensure that navigating back and forth between interconnected rooms results in consistent, accurate `LOOK` data.
- **Test:** Send `LOOK` and record the room's exits and items. Send `MOVE <direction>` to enter a neighboring room. Send `MOVE <opposite_direction>` to return to the original room, and send `LOOK` again.
- **Expected Result:** The data returned by the final `LOOK` command precisely matches the data from the first `LOOK` command (assuming no other players interacted with the room). The geometric map safely preserves spatial consistency.

# Server Logging
The server must implement a comprehensive, structured logging system (MANDATORY).
- **Test:** Actively play the game, trigger errors, engage in combat, complete quests, and simulate abuse (e.g., command flooding) while monitoring the server's output streams.
- **Expected Result:** The server logs must rigorously adhere to the following criteria:
  - **Connections:** All client connections and disconnections strictly log the timestamp and IP address.
  - **Inputs:** Every command received logs the player name, action, and parameters.
  - **Outputs:** All server responses and outgoing error codes are logged.
  - **State Changes:** Critical world state changes (item movements, combat results, quest progression/completion) are logged.
  - **Formatting:** The logs utilize a structured, easily parsable format (JSON recommended).
  - **Levels:** The logs include explicit severity levels (e.g., `INFO`, `WARN`, `ERROR`).
  - **Security:** Potential abuse patterns (like command flooding or rapid connection cycling) are successfully monitored and logged.
  - **Streams:** The logging system writes to appropriate output streams with precise timestamps.
  - **Performance:** The volume of logging does not bottleneck or significantly impact the server's responsiveness under heavy load.
