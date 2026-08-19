# Protocol Implementation

The Answer Protocol (TAP) is strictly designed to follow the RFC 42TAP specifications. Our implementation ensures robust, concurrent handling of TCP connections with full compliance to the required message formats and error codes.

## Handshake and Greeting
The server handshake and initial greeting strictly adhere to the RFC 42TAP specifications upon any client connection (whether CLI or GUI). Immediately after a connection is established, the server sends the expected protocol version greeting before awaiting the `CONNECT` command.

Example greeting:
```text
S: OK hello proto=1
```

## ABNF Syntax Compliance

The server's message parser is strictly built against the ABNF syntax definitions provided in the RFC:
- All incoming client commands are validated against their expected ABNF structure.
- Malformed commands, missing arguments, or invalid syntax immediately trigger protocol-compliant error responses (e.g., `ERR`) rather than causing server instability.
- All outbound events and responses from the server are formatted strictly according to the RFC ABNF rules, ensuring perfect compatibility with any compliant third-party clients.

## Deviations from RFC 42TAP

While our core architecture strictly follows the RFC, we made a few deliberate extensions to the protocol to improve debugging clarity and support our extended combat and questing mechanics:

### Custom Error Codes
To provide more actionable feedback to clients, we introduced specific error codes beyond the generic RFC errors:
- `400 ErrInvalidCommandFormat`: Triggered when ABNF parsing fails for a recognized command.
- `404 ErrRoomNotFound`: Used internally and externally when a requested exit or lookup fails.
- `902 ErrNotAuthenticated`: Returned when a client attempts to execute world commands without completing the `CONNECT` handshake.
- `903 ErrUnknownCommand`: Returned when a command completely falls outside the known router dictionary.

### Custom Events
To support our real-time combat system without breaking the RFC's standard text broadcasting, we extended the `ROOM COMBAT` event channel:
- `DEFEAT <username> <npc_id>`: An explicit event dispatched to the room when an NPC is slain, allowing our GUI client to trigger specific animations and update the local room state seamlessly.

The RFC's `PRESENCE ENTER`/`PRESENCE LEAVE` events are scoped to `ROOM`, so a client only learns about players entering or leaving its own room — there is no way to keep a server-wide online count live without polling `WHO` on a timer. We added a server-wide counterpart:
- `EVT GLOBAL PRESENCE ENTER <username>` / `EVT GLOBAL PRESENCE LEAVE <username>`: Broadcast to every connected client (in addition to, not instead of, the existing `ROOM PRESENCE` event) whenever a player connects or disconnects. The GUI client uses this purely as a trigger to re-issue `WHO`, so the top-bar player count updates the instant someone joins or leaves anywhere on the server, not just in the viewer's own room.

### GROUP Semantics

The RFC's `GROUP JOIN <leader-name>` syntax and `EVT GROUP INVITE <leader>` event both name their argument "leader", but the RFC never defines what a leader is or what privileges it has — this is the only place in the document the word appears. We deliberately did not build a leader role: any current member of a group may invite others, and a group persists as long as it has at least one member, regardless of who created it.

**`GROUP JOIN` takes a group id, not a username.** Because our groups have no leader role and can outlive their creator, resolving "join by naming a member" breaks as soon as that member has since left or disconnected, even though the group itself may still exist. We instead require the server-generated opaque group id returned by `GROUP CREATE` (e.g. `group-3`). **This breaks wire compatibility with a strictly RFC-literal client from another group for this one command** — a conscious tradeoff, not an oversight.

**`EVT GROUP INVITE` carries the group id in addition to the inviter.** Format: `EVT GROUP INVITE <inviter-username> <group-id>` instead of the RFC's `EVT GROUP INVITE <leader>`. This is a direct consequence of the previous point: since `JOIN` needs a group id, the invitee has to learn it from somewhere, and the invite event is the only place that can carry it.

**Invitations are single-use.** A pending invitation (tracked per-client in `isInvited []string`) is consumed the moment `GROUP JOIN` succeeds against it; re-joining after leaving requires a fresh invite.
