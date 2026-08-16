# Testing

## Verifying Protocol Handshake

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

## Commands and Events Testing

To confirm that the implementation behaves strictly as specified by the RFC, you must manually execute every mandatory command and verify the broadcasted events.

### Mandatory Commands

In either the CLI or GUI client, execute the following commands and verify the server's expected `S: OK ...` responses:

| Command | Description |
| :--- | :--- |
| `CONNECT <username>` | Connects a player to the server. |
| `LOOK` | Displays the current room description, exits, items, and NPCs. |
| `MOVE <direction>` | Moves the player to an adjacent room. |
| `CHAT <target> <message>` | Sends a message globally, to a room, or to a group. |
| `TAKE <item>` | Picks up an obtainable item from the current room. |
| `DROP <item>` | Drops an item from the inventory into the room. |
| `INVENTORY` | Lists items currently held by the player. |
| `TALK <npc>` | Initiates dialogue with an NPC. |
| `ATTACK <target>` | Initiates or continues combat with an enemy NPC. |
| `STATUS` | Displays current health and combat status. |
| `QUEST <action>` | Manages specific quest interactions. |
| `QUESTS` | Lists active and completed quests. |
| `WHO` | Lists players currently online or in the room. |
| `GROUP <action>` | Manages party creation and invites. |
| `QUIT` | Safely disconnects from the server. |

### Expected Broadcast Events

While executing the above commands with multiple connected clients, verify that the server correctly pushes the following asynchronous events to the appropriate clients:

| Event | Description |
| :--- | :--- |
| `EVT ROOM PRESENCE ENTER <username>` | Broadcasted when a player enters your room. |
| `EVT ROOM PRESENCE LEAVE <username>` | Broadcasted when a player leaves your room. |
| `EVT GLOBAL CHAT <username> <message>` | Broadcasted globally. |
| `EVT ROOM CHAT <username> <message>` | Broadcasted to the room. |
| `EVT COMBAT <details>` | Broadcasted during attack rounds. |

If any command returns an `ERR` instead of `OK` (or if an event fails to broadcast to other connected clients), cross-reference the exact syntax with the RFC 42TAP specification document.
