# Acceptance Test Plan

An **Acceptance Test Plan (ATP)** is a formal document that outlines how a system or software product will be tested to verify if it meets business requirements and user needs. It acts as a guide for checking whether a system is ready for final delivery and client sign-off.

## Project Structure
The repository shall contain all mandatory deliverables:
- [ ] Building tool (Makefile or equivalent) at the root with targets/commands
- [ ] Targets exist for: install dependencies, run-server, run-client, run-client-gui, lint, clean
- [ ] TCP server implementation
- [ ] CLI client
- [ ] GUI client
- [ ] Static world data (YAML or JSON)
- [ ] README.md file at the root
- [ ] No unauthorised files are present

## README Documentation
The README.md file shall contain all the required sections:
- [ ] First line in italics: `This project has been created as part of the 42 curriculum by <login1>, <login2>...`. All group members’ logins must be present
- [ ] Description section: Clearly presenting the project goal and overview
- [ ] Instructions section: With compilation, installation, and execution information
- [ ] Resources section: Listing references and describing AI usage (if any)
- [ ] Architecture section: Explaining server design choices
- [ ] Protocol Implementation section: Documenting any RFC deviations
- [ ] Combat System section: Describing combat mechanics and design choices
- [ ] Quest System section: Explaining quest progression and implementation
- [ ] World Design section: Describing world layout and NPC/item distribution
- [ ] Server Logging section: Documenting logging implementation and monitoring capabilities
- [ ] Group Contributions section: Indicating each member's responsibilities
- [ ] Building and Running section: With detailed instructions
- [ ] Testing section: Explaining how to test functionality

## World Requirements

- [ ] [Room interconnection](/docs/world_layout.md): At least 8 interconnected rooms forming loops with at least one optional branch.
- [ ] NPC roles: At least 3 distinct NPC roles.
  - [ ] Dialogue NPCs: `leon`
  - [ ] Quest-giver NPCs: `yen_sid`, 
  - [ ] Enemy NPCs: `shadow_heartless`, `large_body`, `sephiroth`
- [ ] Item availability: At least 4 distinct items with at least 2 obtainable in-world: `potion`, `ether`, `keyblade`, `wayfinder`
- [ ] At least 2 implemented quests of different types: `find_wayfinder`, `defeat_shadow`
- [ ] Movement allows full circuit exploration (no "line-only" maps).
- [ ] Definition in world data: All NPCs and items referenced in rooms are properly defined in [world data](/data/world.yaml).

## Build
- [ ]  Project implementation language: Verify the project is implemented in one of the allowed languages: C, C++, Rust, Go, or Zig.
- [ ]  Building tool targets: On a clean environment, verify the building tool provides targets/commands for: install dependencies, run-server, run-client, run-client-gui, lint, clean.
- [ ]  `install` target: Run install must complete without errors.
- [ ]  Component run targets: Run-server, run-client, and run-client-gui must start the expected components.
- [ ]  `lint` and `clean` targets: Run lint and clean must complete without errors.
- [ ]  Allowed dependencies: Only appropriate dependencies for the chosen language are allowed (networking libraries, GUI toolkit, data parsers such as YAML/JSON parsers).

## Protocol Compliance
- [ ]  Connect with clients: Start the server and connect with both the CLI and GUI clients.
- [ ]  Greeting match: Verify that the greeting matches RFC 42TAP specifications.
- [ ]  Command execution: Execute every command and event defined in the RFC document to confirm they behave as specified.
- [ ]  Message formats: Verify that all message formats follow the ABNF syntax definitions in the RFC.
- [ ]  Malformed commands: Send malformed or invalid commands and ensure the server returns the correct ERR response without crashing.
- [ ]  Error codes: Check that error codes match exactly those specified in Section 7.2 of the RFC.

## Server Behaviour

- [ ]  Load world data: The server loads the world data and validates exits and references.
- [ ]  Room presence: The room is present.
- [ ]  Chat event broadcasting: The chat events are broadcast only to the intended recipients.
- [ ]  Handle client disconnect: If a client is disconnected while an event is being broadcast, the server removes the client cleanly and continues operating.

## CLI Client
Connect using the CLI client:
- [ ]  Interactive commands: Commands can be sent interactively.
- [ ]  Response time: Responses are displayed immediately.
- [ ]  Asynchronous events: Asynchronous events (chat, presence) appear while waiting for input.
- [ ]  Full flow: Full flow works as expected: `CONNECT > LOOK > MOVE > CHAT > WHO > QUIT`.

## GUI Client
Connect using the GUI client:
- [ ]  UI Display: Room details, items, NPCs, and exits are displayed.
- [ ]  Chat separation: Chat is separated by scope (Global, Room, Group).
- [ ]  Buttons function: Buttons for actions send the correct commands.
- [ ]  Player counts: Player counts in the room and on the server update in real time.
- [ ]  Responsiveness: The GUI remains responsive while receiving events.

## Robustness
Server and client edge cases:
- [ ]  Abrupt client disconnect: Disconnect a client abruptly: server must continue and remove the session.
- [ ]  Command spam: Send multiple commands quickly from different clients: responses remain correct.
- [ ]  Simultaneous moves: Attempt simultaneous MOVE actions from different clients and check for correct presence events.
- [ ]  Group management: Create and disband groups while members join/leave quickly: group state remains consistent.

## Network Features: Server
The server must handle:
- [ ]  Single packet: Multiple commands in a single TCP packet.
- [ ]  Split packets: Commands split across packets.
- [ ]  Unicode: Unicode characters in usernames or messages without encoding errors.
- [ ]  Control characters: Control characters in messages are either rejected or safely handled.

## Network Features: Inventory and NPC Interactions
- [ ]  `TAKE` command: Pick up items from rooms.
- [ ]  `DROP` command: Drop items from inventory.
- [ ]  `INVENTORY` command: List player's items.
- [ ]  `TALK` command: Interact with NPCs.
- [ ]  `ATTACK` command: Initiate combat with enemy NPCs.
- [ ]  `STATUS` command: Check player health and combat status.
- [ ]  `QUEST` command: Request quests from quest-giver NPCs.
- [ ]  `QUESTS` command: List active and completed quests.

## Data Integrity
- [ ]  `LOOK` outputs format: Are valid JSON, contain consistent IDs, and match the current game state.
- [ ]  Room definitions: All items and NPCs in rooms are defined in the world data.
- [ ]  Room traversal: Moving between rooms and back is consistent and accurate in `LOOK` data.

## Dynamic Item Management
- [ ]  Take an item from a room: The item disappears from `LOOK` output.
- [ ]  Attempt to take the same item again: The item returns `ITEM_NOT_FOUND` error.
- [ ]  Have a second player try to take the same item: The second player cannot take the same item.
- [ ]  Drop the item: The item reappears in the room for other players.
- [ ]  Test item IDs (e.g., "item.herbs"): Commands work correctly.
- [ ]  Test display names (e.g., "Herbs"): Commands work correctly.
- [ ]  Multi-word item names: Work correctly (e.g., "Loaf of Bread").

## Combat System
- [ ]  Use `STATUS` command at the beginning of the fight: Players start with 100 HP.
- [ ]  Use `ATTACK` command on enemy NPCs: Damage is dealt.
- [ ]  Enemy NPCs counter-attack: Reduce player HP.
- [ ]  Reach 0 HP: Players respawn at a safe location with reduced health.
- [ ]  Attack non-hostile NPCs: Cannot be attacked (should return `NPC_NOT_HOSTILE` error).
- [ ]  Combat logs: Combat results are logged and broadcast to relevant players.
- [ ]  Document combat mechanics: The group's combat mechanics design is documented in `README` with clear justification.

## Quest System
- [ ]  Receive quests: Use `QUEST` command on quest-giver NPCs to receive quests.
- [ ]  List quests: Use `QUESTS` command to list active and completed quests.
- [ ]  Quest types: Verify at least 2 different quest types are implemented (fetch item, defeat NPC, deliver item).
- [ ]  Quest completion: Test quest completion validation and reward systems.
- [ ]  NPCs without quests: Confirm NPCs without quests return `NO_QUEST_AVAILABLE` error.
- [ ]  Document quest mechanics: Verify the group's quest progression mechanics are documented in `README` with implementation approach.

## Server Logging
- [ ]  Client connections and disconnections: Logged with timestamps and IP addresses.
- [ ]  Commands received from clients: Logged with player name and parameters.
- [ ]  Server responses and error codes: Confirmed to be logged.
- [ ]  World state changes: Test that item movements, NPC interactions, combat results are logged.
- [ ]  Quest progress and completion: Events are logged.
- [ ]  Structured logging format: JSON recommended is used for easy parsing.
- [ ]  Log levels: (INFO, WARN, ERROR) are included for different event types.
- [ ]  Abuse patterns: Test that command flooding, rapid connections are monitored and logged.
- [ ]  Timestamps and streams: Verify all logs include precise timestamps and are written to appropriate output streams.
- [ ]  Performance impact: Check that logging does not significantly impact server performance or responsiveness.

## Recode Exercise
- [ ]  Request a brief modification to verify understanding: Ask the group to make a small change to one of the systems (e.g., modify NPC dialogue, adjust combat damage, add a simple quest step).
- [ ]  Modification feasibility: The modification should be feasible within a few minutes.
- [ ]  Team understanding: Verify that both group members understand the codebase and can explain their implementation choices.
- [ ]  Modification validation: Confirm the modification works correctly and doesn't break existing functionality.
