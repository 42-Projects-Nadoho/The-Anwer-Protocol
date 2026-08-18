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
| Resource Interaction | `TALK <npc>` | Initiates dialogue with an NPC. |
| Resource Interaction | `ATTACK <target>` | Initiates or continues combat with an enemy NPC. |
| Resource Interaction | `STATUS` | Displays current health and combat status. |
| Resource Interaction | `QUEST <npc>` | Manages specific quest interactions. |
| Resource Interaction | `QUESTS` | Lists active and completed quests. |

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

# Server Behaviour
This section focuses on testing the internal game logic, mechanics, and state management of the server. While protocol compliance ensures we speak the right language, server behaviour testing ensures the actual "game" functions correctly—validating that events don't leak across boundaries, configuration files are structurally sound, and gameplay mechanics work as intended.

## World Data Validation Testing

The server actively validates the integrity of `data/world.yaml` during boot-up to ensure that no references point to missing entities. 

To test that this validation is working:
1. Open `data/world.yaml` in your editor.
2. Find any room (e.g., `destiny_islands`) and intentionally corrupt one of its exits by pointing it to a fake room ID (e.g., change `exits: north: traverse_town` to `exits: north: fake_room`).
3. Attempt to start the server with `make run-server`.
4. **Expected Result**: The server should immediately crash and refuse to start, printing an error such as: `Fatal error loading world: world data invalid: room "destiny_islands" exit "north" points to unknown room "fake_room"`.
5. Restore the broken exit back to normal after confirming.

You can perform similar tests by adding fake items to a room's `items` list or a fake NPC to a room's `spawns` list. The server will catch and reject all of them before booting up.

## Event Isolation Testing

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

## Disconnection Resilience Testing

To verify that the server's broadcast loop does not block or crash if a client abruptly disconnects during heavy event traffic:

1. **Setup**: Start the server and launch 2 separate CLI client tabs. Connect `Alice` and `Bob`.
2. **Trigger Broadcasts**: Have `Alice` initiate combat with an enemy (e.g., `ATTACK shadow_heartless_1`). This will trigger a recurring 3-second `EVT ROOM COMBAT` broadcast to both players.
3. **Abrupt Disconnect**: While the combat is running (and events are rapidly broadcasting), abruptly kill `Bob`'s terminal (e.g., press `Ctrl+C` or completely close the terminal window) instead of typing `QUIT`.
4. **Expected Result**: 
   - `Alice` should continue receiving combat events without any lag, interruption, or server crash. 
   - The server logs should show `Bob` disconnecting and being cleanly unregistered, proving that dead sockets do not stall the global event loop.

# CLI Client Testing

The CLI client (`bin/tap-cli`) must be tested to ensure it correctly handles user input and asynchronous server communication.

## Interactive Commands
To verify that players can organically send commands in real-time, test the CLI client's interactive prompt:
1. **Setup**: Start the server (`make run-server`) and the CLI client (`make run-client`).
2. **Connect**: Type `CONNECT <username>` and press Enter.
3. **Interact**: Type any command at your own pace.
4. **Expected Result**: Every time you press Enter, the CLI should instantly forward your command to the server.

## Immediate Responses & Asynchronous Events
To verify that the CLI correctly handles the background goroutine reading from the TCP socket:
1. **Setup**: Have `Alice` and `Bob` connected via two separate CLI tabs.
2. **Test Immediate Responses**: Have `Alice` type `LOOK`. She should immediately receive the room description.
3. **Test Asynchronous Events**: While `Alice` is idling and waiting for input at her terminal prompt, have `Bob` type `CHAT ROOM Hello!`.
4. **Expected Result**: `Alice` should immediately see `EVT ROOM CHAT Bob Hello!` pop up on her screen asynchronously, without it interrupting her own pending prompt input.

## Full Flow Integration
To verify that a full standard gameplay loop functions seamlessly from start to finish without breaking the client:

**Flow**: Execute the following commands in order:
```bash
> CONNECT Alice
> LOOK
> MOVE north
> CHAT GLOBAL Testing full flow
> WHO
> QUIT
```
**Expected Result**: The CLI should handle every command, display the results cleanly, and finally safely terminate its process upon receiving the `QUIT` acknowledgement.
