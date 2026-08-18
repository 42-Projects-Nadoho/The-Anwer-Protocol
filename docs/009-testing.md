# Testing

### Protocol Handshake

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

### ABNF Syntax Compliance

To verify that message formats strictly follow the [ABNF syntax definitions](https://en.wikipedia.org/wiki/Augmented_Backus%E2%80%93Naur_form), test the server's error handling by sending malformed commands from the CLI client:

1. **Unknown Command:** Send an undefined command (e.g., `GIBBERISH`) and verify the server responds with a protocol-compliant error.
2. **Missing Arguments:** Send a command that requires arguments without them (e.g., `MOVE` or `CHAT`) and verify it is rejected.
3. **Invalid Formatting:** Send commands with invalid characters to test the robustness of the ABNF parser.

The server should gracefully return `ERR` messages and never crash, proving that the ABNF syntax definitions are strictly enforced.

### Commands Testing

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
| Resource Interaction | `QUEST <action>` | Manages specific quest interactions. |
| Resource Interaction | `QUESTS` | Lists active and completed quests. |

### Events Testing

While executing the above commands with multiple connected clients, verify that the server correctly pushes the following asynchronous events to the appropriate clients.

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
| `ROOM` | `CHAT` | `EVT ROOM CHAT CombatSys <details>` | Broadcasted during attack rounds. |
| `STATS` | `PLAYERS` | `EVT STATS players=<count>` | Updated server player count. |

If any command returns an `ERR` instead of `OK` (or if an event fails to broadcast to other connected clients), cross-reference the exact syntax with the RFC 42TAP specification document.
