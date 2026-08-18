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
