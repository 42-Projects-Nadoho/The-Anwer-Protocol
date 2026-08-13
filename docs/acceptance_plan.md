# Acceptance Test Plan

An **Acceptance Test Plan (ATP)** is a formal document that outlines how a system or software product will be tested to verify if it meets business requirements and user needs. It acts as a guide for checking whether a system is ready for final delivery and client sign-off.

## Project Structure
The repository shall contain all mandatory deliverables:
- [ ] Building tool (Makefile or equivalent) at the root with targets/commands:
  - [ ] install dependencies
  - [ ] run-server
  - [ ] run-client
  - [ ] run-client-gui
  - [ ] lint
  - [ ] clean
- [ ] TCP server implementation
- [ ] CLI client
- [ ] GUI client
- [ ] Static world data (YAML or JSON)
- [ ] README.md file at the root
- [ ] No unauthorised files are present

## README Documentation
The README.md file shall contain all the required sections:
- [ ] First line in italics: `This project has been created as part of the 42 curriculum by <login1>, <login2>...`. All group members’ logins must be present.
- [ ] Description section clearly presenting the project goal and overview
- [ ] Instructions section with compilation, installation, and execution information
- [ ] Resources section listing references and describing AI usage (if any)
- [ ] Architecture section explaining server design choices
- [ ] Protocol Implementation section documenting any RFC deviations
- [ ] Combat System section describing combat mechanics and design choices
- [ ] Quest System section explaining quest progression and implementation
- [ ] World Design section describing world layout and NPC/item distribution
- [ ] Server Logging section documenting logging implementation and monitoring capabilities
- [ ] Group Contributions section indicating each member's responsibilities
- [ ] Building and Running section with detailed instructions
- [ ] Testing section explaining how to test functionality

## World Requirements
- [ ] At least 8 interconnected rooms forming loops with at least one optional branch
- [ ] At least 3 distinct NPC roles: dialogue NPCs, quest-giver NPCs, and enemy NPCs
- [ ] At least 4 distinct items with at least 2 obtainable in-world
- [ ] At least 2 implemented quests of different types
- [ ] Movement allows full circuit exploration (no "line-only" maps)
- [ ] All NPCs and items referenced in rooms are properly defined in world data

## Build
- [ ] Verify the project is implemented in one of the allowed languages: C, C++, Rust, Go, or Zig
- [ ] On a clean environment, verify the building tool provides targets/commands for:
  - [ ] install dependencies
  - [ ] run-server
  - [ ] run-client
  - [ ] run-client-gui
  - [ ] lint
  - [ ] clean
- [ ] Run install must complete without errors
- [ ] Run-server, run-client, and run-client-gui must start the expected components
- [ ] Run lint and clean must complete without errors
- [ ] Only appropriate dependencies for the chosen language are allowed (networking libraries, GUI toolkit, data parsers such as YAML/JSON parsers)

## Protocol Compliance
- [ ] Start the server and connect with both the CLI and GUI clients
- [ ] Verify that the greeting matches RFC 42TAP specifications
- [ ] Execute every command and event defined in the RFC document to confirm they behave as specified
- [ ] Verify that all message formats follow the ABNF syntax definitions in the RFC
- [ ] Send malformed or invalid commands and ensure the server returns the correct ERR response without crashing
- [ ] Check that error codes match exactly those specified in Section 7.2 of the RFC

## Server Behaviour
- [ ] The server loads the world data and validates exits and references
- [ ] The room is present
- [ ] The chat events are broadcast only to the intended recipients
- [ ] Disconnect a client while an event is being broadcast: the server removes the client cleanly and continues operating

## CLI Client
Connect using the CLI client:
- [ ] Commands can be sent interactively
- [ ] Responses are displayed immediately
- [ ] Asynchronous events (chat, presence) appear while waiting for input
- [ ] Full flow works as expected: `CONNECT > LOOK > MOVE > CHAT > WHO > QUIT`

## GUI Client
Connect using the GUI client:
- [ ] Room details, items, NPCs, and exits are displayed
- [ ] Chat is separated by scope (Global, Room, Group)
- [ ] Buttons for actions send the correct commands
- [ ] Player counts in the room and on the server update in real time
- [ ] The GUI remains responsive while receiving events

## Robustness
Server and client edge cases:
- [ ] Disconnect a client abruptly: server must continue and remove the session
- [ ] Send multiple commands quickly from different clients: responses remain correct
- [ ] Attempt simultaneous MOVE actions from different clients and check for correct presence events
- [ ] Create and disband groups while members join/leave quickly: group state remains consistent

## Network Features: Server
The server must handle:
- [ ] Multiple commands in a single TCP packet
- [ ] Commands split across packets
- [ ] Unicode characters in usernames or messages without encoding errors
- [ ] Control characters in messages are either rejected or safely handled

## Network Features: Inventory and NPC Interactions
- [ ] `TAKE` command: pick up items from rooms
- [ ] `DROP` command: drop items from inventory
- [ ] `INVENTORY` command: list player's items
- [ ] `TALK` command: interact with NPCs
- [ ] `ATTACK` command: initiate combat with enemy NPCs
- [ ] `STATUS` command: check player health and combat status
- [ ] `QUEST` command: request quests from quest-giver NPCs
- [ ] `QUESTS` command: list active and completed quests

## Data Integrity
- [ ] 

## Dynamic Item Management
- [ ] 

## Combat System
- [ ] 

## Quest System
- [ ] 

## Server Logging
- [ ] 

## Recode Exercise
- [ ] 
